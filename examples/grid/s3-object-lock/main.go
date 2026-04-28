package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bedag/storagegrid-sdk-go/client"
	"github.com/bedag/storagegrid-sdk-go/models"
)

// End-to-end S3 Object Lock example.
//
//  1. Reads the grid-wide compliance-global settings via the GridClient.
//  2. Optionally enables S3 Object Lock grid-wide (irreversible).
//  3. Updates the tenant policy to allow S3 Object Lock for the target tenant
//     and caps the maximum retention a bucket may specify.
//  4. Creates a bucket via the TenantClient with S3 Object Lock enabled and a
//     default compliance-mode retention period.
//  5. Reads the per-bucket Object Lock settings back using the new
//     /org/containers/{bucketName}/object-lock endpoint.
//
// See: https://docs.netapp.com/us-en/storagegrid/ilm/managing-objects-with-s3-object-lock.html
func main() {
	endpoint := os.Getenv("STORAGEGRID_ENDPOINT")
	gridUser := os.Getenv("STORAGEGRID_USERNAME")
	gridPass := os.Getenv("STORAGEGRID_PASSWORD")
	tenantUser := os.Getenv("STORAGEGRID_TENANT_USERNAME")
	tenantPass := os.Getenv("STORAGEGRID_TENANT_PASSWORD")
	accountID := os.Getenv("STORAGEGRID_ACCOUNT_ID")
	skipSSL := os.Getenv("STORAGEGRID_SKIP_SSL") == "true"
	apply := os.Getenv("STORAGEGRID_APPLY") == "true"

	if endpoint == "" || gridUser == "" || gridPass == "" {
		log.Fatal("Required env vars: STORAGEGRID_ENDPOINT, STORAGEGRID_USERNAME, STORAGEGRID_PASSWORD")
	}

	ctx := context.Background()

	gridClient, err := client.NewGridClient(buildOpts(endpoint, gridUser, gridPass, nil, skipSSL)...)
	if err != nil {
		log.Fatalf("Failed to create grid client: %v", err)
	}

	// 1. Inspect grid-wide settings.
	fmt.Println("🔍 Grid-wide S3 Object Lock settings (/grid/compliance-global)")
	current, err := gridClient.S3ObjectLock().Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get S3 Object Lock settings: %v", err)
	}
	printGridSettings("Current", current)

	// 2. Optionally enable grid-wide.
	if apply {
		enable := true
		fmt.Println("\n✏️  Enabling complianceEnabled grid-wide (irreversible)...")
		updated, err := gridClient.S3ObjectLock().Update(ctx, &models.S3ObjectLock{ComplianceEnabled: &enable})
		if err != nil {
			log.Fatalf("Failed to update S3 Object Lock settings: %v", err)
		}
		printGridSettings("Updated", updated)
		current = updated
	} else {
		fmt.Println("\nℹ️  Set STORAGEGRID_APPLY=true to enable S3 Object Lock grid-wide.")
	}

	if current.ComplianceEnabled == nil || !*current.ComplianceEnabled {
		fmt.Println("\n⚠️  Grid-wide S3 Object Lock is not enabled — skipping bucket demo.")
		return
	}

	if tenantUser == "" || tenantPass == "" || accountID == "" {
		fmt.Println("\nℹ️  Set STORAGEGRID_TENANT_USERNAME, STORAGEGRID_TENANT_PASSWORD, and STORAGEGRID_ACCOUNT_ID to run the bucket demo.")
		return
	}

	// 3. Update the tenant policy: allow S3 Object Lock and cap max retention.
	//    Per-tenant settings live on TenantPolicy and gate what tenants can do
	//    with S3 Object Lock once it is enabled grid-wide.
	fmt.Println("\n🔐 Updating tenant policy to allow S3 Object Lock...")
	tenant, err := gridClient.Tenant().GetById(ctx, accountID)
	if err != nil {
		log.Fatalf("Failed to get tenant %s: %v", accountID, err)
	}
	if tenant.Policy == nil {
		tenant.Policy = &models.TenantPolicy{}
	}
	allow := true
	maxYears := 10
	tenant.Policy.AllowComplianceMode = &allow
	tenant.Policy.MaxRetentionYears = &maxYears
	tenant.Policy.MaxRetentionDays = nil // no per-day cap (years cap is sufficient)
	if _, err := gridClient.Tenant().Update(ctx, tenant); err != nil {
		log.Fatalf("Failed to update tenant policy: %v", err)
	}
	fmt.Printf("  ✅ allowComplianceMode=true, maxRetentionYears=%d on tenant %s\n", maxYears, accountID)

	tenantClient, err := client.NewTenantClient(buildOpts(endpoint, tenantUser, tenantPass, &accountID, skipSSL)...)
	if err != nil {
		log.Fatalf("Failed to create tenant client: %v", err)
	}

	// 4. Create a bucket with S3 Object Lock enabled.
	bucketName := fmt.Sprintf("object-lock-demo-%s", time.Now().Format("20060102-150405"))
	enabled := true
	bucket := &models.Bucket{
		Name:             bucketName,
		Region:           "us-east-1",
		EnableVersioning: &enabled, // S3 Object Lock requires versioning
		S3ObjectLock: &models.BucketS3ObjectLockSettings{
			Enabled: &enabled,
			DefaultRetentionSetting: &models.BucketS3ObjectLockDefaultRetentionSettings{
				Mode: "compliance",
				Days: 30,
			},
		},
	}

	fmt.Printf("\n🏗️  Creating bucket %q with S3 Object Lock enabled...\n", bucketName)
	created, err := tenantClient.Bucket().Create(ctx, bucket)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}
	fmt.Printf("  ✅ Created %s in %s\n", created.Name, created.Region)

	// 5. Read back per-bucket Object Lock settings via the dedicated sub-resource.
	fmt.Println("\n🔍 Reading bucket Object Lock settings (/org/containers/<name>/object-lock)...")
	settings, err := tenantClient.Bucket().GetObjectLock(ctx, bucketName)
	if err != nil {
		log.Fatalf("Failed to get bucket Object Lock settings: %v", err)
	}
	printBucketObjectLock(settings)
}

func buildOpts(endpoint, user, pass string, accountID *string, skipSSL bool) []client.ClientOption {
	opts := []client.ClientOption{
		client.WithEndpoint(endpoint),
		client.WithCredentials(&models.Credentials{
			Username:  user,
			Password:  pass,
			AccountId: accountID,
		}),
	}
	if skipSSL {
		opts = append(opts, client.WithSkipSSL())
	}
	return opts
}

func printGridSettings(label string, s *models.S3ObjectLock) {
	fmt.Printf("\n📊 %s:\n", label)
	fmt.Printf("  complianceEnabled:             %s\n", boolStr(s.ComplianceEnabled))
	fmt.Printf("  legacyComplianceEnabled:       %s\n", boolStr(s.LegacyComplianceEnabled))
	fmt.Printf("  createLegacyComplianceBuckets: %s\n", boolStr(s.CreateLegacyComplianceBuckets))
}

func printBucketObjectLock(s *models.BucketS3ObjectLockSettings) {
	fmt.Printf("  enabled: %s\n", boolStr(s.Enabled))
	if s.DefaultRetentionSetting == nil {
		fmt.Println("  defaultRetentionSetting: (none)")
		return
	}
	fmt.Println("  defaultRetentionSetting:")
	fmt.Printf("    mode:  %s\n", s.DefaultRetentionSetting.Mode)
	if s.DefaultRetentionSetting.Days != 0 {
		fmt.Printf("    days:  %d\n", s.DefaultRetentionSetting.Days)
	}
	if s.DefaultRetentionSetting.Years != 0 {
		fmt.Printf("    years: %d\n", s.DefaultRetentionSetting.Years)
	}
}

func boolStr(b *bool) string {
	if b == nil {
		return "(unset)"
	}
	if *b {
		return "true"
	}
	return "false"
}

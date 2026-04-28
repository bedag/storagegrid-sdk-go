# S3 Object Lock Example (end-to-end)

This example walks through the full S3 Object Lock flow against a StorageGRID
deployment:

1. Reads the grid-wide compliance-global settings via the **GridClient**
   (`/grid/compliance-global`).
2. Optionally enables S3 Object Lock grid-wide (set `STORAGEGRID_APPLY=true`).
3. Updates the target tenant's policy (`allowComplianceMode`,
   `maxRetentionYears`) so the tenant can use S3 Object Lock and bucket
   retention is capped grid-side.
4. Creates a bucket via the **TenantClient** with S3 Object Lock enabled and a
   default compliance-mode retention period.
5. Reads the per-bucket Object Lock settings back via the dedicated
   `/org/containers/{bucketName}/object-lock` sub-resource.

> ⚠️ **Enabling S3 Object Lock grid-wide is irreversible.** See the
> [official documentation](https://docs.netapp.com/us-en/storagegrid/ilm/managing-objects-with-s3-object-lock.html#what-is-s3-object-lock).
 and tenant policy
  update)
- Tenant administrator credentials and an account ID (for the bucket stepsdentials (for the grid-wide settings)
- Tenant administrator credentials and an account ID (for the bucket steps)
- The tenant must have S3 Object Lock allowed via its tenant policy
  (`allowComplianceMode`)

## Environment Variables

```bash
# Required for the grid-side steps
export STORAGEGRID_ENDPOINT="https://your-storagegrid.example.com"
export STORAGEGRID_USERNAME="grid-admin"
export STORAGEGRID_PASSWORD="grid-password"

# Required for the tenant/bucket steps
export STORAGEGRID_TENANT_USERNAME="tenant-admin"
export STORAGEGRID_TENANT_PASSWORD="tenant-password"
export STORAGEGRID_ACCOUNT_ID="12345678901234567890"

# Optional
export STORAGEGRID_SKIP_SSL="true"   # development only
export STORAGEGRID_APPLY="true"      # actually enable S3 Object Lock grid-wide
```

If `STORAGEGRID_APPLY` is not set, the example only reads grid-wide settings
and does not create any buckets.

## Running

```bash
cd examples/grid/s3-object-lock
go mod init s3-object-lock-example
go mod tidy
go run main.go
```

## What you should see

```
🔍 Grid-wide S3 Object Lock settings (/grid/compliance-global)

📊 Current:
  complianceEnabled:             true
  legacyComplianceEnabled:       false
  createLegacyComplianceBuckets: false

🏗️  Creating bucket "object-lock-demo-20260428-101530" with S3 Object Lock enabled...
🔐 Updating tenant policy to allow S3 Object Lock...
  ✅ allowComplianceMode=true, maxRetentionYears=10 on tenant 12345678901234567890

  ✅ Created object-lock-demo-20260428-101530 in us-east-1

🔍 Reading bucket Object Lock settings (/org/containers/<name>/object-lock)...
  enabled: true
  defaultRetentionSetting:
    mode:  compliance
    days:  30
```

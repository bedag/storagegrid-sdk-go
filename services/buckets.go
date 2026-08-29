package services

import (
	"context"
	"fmt"

	"github.com/bedag/storagegrid-sdk-go/models"
)

const (
	bucketEndpoint      string = "/org/containers"
	tenantUsageEndpoint string = "/org/usage"
)

// BucketServiceInterface defines the contract for bucket service operations
type BucketServiceInterface interface {
	List(ctx context.Context) (*[]models.Bucket, error)
	GetByName(ctx context.Context, name string) (*models.Bucket, error)
	Create(ctx context.Context, bucket *models.Bucket) (*models.Bucket, error)
	// Deprecated: GetUsage retrieves the full tenant usage and filters by bucket
	// name. Use GetBucketUsage instead, which targets the per-bucket usage
	// endpoint (/org/containers/{name}/usage).
	GetUsage(ctx context.Context, name string) (*models.BucketStats, error)
	Delete(ctx context.Context, name string) error
	Drain(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error)
	CancelDrain(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error)
	DrainStatus(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error)

	// Per-bucket sub-resources

	GetBucketUsage(ctx context.Context, name string) (*models.BucketUsage, error)
	GetRegion(ctx context.Context, name string) (string, error)

	GetObjectLock(ctx context.Context, name string) (*models.BucketS3ObjectLockSettings, error)
	UpdateObjectLock(ctx context.Context, name string, settings *models.BucketS3ObjectLockSettings) (*models.BucketS3ObjectLockSettings, error)

	GetConsistency(ctx context.Context, name string) (*models.BucketConsistencySetting, error)
	UpdateConsistency(ctx context.Context, name string, settings *models.BucketConsistencySetting) (*models.BucketConsistencySetting, error)

	GetNotification(ctx context.Context, name string) (*models.BucketNotificationConfiguration, error)
	UpdateNotification(ctx context.Context, name string, config *models.BucketNotificationConfiguration) (*models.BucketNotificationConfiguration, error)

	GetPolicy(ctx context.Context, name string) (*models.BucketPolicyConfiguration, error)
	UpdatePolicy(ctx context.Context, name string, policy *models.BucketPolicyConfiguration) (*models.BucketPolicyConfiguration, error)

	GetCors(ctx context.Context, name string) (*models.BucketCorsConfiguration, error)
	UpdateCors(ctx context.Context, name string, config *models.BucketCorsConfiguration) (*models.BucketCorsConfiguration, error)

	GetCompliance(ctx context.Context, name string) (*models.BucketComplianceSettings, error)
	UpdateCompliance(ctx context.Context, name string, settings *models.BucketComplianceSettings) (*models.BucketComplianceSettings, error)
}

type BucketService struct {
	client HTTPClient
}

func NewBucketService(client HTTPClient) *BucketService {
	return &BucketService{client: client}
}

func (s *BucketService) List(ctx context.Context) (*[]models.Bucket, error) {
	response := models.Response{}
	response.Data = &[]models.Bucket{}
	err := s.client.DoParsed(ctx, "GET", bucketEndpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	buckets := response.Data.(*[]models.Bucket)

	return buckets, nil
}

func (s *BucketService) GetByName(ctx context.Context, name string) (*models.Bucket, error) {
	// the bucket endpoint doesn't have a simple get by name, so we have to list all buckets and find the one we want
	buckets, err := s.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, bucket := range *buckets {
		if bucket.Name == name {
			return &bucket, nil
		}
	}

	return nil, fmt.Errorf("bucket with name %s not found", name)
}

func (s *BucketService) Create(ctx context.Context, bucket *models.Bucket) (*models.Bucket, error) {
	response := models.Response{}
	response.Data = &models.Bucket{}
	err := s.client.DoParsed(ctx, "POST", bucketEndpoint, bucket, &response)
	if err != nil {
		return nil, err
	}

	bucket = response.Data.(*models.Bucket)

	return bucket, nil
}

// Deprecated: prefer GetBucketUsage which targets the per-bucket usage endpoint
// (/org/containers/{name}/usage) instead of fetching and filtering the full
// tenant usage payload.
func (s *BucketService) GetUsage(ctx context.Context, name string) (*models.BucketStats, error) {
	response := models.Response{}
	response.Data = &models.TenantUsage{}
	err := s.client.DoParsed(ctx, "GET", tenantUsageEndpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	tenantUsage := response.Data.(*models.TenantUsage)

	for _, bucket := range tenantUsage.Buckets {
		if *bucket.Name == name {
			return bucket, nil
		}
	}

	return nil, fmt.Errorf("usage for bucket with name %s not found", name)
}

func (s *BucketService) Delete(ctx context.Context, name string) error {
	err := s.client.DoParsed(ctx, "DELETE", bucketEndpoint+"/"+name, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// Drain a bucket by name. This will delete all objects in the bucket but leave the bucket itself intact.
//
// deleteObjects must be sent as a JSON boolean. Sending it as a string makes the API
// reject the request with 422 "Deleteobjects is not included in the list."
func (s *BucketService) Drain(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error) {
	response := models.Response{}
	response.Data = &models.BucketDeleteObjectStatus{}
	body := map[string]bool{"deleteObjects": true}

	err := s.client.DoParsed(ctx, "POST", bucketEndpoint+"/"+name+"/delete-objects", body, &response)
	if err != nil {
		return nil, err
	}

	deleteObjectStatus := response.Data.(*models.BucketDeleteObjectStatus)

	return deleteObjectStatus, nil
}

// CancelDrain stops an in-progress drain for the named bucket.
//
// deleteObjects must be sent as a JSON boolean; see Drain.
func (s *BucketService) CancelDrain(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error) {
	response := models.Response{}
	response.Data = &models.BucketDeleteObjectStatus{}
	body := map[string]bool{"deleteObjects": false}

	err := s.client.DoParsed(ctx, "POST", bucketEndpoint+"/"+name+"/delete-objects", body, &response)
	if err != nil {
		return nil, err
	}

	deleteObjectStatus := response.Data.(*models.BucketDeleteObjectStatus)

	return deleteObjectStatus, nil
}

func (s *BucketService) DrainStatus(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error) {
	response := models.Response{}
	response.Data = &models.BucketDeleteObjectStatus{}

	err := s.client.DoParsed(ctx, "GET", bucketEndpoint+"/"+name+"/delete-objects", nil, &response)
	if err != nil {
		return nil, err
	}

	deleteObjectStatus := response.Data.(*models.BucketDeleteObjectStatus)

	return deleteObjectStatus, nil
}

// bucketSubresource builds the URL path for a per-bucket sub-resource endpoint.
func bucketSubresource(name, subresource string) string {
	return bucketEndpoint + "/" + name + "/" + subresource
}

// GetBucketUsage retrieves the per-bucket usage metrics from
// GET /org/containers/{name}/usage.
func (s *BucketService) GetBucketUsage(ctx context.Context, name string) (*models.BucketUsage, error) {
	response := models.Response{}
	response.Data = &models.BucketUsage{}
	if err := s.client.DoParsed(ctx, "GET", bucketSubresource(name, "usage"), nil, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketUsage), nil
}

// GetRegion retrieves the region for a bucket from GET /org/containers/{name}/region.
func (s *BucketService) GetRegion(ctx context.Context, name string) (string, error) {
	response := models.Response{}
	data := struct {
		Region string `json:"region"`
	}{}
	response.Data = &data
	if err := s.client.DoParsed(ctx, "GET", bucketSubresource(name, "region"), nil, &response); err != nil {
		return "", err
	}
	return data.Region, nil
}

// GetObjectLock retrieves the S3 Object Lock settings for a bucket.
func (s *BucketService) GetObjectLock(ctx context.Context, name string) (*models.BucketS3ObjectLockSettings, error) {
	response := models.Response{}
	response.Data = &models.BucketS3ObjectLockSettings{}
	if err := s.client.DoParsed(ctx, "GET", bucketSubresource(name, "object-lock"), nil, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketS3ObjectLockSettings), nil
}

// UpdateObjectLock updates the S3 Object Lock settings for a bucket.
func (s *BucketService) UpdateObjectLock(ctx context.Context, name string, settings *models.BucketS3ObjectLockSettings) (*models.BucketS3ObjectLockSettings, error) {
	response := models.Response{}
	response.Data = &models.BucketS3ObjectLockSettings{}
	if err := s.client.DoParsed(ctx, "PUT", bucketSubresource(name, "object-lock"), settings, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketS3ObjectLockSettings), nil
}

// GetConsistency retrieves the consistency value for a bucket from
// GET /org/containers/{name}/consistency.
func (s *BucketService) GetConsistency(ctx context.Context, name string) (*models.BucketConsistencySetting, error) {
	response := models.Response{}
	response.Data = &models.BucketConsistencySetting{}
	if err := s.client.DoParsed(ctx, "GET", bucketSubresource(name, "consistency"), nil, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketConsistencySetting), nil
}

// UpdateConsistency sets the consistency value for a bucket. StorageGRID applies the change
// only to objects ingested after it; objects already in the bucket keep their prior behavior.
//
// The reducedConsistency query flag the endpoint also accepts is deliberately not exposed:
// NetApp documents it as for use only when directed by technical support.
func (s *BucketService) UpdateConsistency(ctx context.Context, name string, settings *models.BucketConsistencySetting) (*models.BucketConsistencySetting, error) {
	response := models.Response{}
	response.Data = &models.BucketConsistencySetting{}
	if err := s.client.DoParsed(ctx, "PUT", bucketSubresource(name, "consistency"), settings, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketConsistencySetting), nil
}

// GetNotification retrieves the notification configuration for a bucket.
func (s *BucketService) GetNotification(ctx context.Context, name string) (*models.BucketNotificationConfiguration, error) {
	response := models.Response{}
	response.Data = &models.BucketNotificationConfiguration{}
	if err := s.client.DoParsed(ctx, "GET", bucketSubresource(name, "notification"), nil, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketNotificationConfiguration), nil
}

// UpdateNotification updates the notification configuration for a bucket.
func (s *BucketService) UpdateNotification(ctx context.Context, name string, config *models.BucketNotificationConfiguration) (*models.BucketNotificationConfiguration, error) {
	response := models.Response{}
	response.Data = &models.BucketNotificationConfiguration{}
	if err := s.client.DoParsed(ctx, "PUT", bucketSubresource(name, "notification"), config, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketNotificationConfiguration), nil
}

// GetPolicy retrieves the bucket policy.
func (s *BucketService) GetPolicy(ctx context.Context, name string) (*models.BucketPolicyConfiguration, error) {
	response := models.Response{}
	response.Data = &models.BucketPolicyConfiguration{}
	if err := s.client.DoParsed(ctx, "GET", bucketSubresource(name, "policy"), nil, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketPolicyConfiguration), nil
}

// UpdatePolicy updates the bucket policy.
func (s *BucketService) UpdatePolicy(ctx context.Context, name string, policy *models.BucketPolicyConfiguration) (*models.BucketPolicyConfiguration, error) {
	response := models.Response{}
	response.Data = &models.BucketPolicyConfiguration{}
	if err := s.client.DoParsed(ctx, "PUT", bucketSubresource(name, "policy"), policy, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketPolicyConfiguration), nil
}

// GetCors retrieves the CORS configuration for a bucket.
func (s *BucketService) GetCors(ctx context.Context, name string) (*models.BucketCorsConfiguration, error) {
	response := models.Response{}
	response.Data = &models.BucketCorsConfiguration{}
	if err := s.client.DoParsed(ctx, "GET", bucketSubresource(name, "cors"), nil, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketCorsConfiguration), nil
}

// UpdateCors updates the CORS configuration for a bucket.
func (s *BucketService) UpdateCors(ctx context.Context, name string, config *models.BucketCorsConfiguration) (*models.BucketCorsConfiguration, error) {
	response := models.Response{}
	response.Data = &models.BucketCorsConfiguration{}
	if err := s.client.DoParsed(ctx, "PUT", bucketSubresource(name, "cors"), config, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketCorsConfiguration), nil
}

// GetCompliance retrieves the legacy Compliance settings for a bucket.
func (s *BucketService) GetCompliance(ctx context.Context, name string) (*models.BucketComplianceSettings, error) {
	response := models.Response{}
	response.Data = &models.BucketComplianceSettings{}
	if err := s.client.DoParsed(ctx, "GET", bucketSubresource(name, "compliance"), nil, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketComplianceSettings), nil
}

// UpdateCompliance updates the legacy Compliance settings for a bucket.
func (s *BucketService) UpdateCompliance(ctx context.Context, name string, settings *models.BucketComplianceSettings) (*models.BucketComplianceSettings, error) {
	response := models.Response{}
	response.Data = &models.BucketComplianceSettings{}
	if err := s.client.DoParsed(ctx, "PUT", bucketSubresource(name, "compliance"), settings, &response); err != nil {
		return nil, err
	}
	return response.Data.(*models.BucketComplianceSettings), nil
}

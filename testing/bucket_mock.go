package testing

import (
	"context"

	"github.com/bedag/storagegrid-sdk-go/models"
	"github.com/bedag/storagegrid-sdk-go/services"
)

// MockBucketService implements services.BucketServiceInterface for testing
type MockBucketService struct {
	ListFunc        func(ctx context.Context) (*[]models.Bucket, error)
	GetByNameFunc   func(ctx context.Context, name string) (*models.Bucket, error)
	CreateFunc      func(ctx context.Context, bucket *models.Bucket) (*models.Bucket, error)
	GetUsageFunc    func(ctx context.Context, name string) (*models.BucketStats, error)
	DeleteFunc      func(ctx context.Context, name string) error
	DrainFunc       func(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error)
	CancelDrainFunc func(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error)
	DrainStatusFunc func(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error)

	GetBucketUsageFunc     func(ctx context.Context, name string) (*models.BucketUsage, error)
	GetRegionFunc          func(ctx context.Context, name string) (string, error)
	GetObjectLockFunc      func(ctx context.Context, name string) (*models.BucketS3ObjectLockSettings, error)
	UpdateObjectLockFunc   func(ctx context.Context, name string, settings *models.BucketS3ObjectLockSettings) (*models.BucketS3ObjectLockSettings, error)
	GetConsistencyFunc     func(ctx context.Context, name string) (*models.BucketConsistencySetting, error)
	UpdateConsistencyFunc  func(ctx context.Context, name string, settings *models.BucketConsistencySetting) (*models.BucketConsistencySetting, error)
	GetNotificationFunc    func(ctx context.Context, name string) (*models.BucketNotificationConfiguration, error)
	UpdateNotificationFunc func(ctx context.Context, name string, config *models.BucketNotificationConfiguration) (*models.BucketNotificationConfiguration, error)
	GetPolicyFunc          func(ctx context.Context, name string) (*models.BucketPolicyConfiguration, error)
	UpdatePolicyFunc       func(ctx context.Context, name string, policy *models.BucketPolicyConfiguration) (*models.BucketPolicyConfiguration, error)
	GetCorsFunc            func(ctx context.Context, name string) (*models.BucketCorsConfiguration, error)
	UpdateCorsFunc         func(ctx context.Context, name string, config *models.BucketCorsConfiguration) (*models.BucketCorsConfiguration, error)
	GetComplianceFunc      func(ctx context.Context, name string) (*models.BucketComplianceSettings, error)
	UpdateComplianceFunc   func(ctx context.Context, name string, settings *models.BucketComplianceSettings) (*models.BucketComplianceSettings, error)
}

func (m *MockBucketService) List(ctx context.Context) (*[]models.Bucket, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return &[]models.Bucket{}, nil
}

func (m *MockBucketService) GetByName(ctx context.Context, name string) (*models.Bucket, error) {
	if m.GetByNameFunc != nil {
		return m.GetByNameFunc(ctx, name)
	}
	return &models.Bucket{Name: name}, nil
}

func (m *MockBucketService) Create(ctx context.Context, bucket *models.Bucket) (*models.Bucket, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, bucket)
	}
	return bucket, nil
}

func (m *MockBucketService) GetUsage(ctx context.Context, name string) (*models.BucketStats, error) {
	if m.GetUsageFunc != nil {
		return m.GetUsageFunc(ctx, name)
	}
	bucketName := name
	return &models.BucketStats{Name: &bucketName}, nil
}

func (m *MockBucketService) Delete(ctx context.Context, name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, name)
	}
	return nil
}

func (m *MockBucketService) Drain(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error) {
	if m.DrainFunc != nil {
		return m.DrainFunc(ctx, name)
	}
	return &models.BucketDeleteObjectStatus{}, nil
}

func (m *MockBucketService) CancelDrain(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error) {
	if m.DrainFunc != nil {
		return m.CancelDrainFunc(ctx, name)
	}
	return &models.BucketDeleteObjectStatus{}, nil
}

func (m *MockBucketService) DrainStatus(ctx context.Context, name string) (*models.BucketDeleteObjectStatus, error) {
	if m.DrainStatusFunc != nil {
		return m.DrainStatusFunc(ctx, name)
	}
	return &models.BucketDeleteObjectStatus{}, nil
}

func (m *MockBucketService) GetBucketUsage(ctx context.Context, name string) (*models.BucketUsage, error) {
	if m.GetBucketUsageFunc != nil {
		return m.GetBucketUsageFunc(ctx, name)
	}
	return &models.BucketUsage{}, nil
}

func (m *MockBucketService) GetRegion(ctx context.Context, name string) (string, error) {
	if m.GetRegionFunc != nil {
		return m.GetRegionFunc(ctx, name)
	}
	return "", nil
}

func (m *MockBucketService) GetObjectLock(ctx context.Context, name string) (*models.BucketS3ObjectLockSettings, error) {
	if m.GetObjectLockFunc != nil {
		return m.GetObjectLockFunc(ctx, name)
	}
	return &models.BucketS3ObjectLockSettings{}, nil
}

func (m *MockBucketService) UpdateObjectLock(ctx context.Context, name string, settings *models.BucketS3ObjectLockSettings) (*models.BucketS3ObjectLockSettings, error) {
	if m.UpdateObjectLockFunc != nil {
		return m.UpdateObjectLockFunc(ctx, name, settings)
	}
	return settings, nil
}

func (m *MockBucketService) GetConsistency(ctx context.Context, name string) (*models.BucketConsistencySetting, error) {
	if m.GetConsistencyFunc != nil {
		return m.GetConsistencyFunc(ctx, name)
	}
	return &models.BucketConsistencySetting{}, nil
}

func (m *MockBucketService) UpdateConsistency(ctx context.Context, name string, settings *models.BucketConsistencySetting) (*models.BucketConsistencySetting, error) {
	if m.UpdateConsistencyFunc != nil {
		return m.UpdateConsistencyFunc(ctx, name, settings)
	}
	return settings, nil
}

func (m *MockBucketService) GetNotification(ctx context.Context, name string) (*models.BucketNotificationConfiguration, error) {
	if m.GetNotificationFunc != nil {
		return m.GetNotificationFunc(ctx, name)
	}
	return &models.BucketNotificationConfiguration{}, nil
}

func (m *MockBucketService) UpdateNotification(ctx context.Context, name string, config *models.BucketNotificationConfiguration) (*models.BucketNotificationConfiguration, error) {
	if m.UpdateNotificationFunc != nil {
		return m.UpdateNotificationFunc(ctx, name, config)
	}
	return config, nil
}

func (m *MockBucketService) GetPolicy(ctx context.Context, name string) (*models.BucketPolicyConfiguration, error) {
	if m.GetPolicyFunc != nil {
		return m.GetPolicyFunc(ctx, name)
	}
	return &models.BucketPolicyConfiguration{}, nil
}

func (m *MockBucketService) UpdatePolicy(ctx context.Context, name string, policy *models.BucketPolicyConfiguration) (*models.BucketPolicyConfiguration, error) {
	if m.UpdatePolicyFunc != nil {
		return m.UpdatePolicyFunc(ctx, name, policy)
	}
	return policy, nil
}

func (m *MockBucketService) GetCors(ctx context.Context, name string) (*models.BucketCorsConfiguration, error) {
	if m.GetCorsFunc != nil {
		return m.GetCorsFunc(ctx, name)
	}
	return &models.BucketCorsConfiguration{}, nil
}

func (m *MockBucketService) UpdateCors(ctx context.Context, name string, config *models.BucketCorsConfiguration) (*models.BucketCorsConfiguration, error) {
	if m.UpdateCorsFunc != nil {
		return m.UpdateCorsFunc(ctx, name, config)
	}
	return config, nil
}

func (m *MockBucketService) GetCompliance(ctx context.Context, name string) (*models.BucketComplianceSettings, error) {
	if m.GetComplianceFunc != nil {
		return m.GetComplianceFunc(ctx, name)
	}
	return &models.BucketComplianceSettings{}, nil
}

func (m *MockBucketService) UpdateCompliance(ctx context.Context, name string, settings *models.BucketComplianceSettings) (*models.BucketComplianceSettings, error) {
	if m.UpdateComplianceFunc != nil {
		return m.UpdateComplianceFunc(ctx, name, settings)
	}
	return settings, nil
}

// Compile-time interface compliance check
var _ services.BucketServiceInterface = (*MockBucketService)(nil)

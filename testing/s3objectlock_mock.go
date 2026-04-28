package testing

import (
	"context"

	"github.com/bedag/storagegrid-sdk-go/models"
	"github.com/bedag/storagegrid-sdk-go/services"
)

// MockS3ObjectLockService implements services.S3ObjectLockServiceInterface for testing
type MockS3ObjectLockService struct {
	GetFunc    func(ctx context.Context) (*models.S3ObjectLock, error)
	UpdateFunc func(ctx context.Context, settings *models.S3ObjectLock) (*models.S3ObjectLock, error)
}

func (m *MockS3ObjectLockService) Get(ctx context.Context) (*models.S3ObjectLock, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx)
	}
	return &models.S3ObjectLock{}, nil
}

func (m *MockS3ObjectLockService) Update(ctx context.Context, settings *models.S3ObjectLock) (*models.S3ObjectLock, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, settings)
	}
	return settings, nil
}

// Compile-time interface compliance check
var _ services.S3ObjectLockServiceInterface = (*MockS3ObjectLockService)(nil)

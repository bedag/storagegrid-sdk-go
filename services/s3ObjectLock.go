package services

import (
	"context"

	"github.com/bedag/storagegrid-sdk-go/models"
)

const (
	s3ObjectLockEndpoint string = "/grid/compliance-global"
)

// S3ObjectLockServiceInterface defines the contract for grid-wide S3 Object Lock
// (compliance-global) operations.
type S3ObjectLockServiceInterface interface {
	Get(ctx context.Context) (*models.S3ObjectLock, error)
	Update(ctx context.Context, settings *models.S3ObjectLock) (*models.S3ObjectLock, error)
}

type S3ObjectLockService struct {
	client HTTPClient
}

func NewS3ObjectLockService(client HTTPClient) *S3ObjectLockService {
	return &S3ObjectLockService{client: client}
}

func (s *S3ObjectLockService) Get(ctx context.Context) (*models.S3ObjectLock, error) {
	response := models.Response{}
	response.Data = &models.S3ObjectLock{}
	err := s.client.DoParsed(ctx, "GET", s3ObjectLockEndpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	settings := response.Data.(*models.S3ObjectLock)

	return settings, nil
}

func (s *S3ObjectLockService) Update(ctx context.Context, settings *models.S3ObjectLock) (*models.S3ObjectLock, error) {
	response := models.Response{}
	response.Data = &models.S3ObjectLock{}
	err := s.client.DoParsed(ctx, "PUT", s3ObjectLockEndpoint, settings, &response)
	if err != nil {
		return nil, err
	}

	updated := response.Data.(*models.S3ObjectLock)

	return updated, nil
}

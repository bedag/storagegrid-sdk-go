package testing

import (
	"context"

	"github.com/bedag/storagegrid-sdk-go/models"
)

// MockMetricsService implements services.MetricsServiceInterface for testing
type MockMetricsService struct {
	QueryFunc func(ctx context.Context, promql string) (*models.MetricQueryResult, error)
}

func (m *MockMetricsService) Query(ctx context.Context, promql string) (*models.MetricQueryResult, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, promql)
	}
	return &models.MetricQueryResult{}, nil
}

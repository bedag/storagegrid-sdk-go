package services

import (
	"context"
	"net/url"

	"github.com/bedag/storagegrid-sdk-go/models"
)

const metricQueryEndpoint string = "/grid/metric-query"

// MetricsServiceInterface defines the contract for grid metric queries.
type MetricsServiceInterface interface {
	Query(ctx context.Context, promql string) (*models.MetricQueryResult, error)
}

type MetricsService struct {
	client HTTPClient
}

func NewMetricsService(client HTTPClient) *MetricsService {
	return &MetricsService{client: client}
}

func (s *MetricsService) Query(ctx context.Context, promql string) (*models.MetricQueryResult, error) {
	response := models.Response{}
	response.Data = &models.MetricQueryResult{}

	err := s.client.DoParsed(ctx, "GET", metricQueryEndpoint+"?query="+url.QueryEscape(promql), nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Data.(*models.MetricQueryResult), nil
}

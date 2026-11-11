package models

import (
	"fmt"
	"strconv"
)

// MetricQueryResult is the data section of a /grid/metric-query response (Prometheus instant query).
type MetricQueryResult struct {
	ResultType string         `json:"resultType"`
	Result     []MetricSample `json:"result"`
}

// MetricSample is one series of an instant query result.
type MetricSample struct {
	Metric map[string]string `json:"metric"`
	Value  [2]any            `json:"value"` // [unix timestamp, value as string]
}

// Float parses the sample value.
func (s MetricSample) Float() (float64, error) {
	str, ok := s.Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("metric sample value is %T, expected string", s.Value[1])
	}
	return strconv.ParseFloat(str, 64)
}

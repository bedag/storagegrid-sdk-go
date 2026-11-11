package services

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func TestMetricQueryEncodesAndParses(t *testing.T) {
	const q = "min by (endpoint_id) (storagegrid_load_balancer_endpoint_bound)"
	f := &fakeHTTP{body: `{"data":{"resultType":"vector","result":[{"metric":{"__name__":"x","endpoint_id":"c47f902c"},"value":[1762800000.123,"1"]}]}}`}
	got, err := NewMetricsService(f).Query(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	raw, ok := strings.CutPrefix(f.path, "/grid/metric-query?query=")
	if !ok {
		t.Fatalf("path = %q", f.path)
	}
	if dec, err := url.QueryUnescape(raw); err != nil || dec != q {
		t.Fatalf("decoded query = %q, %v", dec, err)
	}
	s := got.Result[0]
	if s.Metric["endpoint_id"] != "c47f902c" {
		t.Fatalf("metric = %v", s.Metric)
	}
	if v, err := s.Float(); err != nil || v != 1 {
		t.Fatalf("Float() = %v, %v", v, err)
	}
}

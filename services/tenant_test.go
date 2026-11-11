package services

import (
	"context"
	"testing"
)

func TestListUsageCallsAccountsCache(t *testing.T) {
	f := &fakeHTTP{body: `{"data":[{"id":"7173","name":"t","capabilities":["s3"],"policy":{"quotaObjectBytes":10737418240},"objectCount":10,"dataBytes":1569128713}]}`}
	got, err := NewTenantService(f).ListUsage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if f.path != "/grid/accounts-cache" {
		t.Fatalf("path = %q", f.path)
	}
	e := (*got)[0]
	if e.Id != "7173" || *e.DataBytes != 1569128713 || *e.ObjectCount != 10 || *e.Policy.QuotaObjectBytes != 10737418240 {
		t.Fatalf("unexpected entry %+v", e)
	}
}

func TestGetUsageWithBucketsQuery(t *testing.T) {
	f := &fakeHTTP{body: `{"data":{"calculationTime":"2026-11-10T18:58:25.707000Z","objectCount":10,"dataBytes":1569128713,"buckets":[{"name":"flatcar-image-store-dev-2e0af0b5","objectCount":10,"dataBytes":1569128713,"region":"us-east-1"}]}}`}
	got, err := NewTenantService(f).GetUsageWithBuckets(context.Background(), "7173")
	if err != nil {
		t.Fatal(err)
	}
	if f.path != "/grid/accounts/7173/usage?includeBucketDetail=true" {
		t.Fatalf("path = %q", f.path)
	}
	if *got.Buckets[0].Name != "flatcar-image-store-dev-2e0af0b5" {
		t.Fatalf("unexpected buckets %+v", got.Buckets)
	}
}

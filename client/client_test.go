package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type countingTransport struct {
	calls atomic.Int32
	next  http.RoundTripper
}

func (t *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.calls.Add(1)
	return t.next.RoundTrip(req)
}

func TestNewClientDoesNotTouchDefaultClient(t *testing.T) {
	before := http.DefaultClient.Transport

	c1, err := newClient(WithEndpoint("grid.example.com"), WithSkipSSL())
	if err != nil {
		t.Fatal(err)
	}
	c2, err := newClient(WithEndpoint("grid.example.com"))
	if err != nil {
		t.Fatal(err)
	}

	if http.DefaultClient.Transport != before {
		t.Fatal("http.DefaultClient.Transport was modified")
	}
	if c1.httpClient == c2.httpClient || c1.httpClient == http.DefaultClient {
		t.Fatal("clients must not share an http.Client")
	}
	if !c1.httpClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify {
		t.Error("WithSkipSSL client should skip TLS verification")
	}
	if c2.httpClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify {
		t.Error("client without WithSkipSSL must verify TLS")
	}
}

func TestWithHTTPClientIsUsed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	rt := &countingTransport{next: http.DefaultTransport}
	hc := &http.Client{Transport: rt}
	c, err := newClient(WithEndpoint(srv.URL), WithHTTPClient(hc))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.DoUnparsed(context.Background(), http.MethodPost, "/authorize", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if got := rt.calls.Load(); got != 1 {
		t.Errorf("expected 1 call through the provided transport, got %d", got)
	}
	if hc.Transport != rt {
		t.Error("provided client's Transport must not be replaced")
	}
}

func TestWithHTTPClientRejectsSkipSSL(t *testing.T) {
	if _, err := newClient(WithEndpoint("grid.example.com"), WithHTTPClient(&http.Client{}), WithSkipSSL()); err == nil {
		t.Fatal("expected an error when combining WithHTTPClient and WithSkipSSL")
	}
}

func TestRequestHonorsContext(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c, err := newClient(WithEndpoint(srv.URL))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := c.DoUnparsed(ctx, http.MethodPost, "/authorize", nil); err == nil {
		t.Fatal("expected an error for a canceled context")
	}
	if hits.Load() != 0 {
		t.Error("request must not reach the server after the context was canceled")
	}
}

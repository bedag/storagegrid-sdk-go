package services

import (
	"context"
	"encoding/json"
	"net/http"
)

// fakeHTTP records the requested path and decodes canned JSON into output.
type fakeHTTP struct {
	path string
	body string
}

func (f *fakeHTTP) DoParsed(_ context.Context, _, path string, _, output interface{}) error {
	f.path = path
	return json.Unmarshal([]byte(f.body), output)
}

func (f *fakeHTTP) DoUnparsed(context.Context, string, string, interface{}) (*http.Response, error) {
	return nil, nil
}

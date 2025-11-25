package e2e

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makePostCreateRequest(t *testing.T, ts *httptest.Server, endPoint, req string) {
	resp, err := http.Post(ts.URL+endPoint, "application/json", bytes.NewBufferString(req))
	if err != nil {
		t.Fatalf("failed POST %s: %v", endPoint, err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	if err = resp.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}
}

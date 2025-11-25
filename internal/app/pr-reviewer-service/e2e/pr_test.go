package e2e

import (
	"bytes"
	"net/http"
	"testing"
)

func TestPRLifecycle(t *testing.T) {
	_, ts := StartTestServer(t)
	defer ts.Close()

	team := `{
		"team_name": "backend",
		"members": [
			{"user_id":"u1", "username":"Author", "is_active":true},
			{"user_id":"u2", "username":"Reviewer", "is_active":true}
		]
	}`
	makePostCreateRequest(t, ts, "/team/add", team)

	req := `{
		"pull_request_id":"pr1",
		"pull_request_name":"Fix login",
		"author_id":"u1"
	}`
	makePostCreateRequest(t, ts, "/pullRequest/create", req)

	req = `{"pull_request_id":"pr1"}`
	resp, err := http.Post(ts.URL+"/pullRequest/merge", "application/json", bytes.NewBufferString(req))
	if err != nil {
		t.Fatalf("failed POST %s: %v", "/pullRequest/merge", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("merge status: %d", resp.StatusCode)
	}
	if err = resp.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}
}

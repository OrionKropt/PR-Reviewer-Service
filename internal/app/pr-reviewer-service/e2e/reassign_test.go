package e2e

import (
	"bytes"
	"net/http"
	"testing"
)

func TestPRReassign(t *testing.T) {
	_, ts := StartTestServer(t)
	defer ts.Close()

	team := `{
		"team_name": "backend",
		"members": [
			{"user_id":"u1","username":"Alice","is_active":true},
			{"user_id":"u2","username":"Bob","is_active":true},
			{"user_id":"u3","username":"Charlie","is_active":true},
			{"user_id":"u4","username":"David","is_active":true}
		]
	}`
	makePostCreateRequest(t, ts, "/team/add", team)

	prCreate := `{
		"pull_request_id": "42",
		"pull_request_name": "Refactor",
		"author_id": "u1"
	}`
	makePostCreateRequest(t, ts, "/pullRequest/create", prCreate)

	reassignReq := `{
		"pull_request_id":"42",
		"old_user_id":"u2"
	}`
	resp, err := http.Post(ts.URL+"/pullRequest/reassign", "application/json", bytes.NewBufferString(reassignReq))
	if err != nil {
		t.Fatalf("failed POST %s: %v", "pullRequest/reassign", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if err = resp.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}
}

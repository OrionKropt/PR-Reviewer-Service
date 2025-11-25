package e2e

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/OrionKropt/PRReviewerService/api/types"
)

func TestTeamAddAndGet(t *testing.T) {
	_, ts := StartTestServer(t)
	defer ts.Close()

	team := `{
			"team_name": "backend",
			"members": [
				{"user_id": "u1", "username": "Alice", "is_active": true},
				{"user_id": "u2", "username": "Bob", "is_active": true}
			]
		}`

	// === POST /team/add ===
	makePostCreateRequest(t, ts, "/team/add", team)

	// === GET /team/get?team_name=backend ===
	resp, err := http.Get(ts.URL + "/team/get?team_name=backend")
	if err != nil {
		t.Fatalf("failed GET /team/get: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	var got types.Team
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if err = resp.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}

	if got.TeamName != "backend" {
		t.Fatalf("expected team backend, got %s", got.TeamName)
	}
	if len(got.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(got.Members))
	}
	if err = resp.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}
}

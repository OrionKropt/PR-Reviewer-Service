package types

const (
	TeamExists  = "TEAM_EXISTS"
	PRExists    = "PR_EXISTS"
	PRMerged    = "PR_MERGED"
	NotAssigned = "NOT_ASSIGNED"
	NoCandidate = "NO_CANDIDATE"
	NotFound    = "NOT_FOUND"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

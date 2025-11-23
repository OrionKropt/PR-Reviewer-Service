package types

type ErrorCode string

const (
	TeamExists    ErrorCode = "TEAM_EXISTS"
	PRExists      ErrorCode = "PR_EXISTS"
	PRMerged      ErrorCode = "PR_MERGED"
	NotAssigned   ErrorCode = "NOT_ASSIGNED"
	NoCandidate   ErrorCode = "NO_CANDIDATE"
	NotFound      ErrorCode = "NOT_FOUND"
	BadRequest    ErrorCode = "BAD_REQUEST"
	InternalError ErrorCode = "INTERNAL"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

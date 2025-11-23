package domain

import "time"

const (
	PROpen   = "OPEN"
	PRMerged = "MERGED"
)

type PullRequest struct {
	id                string
	Name              string
	AuthorID          string
	AssignedReviewers []string
	Status            string
	CreatedAt         string
	MergedAt          string
}

func CreatePullRequest(id, authorID string, name string) PullRequest {
	return PullRequest{
		id:                id,
		Name:              name,
		AuthorID:          authorID,
		AssignedReviewers: make([]string, 0),
		Status:            PROpen,
		CreatedAt:         time.Now().Format(time.RFC3339),
	}
}

func (p *PullRequest) ID() string {
	return p.id
}

package domain

const (
	PROpen = iota
	PRMerged
)

type PullRequest struct {
	id               string
	Name             string
	AuthorID         string
	AssignedReviewer []string
	Status           int
}

func CreatePullRequest(id, authorID string, name string) PullRequest {
	return PullRequest{
		id:               id,
		Name:             name,
		AuthorID:         authorID,
		AssignedReviewer: make([]string, 0),
		Status:           PROpen,
	}
}

func (p *PullRequest) ID() string {
	return p.id
}

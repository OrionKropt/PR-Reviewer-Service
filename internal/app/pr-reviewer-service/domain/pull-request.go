package domain

const (
	PROpen               = "OPEN"
	PRMerged             = "MERGED"
	MaxAssignedReviewers = 2
)

type PullRequest struct {
	id                string
	Name              string
	AuthorID          string
	AssignedReviewers map[string]*User
	Status            string
	CreatedAt         string
	MergedAt          string
}

func CreatePullRequest(id, name, authorID, createAt string) PullRequest {
	return PullRequest{
		id:                id,
		Name:              name,
		AuthorID:          authorID,
		AssignedReviewers: make(map[string]*User),
		Status:            PROpen,
		CreatedAt:         createAt,
	}
}

func (p *PullRequest) ID() string {
	return p.id
}

func (pr *PullRequest) ReassignReviewer(oldUserID string, newReviewer *User) bool {
	if _, ok := pr.AssignedReviewers[oldUserID]; ok {
		pr.AssignedReviewers[newReviewer.ID()] = newReviewer
		delete(pr.AssignedReviewers, oldUserID)
		return true
	}
	return false
}

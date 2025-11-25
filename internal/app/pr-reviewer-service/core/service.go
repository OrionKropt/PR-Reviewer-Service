package core

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/OrionKropt/PRReviewerService/api/types"
	dom "github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/domain"
)

type PRReviewerService struct {
	mu sync.RWMutex

	teams map[string]*dom.Team
	users map[string]*dom.User
	prs   map[string]*dom.PullRequest
}

func NewPRReviewerService() *PRReviewerService {
	return &PRReviewerService{
		teams: make(map[string]*dom.Team),
		users: make(map[string]*dom.User),
		prs:   make(map[string]*dom.PullRequest),
	}
}

func (s *PRReviewerService) CreateTeam(name string, members []types.TeamMember) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.teams[name]; ok {
		return dom.ErrTeamAlreadyExists
	}
	newTeam := dom.CreateTeam(name)
	for _, m := range members {
		newUser := dom.CreateUser(m.UserID, m.Username, name, m.IsActive)
		newTeam.AddMember(&newUser)
		s.users[m.UserID] = &newUser
	}
	s.teams[name] = &newTeam
	return nil
}

func (s *PRReviewerService) GetTeam(name string) (*types.Team, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	team, ok := s.teams[name]
	if !ok {
		return nil, dom.ErrTeamNotFound
	}
	members := team.Members
	outTeam := &types.Team{
		TeamName: team.Name,
		Members:  make([]types.TeamMember, 0),
	}
	for _, m := range members {
		outTeam.Members = append(outTeam.Members, types.TeamMember{
			Username: m.Name,
			UserID:   m.ID(),
			IsActive: m.IsActive,
		})
	}
	return outTeam, nil
}

func (s *PRReviewerService) GetPullRequestsAsReviewer(userID string) ([]types.PullRequestShort, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return nil, fmt.Errorf("user: %s, error: %w", userID, dom.ErrUserNotFound)
	}
	team, ok := s.teams[user.TeamName]
	if !ok {
		return nil, fmt.Errorf("team: %s, error: %w", user.TeamName, dom.ErrTeamNotFound)
	}

	prAsReviewer := make([]types.PullRequestShort, 0, len(team.PullRequests))
	for _, pr := range team.PullRequests {
		if _, ok = pr.AssignedReviewers[userID]; ok {
			prAsReviewer = append(prAsReviewer, types.PullRequestShort{
				PullRequestID:   pr.ID(),
				PullRequestName: pr.Name,
				AuthorID:        pr.AuthorID,
				Status:          pr.Status,
			})
		}
	}
	return prAsReviewer, nil
}

func (s *PRReviewerService) SetIsActive(userId string, isActive bool) (*types.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[userId]
	if !ok {
		return nil, fmt.Errorf("user: %s, error: %w", userId, dom.ErrUserNotFound)
	}
	user.IsActive = isActive
	outUser := &types.User{
		UserID:   user.ID(),
		Username: user.Name,
		TeamName: user.TeamName,
		IsActive: isActive,
	}
	return outUser, nil
}

func (s *PRReviewerService) CreatePullRequest(id, name, authorID string) (*types.PullRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.prs[id]; ok {
		return nil, dom.ErrPRExists
	}
	author, ok := s.users[authorID]
	if !ok {
		return nil, dom.ErrAuthorNotFound
	}
	team, ok := s.teams[author.TeamName]
	if !ok {
		return nil, dom.ErrTeamNotFound
	}
	newPR := dom.CreatePullRequest(id, name, authorID, currentTime())
	possibleReviewers, err := getPossibleReviewers(team, &newPR, author)
	reviewers := make([]string, 0, dom.MaxAssignedReviewers)
	if err == nil {
		for i := 0; i < len(possibleReviewers) && i < dom.MaxAssignedReviewers; i++ {
			reviewers = append(reviewers, possibleReviewers[i].ID())
			newPR.AssignedReviewers[possibleReviewers[i].ID()] = possibleReviewers[i]
		}
	}

	s.prs[id] = &newPR
	team.PullRequests[id] = &newPR

	outPR := types.CreatePullRequest(id, name, authorID, newPR.Status, newPR.CreatedAt, newPR.MergedAt, reviewers)

	return outPR, nil
}

func (s *PRReviewerService) MergePullRequest(id string) (*types.PullRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pr, ok := s.prs[id]
	if !ok {
		return nil, fmt.Errorf("pr: %s, error: %w", id, dom.ErrPRNotFound)
	}
	if pr.Status != dom.PRMerged {
		pr.Status = dom.PRMerged
		pr.MergedAt = currentTime()
	}
	reviewers := make([]string, 0)
	for _, r := range pr.AssignedReviewers {
		reviewers = append(reviewers, r.ID())
	}
	outPR := types.CreatePullRequest(pr.ID(), pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt, reviewers)
	return outPR, nil
}

func (s *PRReviewerService) ReassignReviewerPullRequest(id, oldReviewerID string) (*types.PullRequest, string, error) {

	chooseRandomUser := func(users []*dom.User) *dom.User {
		// #nosec G404 - non-crypto random is acceptable for reviewer selection
		randomIndex := rand.Intn(len(users))
		return users[randomIndex]
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[oldReviewerID]
	if !ok {
		return nil, "", fmt.Errorf("user: %s, error: %w", oldReviewerID, dom.ErrUserNotFound)
	}
	team, ok := s.teams[user.TeamName]
	if !ok {
		return nil, "", fmt.Errorf("team: %s, error: %w", user.TeamName, dom.ErrTeamNotFound)
	}
	pr, ok := team.PullRequests[id]
	if !ok {
		return nil, "", fmt.Errorf("pr: %s, error: %w", id, dom.ErrPRNotFound)
	}
	if _, ok = pr.AssignedReviewers[user.ID()]; !ok {
		fmt.Println("user", user.ID())
		return nil, "", fmt.Errorf("user: %s, error: %w", user.ID(), dom.ErrUserNotAssigned)
	}
	if pr.Status == dom.PRMerged {
		return nil, "", fmt.Errorf("pr: %s, error: %w", pr.Name, dom.ErrPRMerged)
	}

	prAuthor := s.users[pr.AuthorID]

	possibleReviewers, err := getPossibleReviewers(team, pr, prAuthor)
	if err != nil {
		return nil, "", err
	}
	newReviewer := chooseRandomUser(possibleReviewers)
	pr.ReassignReviewer(oldReviewerID, newReviewer)

	reviewers := make([]string, 0, dom.MaxAssignedReviewers)
	for _, r := range pr.AssignedReviewers {
		reviewers = append(reviewers, r.ID())
	}
	outPR := types.CreatePullRequest(pr.ID(), pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt, reviewers)
	return outPR, newReviewer.ID(), nil
}

func getPossibleReviewers(team *dom.Team, pr *dom.PullRequest, authorPR *dom.User) ([]*dom.User, error) {
	disallowedReviewers := make(map[string]*dom.User, len(pr.AssignedReviewers)+1)
	for k, v := range pr.AssignedReviewers {
		disallowedReviewers[k] = v
	}
	disallowedReviewers[authorPR.ID()] = authorPR

	possibleReviewers := team.GetOtherActiveMembers(disallowedReviewers)
	if len(possibleReviewers) == 0 {
		return nil, fmt.Errorf("team: %s, error: %w", team.Name, dom.ErrUserNotActive)
	}
	return possibleReviewers, nil
}

func currentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

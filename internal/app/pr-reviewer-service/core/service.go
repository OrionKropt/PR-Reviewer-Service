package core

import (
	"fmt"
	"sync"

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

func (s *PRReviewerService) AddTeam(name string, members []types.TeamMember) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.teams[name]; ok {
		return fmt.Errorf("team %s already exists", name)
	}
	newTeam := dom.CreateTeam(name)
	for _, m := range members {
		newUser := dom.CreateUser(m.UserID, m.Username, name)
		newTeam.AddMember(newUser)
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
		return nil, fmt.Errorf("team %s not found", name)
	}
	members := team.GetMembers()
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

func (s *PRReviewerService) SetIsActive(userId string, isActive bool) (*types.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.users[userId]
	if !ok {
		return nil, fmt.Errorf("user %s not found", userId)
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

//func (s *PRReviewerService) createPullRequest(id, name, authorID string) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	s.prs[id] = dom.CreatePullRequest(id, authorID, name)
//}

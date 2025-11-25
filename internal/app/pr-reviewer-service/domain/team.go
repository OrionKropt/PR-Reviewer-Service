package domain

type Team struct {
	Name         string
	PullRequests map[string]*PullRequest
	Members      []*User
}

func CreateTeam(name string) Team {
	return Team{
		Name:         name,
		PullRequests: make(map[string]*PullRequest),
		Members:      make([]*User, 0),
	}
}

func (t *Team) AddMember(u *User) {
	t.Members = append(t.Members, u)
}

func (t *Team) GetActiveMembers() []*User {
	activeMembers := make([]*User, 0, len(t.Members))
	for _, member := range t.Members {
		if member.IsActive {
			activeMembers = append(activeMembers, member)
		}
	}
	return activeMembers
}

func (t *Team) GetOtherActiveMembers(users map[string]*User) []*User {
	activeMembers := t.GetActiveMembers()
	otherActiveMembers := make([]*User, 0, len(activeMembers))
	for _, m := range activeMembers {
		if _, ok := users[m.ID()]; !ok {
			otherActiveMembers = append(otherActiveMembers, m)
		}
	}

	return otherActiveMembers
}

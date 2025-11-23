package domain

type Team struct {
	Name    string
	members []User
}

func CreateTeam(name string) Team {
	return Team{
		Name:    name,
		members: make([]User, 0),
	}
}

func (t *Team) GetMembers() []User {
	return t.members
}

func (t *Team) AddMember(u User) {
	t.members = append(t.members, u)
}

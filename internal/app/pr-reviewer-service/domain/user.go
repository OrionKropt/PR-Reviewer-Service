package domain

type User struct {
	id       string
	Name     string
	TeamName string
	IsActive bool
}

func CreateUser(id, name, teamName string, isActive bool) User {
	return User{
		id:       id,
		Name:     name,
		TeamName: teamName,
		IsActive: isActive,
	}
}

func (u User) ID() string {
	return u.id
}

package domain

import "errors"

var (
	ErrPRExists          = errors.New("pull request already exists")
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrAuthorNotFound    = errors.New("author not found")
	ErrTeamNotFound      = errors.New("team not found")
)

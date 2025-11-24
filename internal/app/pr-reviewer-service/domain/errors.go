package domain

import "errors"

var (
	ErrPRExists          = errors.New("pull request already exists")
	ErrPRNotFound        = errors.New("PR not found")
	ErrPRMerged          = errors.New("PR merged")
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrAuthorNotFound    = errors.New("author not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotActive     = errors.New("user not active")
	ErrUserNotAssigned   = errors.New("user not assigned")
	ErrTeamNotFound      = errors.New("team not found")
)

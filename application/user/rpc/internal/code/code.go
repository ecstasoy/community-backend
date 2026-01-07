package code

import "community-backend/pkg/error"

var (
	RegisterNameEmpty = error.New(20001, "Register name cannot be empty")
)

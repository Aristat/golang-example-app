package common

import (
	"fmt"
)

type ErrorCode int

type Error struct {
	ErrorCode   ErrorCode
	Description string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s\n", e.ErrorCode, e.Description)
}

var (
	PasswordIsEmpty  = Error{ErrorCode: 10001, Description: "Password is empty"}
	UserNotFound     = Error{ErrorCode: 10002, Description: "User not found"}
	InvalidReturnUri = Error{ErrorCode: 10003, Description: "ReturnUri is not valid"}
)

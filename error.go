package coinbasepro

import "errors"

var (
	ErrNotFound     = Error{Message: "Route not found"}
	ErrUnauthorized = Error{Message: "Unauthorized."}

	ErrCloseWebsocket = errors.New("close webscoket connection")
)

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

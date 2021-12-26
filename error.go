package coinbasepro

var (
	ErrNotFound     = Error{Message: "Route not found"}
	ErrUnauthorized = Error{Message: "Unauthorized."}
)

type Error struct {
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

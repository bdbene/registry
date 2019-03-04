package server

import "fmt"

type ServerError struct {
	reason string
}

func (error *ServerError) Error() string {
	return fmt.Sprintf("Server failure: %s", error.reason)
}

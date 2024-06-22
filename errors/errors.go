package errors

import "fmt"

var (
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrForbidden    = fmt.Errorf("forbidden access")
	ErrWebsocket    = fmt.Errorf("websocket")
	ErrInvalidToken = fmt.Errorf("invalid token")
)

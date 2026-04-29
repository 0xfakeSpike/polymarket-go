package polymarket

import (
	"errors"
	"fmt"
)

// ApiError is returned when ThrowOnError is set and the API returns a JSON error payload.
type ApiError struct {
	Message string
	Status  int
	Data    any
}

func (e *ApiError) Error() string {
	if e.Status != 0 {
		return fmt.Sprintf("polymarket api: %s (status %d)", e.Message, e.Status)
	}
	return "polymarket api: " + e.Message
}

var (
	ErrL1AuthRequired = errors.New("signer is required for this endpoint")
	ErrL2AuthRequired = errors.New("API credentials are required for this endpoint")
)

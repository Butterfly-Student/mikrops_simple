package hotspot

import "fmt"

// Error types
var (
	ErrProfileNotFound = fmt.Errorf("profile not found")
	ErrUserNotFound    = fmt.Errorf("user not found")
	ErrSessionNotFound = fmt.Errorf("session not found")
	ErrSaleNotFound    = fmt.Errorf("sale record not found")
	ErrInvalidProfile  = fmt.Errorf("invalid profile")
	ErrInvalidUser     = fmt.Errorf("invalid user")
	ErrInvalidValidity = fmt.Errorf("invalid validity format")
	ErrExpiryMode      = fmt.Errorf("invalid expiry mode")
)

// HotspotError wraps errors with additional context
type HotspotError struct {
	Operation string
	Err       error
}

func (e *HotspotError) Error() string {
	return fmt.Sprintf("hotspot %s: %v", e.Operation, e.Err)
}

func (e *HotspotError) Unwrap() error {
	return e.Err
}

// NewError creates new HotspotError
func NewError(operation string, err error) error {
	return &HotspotError{
		Operation: operation,
		Err:       err,
	}
}

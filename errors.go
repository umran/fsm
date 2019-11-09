package fsm

import "fmt"

// GenericError ...
type GenericError struct {
	code    string
	message string
}

// Error ...
func (ge *GenericError) Error() string {
	return fmt.Sprintf("%s: %s", ge.code, ge.message)
}

var (
	// ErrUndefinedState ...
	ErrUndefinedState error = &GenericError{
		code:    "invalid state",
		message: "can't reference undefined state",
	}

	// ErrUndefinedTransition ...
	ErrUndefinedTransition error = &GenericError{
		code:    "invalid transition",
		message: "can't undergo an undefined transition",
	}

	// ErrNilToNonInitialTransition ...
	ErrNilToNonInitialTransition error = &GenericError{
		code:    "invalid transition",
		message: "can't transition from nil state to non-initial state",
	}
)

package fsm

import "fmt"

type genericError struct {
	code    string
	message string
}

// Error ...
func (ge *genericError) Error() string {
	return fmt.Sprintf("%s: %s", ge.code, ge.message)
}

var (
	// ErrIllegalStateName ...
	ErrIllegalStateName error = &genericError{
		code:    "illegal state name",
		message: "can't use reserved name as state name",
	}

	// ErrUndefinedState ...
	ErrUndefinedState error = &genericError{
		code:    "undefined state",
		message: "can't reference undefined state",
	}

	// ErrUndefinedTransition ...
	ErrUndefinedTransition error = &genericError{
		code:    "undefined transition",
		message: "can't undergo an undefined transition",
	}

	// ErrNilToNonInitialTransition ...
	ErrNilToNonInitialTransition error = &genericError{
		code:    "nil to non-initial transition",
		message: "can't transition from nil state to non-initial state",
	}
)

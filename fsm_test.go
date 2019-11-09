package fsm

import (
	"fmt"
	"testing"
)

const (
	OpenState   = "OPEN"
	ClosedState = "CLOSED"
	StoredState = "STORED"
)

type Box struct {
	machine *Machine
}

func (box *Box) ReconcileOpen(previousStateName string, args interface{}) error {
	switch previousStateName {
	case "":
		fmt.Printf("Box initialized to %s", OpenState)
	default:
		fmt.Printf("Box transitioning to %s from %s", OpenState, previousStateName)
	}

	return nil
}

func (box *Box) ReconcileStored(previousStateName string, args interface{}) error {
	fmt.Printf("Box transitioning to %s from %s", StoredState, previousStateName)
	return nil
}

func (box *Box) Initialize() error {
	machine, err := New("", map[string]*StateDefinition{
		OpenState: {
			InitialState: true,
			Transitions: []string{
				ClosedState,
			},
			On: box.ReconcileOpen,
		},
		ClosedState: {
			Transitions: []string{
				OpenState,
				StoredState,
			},
		},
		StoredState: {
			Transitions: []string{
				OpenState,
			},
			On: box.ReconcileStored,
		},
	})

	if err != nil {
		return err
	}

	box.machine = machine
	return nil
}

func NewBox() *Box {
	return &Box{}
}

func TestInitialTransition(t *testing.T) {
	box := NewBox()
	box.Initialize()

	if err := box.machine.ReconcileForState(OpenState, nil); err != nil {
		t.Error(err)
	}

	if box.machine.State() != OpenState {
		t.Error("unexpected state")
	}
}

func TestOnwardTransition(t *testing.T) {
	box := NewBox()
	box.Initialize()
	box.machine.ReconcileForState(OpenState, nil)

	if err := box.machine.ReconcileForState(ClosedState, nil); err != nil {
		t.Error(err)
	}

	if box.machine.State() != ClosedState {
		t.Error("unexpected state")
	}
}

func TestOnwardTransitionToExistingState(t *testing.T) {
	box := NewBox()
	box.Initialize()
	box.machine.ReconcileForState(OpenState, nil)

	if err := box.machine.ReconcileForState(OpenState, nil); err != nil {
		t.Error(err)
	}

	if box.machine.State() != OpenState {
		t.Error("unexpected state")
	}
}

func TestInvalidInitialTransition(t *testing.T) {
	box := NewBox()
	box.Initialize()

	err := box.machine.ReconcileForState(ClosedState, nil)

	if err != ErrNilToNonInitialTransition {
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestInvalidOnwardTransition(t *testing.T) {
	box := NewBox()
	box.Initialize()
	box.machine.ReconcileForState(OpenState, nil)

	err := box.machine.ReconcileForState(StoredState, nil)

	if err != ErrUndefinedTransition {
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestInvalidTransitionToNilState(t *testing.T) {
	box := NewBox()
	box.Initialize()

	err := box.machine.ReconcileForState("Bollocks", nil)

	if err != ErrUndefinedTransition {
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestInvalidMachine(t *testing.T) {
	_, err := New("", map[string]*StateDefinition{
		"ON": {
			Transitions: []string{
				"SOME_UNDEFINED_STATE",
			},
		},
		"OFF": {},
	})

	if err != ErrUndefinedState {
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestGenericError(t *testing.T) {
	err := &GenericError{
		code:    "test code",
		message: "test message",
	}

	fmt.Println(err)
}

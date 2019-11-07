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

func (box *Box) ReconcileOpen(previousState *State, args interface{}) error {
	switch previousState {
	case nil:
		fmt.Printf("Box initialized to %s", OpenState)
	default:
		fmt.Printf("Box transitioning to %s from %s", OpenState, previousState.Name)
	}

	return nil
}

func (box *Box) ReconcileStored(previousState *State, args interface{}) error {
	fmt.Printf("Box transitioning to %s from %s", StoredState, previousState.Name)
	return nil
}

func NewBox() *Box {
	box := &Box{}

	var (
		initialState *State
		states       map[string]*State
	)

	states = map[string]*State{
		OpenState: {
			Name:         OpenState,
			InitialState: true,
			Transitions: []func() *State{
				func() *State { return states[ClosedState] },
			},
			On: box.ReconcileOpen,
		},
		ClosedState: {
			Name:         ClosedState,
			InitialState: false,
			Transitions: []func() *State{
				func() *State { return states[OpenState] },
				func() *State { return states[StoredState] },
			},
		},
		StoredState: {
			Name:         StoredState,
			InitialState: false,
			Transitions: []func() *State{
				func() *State { return states[OpenState] },
			},
			On: box.ReconcileStored,
		},
	}

	box.machine = New(initialState, states)
	return box
}

func TestInitialTransition(t *testing.T) {
	box := NewBox()

	if err := box.machine.ReconcileForState(OpenState, nil); err != nil {
		t.Error(err)
	}

	if box.machine.State().Name != OpenState {
		t.Error("unexpected state")
	}
}

func TestOnwardTransition(t *testing.T) {
	box := NewBox()
	box.machine.ReconcileForState(OpenState, nil)

	if err := box.machine.ReconcileForState(ClosedState, nil); err != nil {
		t.Error(err)
	}

	if box.machine.State().Name != ClosedState {
		t.Error("unexpected state")
	}
}

func TestOnwardTransitionToExistingState(t *testing.T) {
	box := NewBox()
	box.machine.ReconcileForState(OpenState, nil)

	if err := box.machine.ReconcileForState(OpenState, nil); err != nil {
		t.Error(err)
	}

	if box.machine.State().Name != OpenState {
		t.Error("unexpected state")
	}
}

func TestInvalidInitialTransition(t *testing.T) {
	box := NewBox()

	err := box.machine.ReconcileForState(ClosedState, nil)
	if err == nil {
		t.Error("was expecting an error, but received none")
	}

	if err != ErrNilToNonInitialTransition {
		t.Error("unexpected error")
	}
}

func TestInvalidOnwardTransition(t *testing.T) {
	box := NewBox()
	box.machine.ReconcileForState(OpenState, nil)

	err := box.machine.ReconcileForState(StoredState, nil)
	if err == nil {
		t.Error("was expecting an error, but received none")
	}

	if err != ErrUndefinedTransition {
		t.Error("unexpected error")
	}
}

func TestInvalidTransitionToNilState(t *testing.T) {
	box := NewBox()

	err := box.machine.ReconcileForState("Bollocks", nil)
	if err == nil {
		t.Error("was expecting an error, but received none")
	}

	if err != ErrUndefinedTransition {
		t.Error("unexpected error")
	}
}

func TestGenericError(t *testing.T) {
	err := &GenericError{
		code:    "test code",
		message: "test message",
	}

	fmt.Println(err)
}

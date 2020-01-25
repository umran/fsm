package fsm

import (
	"errors"
	"fmt"
	"testing"
)

const (
	OpenState    = "OPEN"
	ClosedState  = "CLOSED"
	StoredState  = "STORED"
	FailingState = "FAILING"
)

type Box struct {
	machine *Machine
}

func (box *Box) ReconcileOpen(previousStateName string, args ...interface{}) error {
	switch previousStateName {
	case "":
		fmt.Printf("Box initialized to %s", OpenState)
	default:
		fmt.Printf("Box transitioning to %s from %s", OpenState, previousStateName)
	}

	return nil
}

func (box *Box) ReconcileStored(previousStateName string, args ...interface{}) error {
	fmt.Printf("Box transitioning to %s from %s", StoredState, previousStateName)
	return nil
}

func (box *Box) ReconcileFailing(previousStateName string, args ...interface{}) error {
	fmt.Printf("Box undergoing a failing transition to %s from %s", FailingState, previousStateName)
	return errors.New("failing")
}

func (box *Box) Initialize() error {
	machine, err := New("", map[string]StateDefinition{
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
		FailingState: {
			InitialState: true,
			On:           box.ReconcileFailing,
		},
	}, nil)

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

func TestTransitionWithFailingOnTransition(t *testing.T) {
	box := NewBox()
	box.Initialize()

	err := box.machine.ReconcileForState(FailingState, nil)

	switch {
	case err == nil:
		t.Error("was expecting an error")
	case err.Error() != "failing":
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestMachineWithReconcileUpdate(t *testing.T) {
	machine, _ := New("", map[string]StateDefinition{
		"STATE_1": {
			InitialState: true,
		},
	}, func(nextStateName string, args ...interface{}) error {
		return nil
	})

	if err := machine.ReconcileForState("STATE_1"); err != nil {
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestMachineWithFailingReconcileUpdate(t *testing.T) {
	machine, _ := New("", map[string]StateDefinition{
		"STATE_1": {
			InitialState: true,
		},
	}, func(nextStateName string, args ...interface{}) error {
		return errors.New("failing")
	})

	err := machine.ReconcileForState("STATE_1")
	switch {
	case err == nil:
		t.Error("was expecting an error")
	case err.Error() != "failing":
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestIllegalStateName(t *testing.T) {
	_, err := New("", map[string]StateDefinition{
		"": {
			Transitions: []string{
				"OFF",
			},
		},
		"OFF": {},
	}, nil)

	if err != ErrIllegalStateName {
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestUndefinedStateReference(t *testing.T) {
	_, err := New("", map[string]StateDefinition{
		"ON": {
			Transitions: []string{
				"SOME_UNDEFINED_STATE",
			},
		},
		"OFF": {},
	}, nil)

	if err != ErrUndefinedState {
		t.Error("unexpected error")
		t.Log(err)
	}
}

func TestGenericError(t *testing.T) {
	err := &genericError{
		code:    "test code",
		message: "test message",
	}

	fmt.Println(err)
}

package fsm

import "errors"

// Machine ...
type Machine struct {
	currentState *State
	states       map[string]*State
}

// State ...
func (machine *Machine) State() *State {
	return machine.currentState
}

// ReconcileForState ...
func (machine *Machine) ReconcileForState(nextStateName string, args interface{}) error {

	nextState := machine.states[nextStateName]

	switch machine.currentState {
	case nil:
		if nextState.InitialState == false {
			return errors.New("invalid transition: can't transition from nil state to non-initial state")
		}
	default:
		if machine.currentState.Name == nextStateName {
			return nil
		}

		if machine.currentState.IsPossibleTransition(nextState) == false {
			return errors.New("invalid transition: can't undergo an undefined transition")
		}
	}

	previousState := machine.currentState
	machine.currentState = nextState

	return machine.currentState.On(previousState, args)
}

// New ...
func New(initialState *State, states map[string]*State) *Machine {
	return &Machine{
		currentState: initialState,
		states:       states,
	}
}

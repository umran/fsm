package fsm

import "sync"

// Machine contains a collection of states and exists in exactly one of those states at any given time
type Machine struct {
	currentState    *state
	states          map[string]*state
	reconcileUpdate func(string, ...interface{}) error
	mux             sync.Mutex
}

// State returns the name of the current state of a given machine
func (machine *Machine) State() string {
	return machine.currentState.name
}

// ReconcileForState transitions the state of a given machine to that specified in the first argument.
// The second argument is an interface{} type that is passed to the 'On' function assigned to the state.
func (machine *Machine) ReconcileForState(nextStateName string, args ...interface{}) error {
	machine.mux.Lock()
	defer machine.mux.Unlock()

	nextState := machine.states[nextStateName]

	// disallow transition to nil state
	if nextState == nil {
		return ErrUndefinedTransition
	}

	switch machine.currentState {
	case nil:
		if nextState.initialState == false {
			return ErrNilToNonInitialTransition
		}
	default:
		if machine.currentState.name == nextStateName {
			return nil
		}

		if machine.currentState.isPossibleTransition(nextStateName) == false {
			return ErrUndefinedTransition
		}
	}

	previousState := machine.currentState
	machine.currentState = nextState
	if machine.reconcileUpdate != nil {
		if err := machine.reconcileUpdate(nextStateName, args...); err != nil {
			return err
		}
	}

	var previousStateName string
	if previousState != nil {
		previousStateName = previousState.name
	}

	if machine.currentState.on != nil {
		if err := machine.currentState.on(previousStateName, args...); err != nil {
			return err
		}
	}

	return nil
}

// New generates a new state machine according to the initial state and state definitions provided
func New(initialStateName string, definitions map[string]StateDefinition, reconcileUpdate func(string, ...interface{}) error) (*Machine, error) {
	states := make(map[string]*state, len(definitions))

	for name, def := range definitions {
		// validate the state name
		if name == "" {
			return nil, ErrIllegalStateName
		}

		// validate all transitions for the state def
		for _, transition := range def.Transitions {
			_, ok := definitions[transition]
			if !ok {
				return nil, ErrUndefinedState
			}
		}

		// create and add new state from definition
		states[name] = newState(name, &def)
	}

	machine := &Machine{
		currentState:    states[initialStateName],
		states:          states,
		reconcileUpdate: reconcileUpdate,
	}

	return machine, nil
}

package fsm

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

	// disallow transition to nil state
	if nextState == nil {
		return ErrUndefinedTransition
	}

	switch machine.currentState {
	case nil:
		if nextState.InitialState == false {
			return ErrNilToNonInitialTransition
		}
	default:
		if machine.currentState.Name == nextStateName {
			return nil
		}

		if machine.currentState.isPossibleTransition(nextState) == false {
			return ErrUndefinedTransition
		}
	}

	previousState := machine.currentState
	machine.currentState = nextState

	if machine.currentState.On != nil {
		return machine.currentState.On(previousState, args)
	}

	return nil
}

// New ...
func New(initialState *State, states map[string]*State) *Machine {
	return &Machine{
		currentState: initialState,
		states:       states,
	}
}

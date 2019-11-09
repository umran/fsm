package fsm

// Machine ...
type Machine struct {
	currentState *state
	states       map[string]*state
}

// State ...
func (machine *Machine) State() string {
	return machine.currentState.name
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

	var previousStateName string
	if previousState != nil {
		previousStateName = previousState.name
	}

	if machine.currentState.on != nil {
		return machine.currentState.on(previousStateName, args)
	}

	return nil
}

// New ...
func New(initialStateName string, definitions map[string]StateDefinition) (*Machine, error) {
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
		currentState: states[initialStateName],
		states:       states,
	}

	return machine, nil
}

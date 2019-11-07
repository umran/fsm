package fsm

// State ...
type State struct {
	Name         string
	InitialState bool
	Transitions  []*State
	On           func(*State, interface{}) error
}

// IsPossibleTransition ...
func (state *State) IsPossibleTransition(nextState *State) bool {
	// disallow transition to nil state
	if nextState == nil {
		return false
	}

	for _, possibleState := range state.Transitions {
		if possibleState.Name == nextState.Name {
			return true
		}
	}

	return false
}

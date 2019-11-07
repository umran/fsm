package fsm

// State ...
type State struct {
	Name         string
	InitialState bool
	Transitions  []func() *State
	On           func(*State, interface{}) error
}

// IsPossibleTransition ...
func (state *State) isPossibleTransition(nextState *State) bool {
	for _, possibleState := range state.Transitions {
		if possibleState().Name == nextState.Name {
			return true
		}
	}

	return false
}

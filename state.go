package fsm

type state struct {
	name         string
	initialState bool
	transitions  []string
	on           func(string, ...interface{}) error
}

func (state *state) isPossibleTransition(nextStateName string) bool {
	for _, possibleState := range state.transitions {
		if possibleState == nextStateName {
			return true
		}
	}

	return false
}

func newState(name string, def *StateDefinition) *state {
	return &state{
		name:         name,
		initialState: def.InitialState,
		transitions:  def.Transitions,
		on:           def.On,
	}
}

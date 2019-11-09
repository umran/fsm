package fsm

// StateDefinition is used to define state properties and behaviour
type StateDefinition struct {
	InitialState bool
	Transitions  []string
	On           func(string, interface{}) error
}

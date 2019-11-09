package fsm

// StateDefinition is used to define states and contain their properties
type StateDefinition struct {
	InitialState bool
	Transitions  []string
	On           func(string, interface{}) error
}

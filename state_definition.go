package fsm

// StateDefinition ...
type StateDefinition struct {
	InitialState bool
	Transitions  []string
	On           func(string, interface{}) error
}

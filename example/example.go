package order

import "github.com/umran/fsm"

const (
	Shipped        = "SHIPPED"
	InDepot        = "IN_DEPOT"
	OutForDelivery = "OUT_FOR_DELIVERY"
	Delivered      = "DELIVERED"
)

// Here we define the Order type which embeds the FSM
type Order struct {
	machine *fsm.Machine
}

// Method to call when transitioning to the SHIPPED state
func (order *Order) OnShipped(previousState *fsm.State, args interface{}) error {
	return nil
}

// Method to call when transitioning to the IN_DEPOT state
func (order *Order) OnInDepot(previousState *fsm.State, args interface{}) error {
	return nil
}

// Method to call when transitioning to the OUT_FOR_DELIVERY state
func (order *Order) OnOutForDelivery(previousState *fsm.State, args interface{}) error {
	return nil
}

// Method to call when transitioning to the DELIVERED state
func (order *Order) OnDelivered(previousState *fsm.State, args interface{}) error {
	return nil
}

func (order *Order) Initialize(initialStateName string) {
	var states map[string]*fsm.State

	states = map[string]*fsm.State{
		Shipped: {
			// The name of the state
			Name: Shipped,
			// Indicates whether the machine can transition from a nil state to this state
			InitialState: true,
			// A list of possible transitions from this state
			Transitions: []func() *fsm.State{
				func() *fsm.State { return states[InDepot] },
			},
			// An optional method that is called on transition to this state
			On: order.OnShipped,
		},

		InDepot: {
			Name:         InDepot,
			InitialState: false,
			Transitions: []func() *fsm.State{
				func() *fsm.State { return states[OutForDelivery] },
			},
			On: order.OnInDepot,
		},

		OutForDelivery: {
			Name:         OutForDelivery,
			InitialState: false,
			Transitions: []func() *fsm.State{
				func() *fsm.State { return states[InDepot] },
				func() *fsm.State { return states[Delivered] },
			},
			On: order.OnOutForDelivery,
		},

		Delivered: {
			Name:         Delivered,
			InitialState: false,
			Transitions:  []func() *fsm.State{},
			On:           order.OnDelivered,
		},
	}

	order.machine = fsm.New(states[initialStateName], states)
}

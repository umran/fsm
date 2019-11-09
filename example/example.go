package example

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
func (order *Order) OnShipped(previousState string, args interface{}) error {
	return nil
}

// Method to call when transitioning to the IN_DEPOT state
func (order *Order) OnInDepot(previousState string, args interface{}) error {
	return nil
}

// Method to call when transitioning to the OUT_FOR_DELIVERY state
func (order *Order) OnOutForDelivery(previousState string, args interface{}) error {
	return nil
}

// Method to call when transitioning to the DELIVERED state
func (order *Order) OnDelivered(previousState string, args interface{}) error {
	return nil
}

func (order *Order) Initialize(initialStateName string) error {
	machine, err := fsm.New("", map[string]fsm.StateDefinition{
		Shipped: {
			// Indicates whether the machine can transition from a nil state to this state
			InitialState: true,
			// A list of possible transitions from this state
			Transitions: []string{
				InDepot,
			},
			// An optional method that is called on transition to this state
			On: order.OnShipped,
		},
		InDepot: {
			Transitions: []string{
				OutForDelivery,
			},
			On: order.OnInDepot,
		},
		OutForDelivery: {
			Transitions: []string{
				InDepot,
				Delivered,
			},
			On: order.OnOutForDelivery,
		},
		Delivered: {
			On: order.OnDelivered,
		},
	})

	if err != nil {
		return err
	}

	order.machine = machine
	return nil
}

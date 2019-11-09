# FSM: A library with a simple API to generate finite state machines in go

## Example
Suppose we wanted to implement the following state machine for some entity, call it Order:

<img width="561" alt="Screen Shot 2019-11-07 at 1 32 13 PM" src="https://user-images.githubusercontent.com/1547890/68429491-0bd32b00-0163-11ea-8893-b35a6a7eda10.png">

We can do so by first creating an Order type that embeds a state machine:
````go
package order

import "github.com/umran/fsm"

// Here we define the Order type which embeds the FSM
type Order struct {
	machine *fsm.Machine
}
````

### Defining state names
Before we define the state machine, it would be handy to have all possible state names defined somewhere as constants:
````go
const (
	Shipped        = "SHIPPED"
	InDepot        = "IN_DEPOT"
	OutForDelivery = "OUT_FOR_DELIVERY"
	Delivered      = "DELIVERED"
)
````

### Defining methods that are called when transitioning to particular states
We can also have the machine do stuff when transitioning to a new state. This can be done by defining methods to be called when transitioning to particular states. Such methods take as arguments the previous state name, an optional args argument, which is an `interface{}` type, and return an error.
````go
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
````

### Defining the state machine
Once the state names and event methods have been defined, we can define the state machine like so:
````go
func (order *Order) Initialize() error {
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
````
The state machine is defined inside the `Initialize()` method of the `Order` type. Inside this method, we create a new state machine initialized to the desired state and assign it to the `order`'s `machine` field:
````go
machine, err := fsm.New("", ...)
````
Note that the `initialStateName` argument passed to `fsm.New()` can be empty: `""`. In fact, due to the way this library is designed, if the initial state's transition method is to be executed, it is necessary that `initialStateName` be empty: `""` during the creation of the state machine. The initial state should instead be set by calling `ReconcileForState()` on the state machine once it has been generated. This is covered in the following section.

### Transitioning State
The state machine transitions state via calls to: `ReconcileForState(nextStateName string, args interface{})`
To continue with our `Order`example, new orders can be initialized to the `Shipped` state like so:
````go
order := new(Order)

order.Initialize()

// Note how the initial state is set by calling ReconcileForState
// This way the OnShipped method is called when the order is
// initialized to the Shipped state
order.machine.ReconcileForState(Shipped, nil)
````
In the real world we would almost always call ReconcileForState within another method that appropriately represents the business logic of our application:
````go
func (order *Order) Ship(trackingID string) error {
	return order.machine.ReconcileForState(Shipped, trackingID)
}

func (order *Order) MarkAsDelivered(signature string) error {
	return order.machine.ReconcileForState(Delivered, signature)
}
````

### Transitioning to the Current State
If `ReconcileForState()` is called where the next state is the current state, it will return immediately, i.e. without calling the method (if one is defined) that would otherwise be called when the the machine transitions to the state in question.

### Dealing with Invalid Transitions
If the requested transition is not a defined transition for the current state, `ReconcileForState()` will return the error: `ErrUndefinedTransition`

If the requested transition is from a `nil` state (as in the case of a machine that has not yet had its state initialized) to a state that is not an initial state, `ReconcileForState()` will return the error: `ErrNilToNonInitialTransition`
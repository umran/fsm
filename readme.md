# FSM: A Simple Library to Generate Finite State Machines in Go

## Usage
Suppose we wanted to implement the following state machine for an entity called Order:

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
Before we define the state machine that models the behaviour of the Order type, it would be handy to have all possible state names defined somewhere as constants:
````go
const (
	Shipped        = "SHIPPED"
	InDepot        = "IN_DEPOT"
	OutForDelivery = "OUT_FOR_DELIVERY"
	Delivered      = "DELIVERED"
)
````

### Defining methods that are called when transitioning to particular states
We can also have the machine do stuff when transitioning to a new state. This can be done by defining methods that are called when transitioning to particular states. Such methods take the previous state as the first argument, an optional args argument, which is an `interface{}` type, as the second argument, and return an error.
````go
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
````

### Defining the state machine
Finally, once the state names and event methods have been defined, we define the state machine like so:
````go
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
````
The state machine is defined inside the `Initialize()` method of the `Order` type. Inside this method, we create a new state machine initialized to the desired state and assign it to the `order`'s `machine` field:
````go
order.machine = fsm.New(states[initialStateName], states)
````
Note that the `initialStateName` variable passed to the `fsm.New()` function can be empty: `""`. In fact, due to the way this library is designed, if the initial state's transition method is to be executed, it is necessary that `initialStateName` be empty: `""` during the creation of the state machine. The initial state should instead be set by calling `ReconcileForState()` on the state machine once it has been generated. This is covered in the following section.

### Transitioning State
The state machine transitions state via calls to: `ReconcileForState(nextStateName string, args interface{})`
To continue with our `Order`example, new orders can be initialized to the `Shipped` state like so:
````go
order := new(Order)

// Note how the initialStateName argument is an empty string
order.Initialize("")

// This way the OnShipped method is called when
// calling ReconcileForState on the Shipped state
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
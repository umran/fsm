# FSM: A library with a simple API to generate finite state machines in go
[![Documentation](https://godoc.org/github.com/umran/fsm?status.svg)](http://godoc.org/github.com/umran/fsm)
[![Build Status](https://travis-ci.com/umran/fsm.svg?branch=master)](https://travis-ci.com/umran/fsm)
[![Coverage Status](https://coveralls.io/repos/github/umran/fsm/badge.svg?branch=master)](https://coveralls.io/github/umran/fsm?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/umran/fsm)](https://goreportcard.com/report/github.com/umran/fsm)

## API
### State Definitions
The `StateDefinition` type allows us to define a state, which consists of the following fields:

1. `InitialState`: a boolean that indicates whether the machine can transition from a nil state to the state in question
2. `Transitions`: a list of the names of states it is possible to transition to from the state in question
3. `On`: A function that is called when transitioning to the state in question. It receives the previous state name (a `string`) as the first argument and an arbitrary `interface{}` type as the second argument

Please note that all of the above fields are optional and do not have to be defined for all states.
For example a particular state might not allow further transitions, in which case its `Transitions` field would be `nil`. Leaving it undefined in such a case would be completely fine. 

The same applies to the `On` field. If there is nothing further to be done on transitioning to a particular state, its `On` field maybe be ignored.

It is worth noting that if `InitialState` is not specified, it defaults to false. So if a particular state is not an initial state we may safely leave out its `InitialState` property.

### Machines
#### New(): Generating a new machine
A `Machine` is simply a collection of states and exists in a particular state at any given time. A machine can be in a `nil` state until it is initialized to an initial state. To create a new machine, one must call `New()` with 2 arguments:

1. the first argument is a `string` that indicates which state the machine should occupy when it is first created. This value can be an empty `string`: `""`, in which case the machine would occupy a `nil` state when it is first created
2. the second argument is a map from state names to `StateDefinitions` and defines all the possible states the machine can occupy over its lifetime

The `New()` function returns a new machine and an error. The only cases where an error is returned are:

1. If any of the states is named the empty `string`: `""`, in which case the following error is returned: `ErrIllegalStateName`
2. if any of the state definitions lists an undefined state under `Transitions`, in which case the following error is returned: `ErrUndefinedState`

````go
machine, err := fsm.New("", map[string]fsm.StateDefinition{
	"STATE_1": {
		InitialState: true,
		Transitions: []string{
			"STATE_2",
		},
		On: func(previousState string, args interface{}) error {
			// this is just a placeholder function that doesn't do anything
			return nil
		},
	},
	"STATE_2": {
		Transitions: []string{
			"STATE_3",
		},
	},
	"STATE_3": {},
})
````

#### State(): Getting the current state of the machine
We can get the current state of a machine by calling its `State()` method. This method returns a `string` specifying the name of the current state.

````go
currentStateName := machine.State()
````

#### ReconcileForState(): Transitioning the state of the machine
To transition a machine's state, we call the machine's `ReconcileForState()` method. This method requires two arguments and returns an error:

1. the first argument is a `string` indicating the name of the state to transition to
2. the second argument is an `interface{}` type and is passed to the state's `On` function (if it is defined)

````go
err := machine.ReconcileforState("STATE_1", nil)
````

If `ReconcileForState()` is called with the machine's current state, it will return immediately, since the machine is already in the desired state. Please note that in this case the `On` function of the state is not called. For this reason, it is sometimes necessary to provide an empty initial state when generating a new machine in order to make sure that the associated `On` function is called when the machine eventually assumes the desired initial state.

When `ReconcileForState()` is called, it determines if the state transition is allowed. If the transition is not allowed, it will return the following error: `ErrUndefinedTransition`.

Alternately, if the current state of the machine is `nil` and the next state does not have its `InitialState` field set to `true`, the following error will be returned: `ErrNilToNonInitialTransition`



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

### Defining "On" methods
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
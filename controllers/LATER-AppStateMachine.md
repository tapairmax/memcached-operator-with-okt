# Application State Machine implementation with OKT

OKT Stands for Operators Karma Tools that will be soon released under Apache 2.0 License, copyright 2021 Orange SA

This state machine GO API described below, through the tests, is already implemented by an "Operators Karma Tools" GO module. 


## Principle

An application can take several states and exposing the current state it has during a run is like pointing the location of a mobile object on a map. 

Having a Graph to describe what is managed by the K8S Operator should be helpful for a human operator or external components in charge to bring some "observability" features to a K8S operator.

Traversing a Graph has also the advantage to maintain a consistent path of all traversed nodes. This can help too to better understand or investigate on what is happening during a run.

Offering a developement framework to build this graph in a consistent way for a K8S operator should allow to normalize the view on the application's lifecycle as long as it is operated by this operator implemented with this framework (using a standard like CNCF CloudEvent could enforce this normalization and interoperability).

A Graph contains nodes and leaf nodes. Evolving between nodes is constraint by defined transitions to pass on known verified events.

The principle is not to traverse all the graph during 1 reconciliation but, at each reconciliation, to:
  + get the current state of the application (a database or whatever application), 
  + collect events that happened, and generate, if needed, the event telling to go to the state accordingly to the expected state in CR (add it to the collected events),
  + trigger (or not) a next state in regards to the collected events
  + in case of new state change, perform asynchroniously the actions related to this state
  + if the expected state is not the current state, requeue a controller runtime event for further Reconciliation...

The proposition is to implement an OKT App resource to manage the application life cycle like we implement an OKT resource mutator to manage different kind of Kubernetes resources: 
+ The application's state is evolving through the multiple Reconciliations,
+ The App resource implements all actions to do at each state, 
+ The expected state is specified in the CR, 
+ The current state is got from a specific Client call to pick up the information and map it into the corresponding state representation.

It based on:
+ A state machine based on what is offered by the OKT's GO module `tools/statemachine`
+ A Client implementing OKT's Client interface (CRUD) that communicate with the application 
+ An OKT application resource that can be registered by the OKT Reconciler


## Example of Life Cycle Graph implementation as currently offered by the OKT tools/statemachine GO module


    // App's States ID
    const (
        Start oktsm.LCGState = iota
        Run
        Service
        Stop
        End
    )

    // Life Cycle Graph (LCG) for the application. We retrieve here the State nodes (having a state ID) and 
    // an Information structure that provides the node's children.
    // A child ID in the children list describe the next State nodes that can be reached right after.
    // The child ID can be used too in an event list that appended dring the current state.
    // The library provide a way to trigger the state change through a list of events (state IDs).
    // Thus, the order of children is important, it defines which next node will be triggered in case of 
    // multiple events matching the children list.
    // I.e. in the children list of the "Start" state, both "End" and "Run" can trigger a state change. However,
    // in this description below, the "End" state has the priority over the "Run" state. 
    var appGraph = oktsm.LCGGraph{
        Start: oktsm.LCGNodeInfo{
            Name: "Start",
            Children: oktsm.LCGChildren{
                End,
                Run,
            },
        },
        Run: oktsm.LCGNodeInfo{
            Name: "Run",
            Children: oktsm.LCGChildren{
                Service,
                Stop,
            },
        },
        Service: oktsm.LCGNodeInfo{
            Name: "Servicing",
            Children: oktsm.LCGChildren{
                Stop,
                Run,
            },
        },
        Stop: oktsm.LCGNodeInfo{
            Name: "Stopping",
            Children: oktsm.LCGChildren{
                End,
            },
        },
        End: oktsm.LCGNodeInfo{Name: "End"}, // Leaf node: End of application's life
    }

    type OKTDatabase struct {
    }

    var _ LCGStateAction = &OKTDatabase{}

    func (db *OKTDatabase) Enter(state LCGState) error {
        stateName := appGraph.StateName(state)
        fmt.Println("Enter in state: ", stateName)

        switch state {
        case Start:
            fmt.Println("The DB is starting")
        case Run:
            fmt.Println("The DB is running now")
        case Service:
            fmt.Println("The DB is in a servicing operation")
        case Stop:
            fmt.Println("The DB is stopping")
        case End:
            fmt.Println("The livecycle machine of the DB is OFF")
        }

        return nil
    }

    // Changing from a node to another is allowed thanks to an GO API in OKT module `tools/statemachine`
    // See the test fonction of this module here
    func TestStateMachine(t *testing.T) {
        db := &OKTDatabase{}
        sm := &Machine{Graph: appGraph, Actions: db}

        off := sm.IsOFF()
        require.True(t, off, "At creation machine should be OFF")

        pStatus := sm.IsPathInGraphEnabled()
        require.False(t, pStatus, "PathInGraph should be disabled by default")
        sm.EnablePathInGraph()
        pStatus = sm.IsPathInGraphEnabled()
        require.True(t, pStatus, "PathInGraph should be enabled afeter Init")

        entered := sm.SetState(Start)
        require.True(t, entered, "Should be entered in Start state")
        off = sm.IsOFF()
        require.False(t, off, "Now, after entering in a state, the machine should be ON")

        var err error

        state := sm.GetState()
        require.Equal(t, Start, state, "Current state must be Start")
        entered, err = sm.EnterNextState(LCGEvents{End})
        require.True(t, entered, "Must be entered in state End")
        require.Nil(t, err, "No error must be raised in End state")

        path := sm.GetPathInGraph()
        fmt.Println(path)

        off = sm.IsOFF()
        require.True(t, off, "The state machine should be OFF after traversing a leaf node in the graph")
        state = sm.GetState()
        require.Equal(t, DefaultState, state, "Current state must be End")

        entered, err = sm.EnterNextState(LCGEvents{Start})
        require.False(t, entered, "Can't go to the next state when the machine is OFF. SetState() must be called before")
        require.Nil(t, err, "No error must be raised when a state has not been browsed")

        entered = sm.SetState(Run)
        require.True(t, entered, "Must be entered in state Run")
        state = sm.GetState()
        require.Equal(t, Run, state, "Current state must be Run")

        entered, err = sm.EnterNextState(LCGEvents{Start})
        require.False(t, entered, "Start is not a child state of Run")
        require.Nil(t, err, "No error must be raised when a state has not been browsed")

        entered, err = sm.EnterNextState(LCGEvents{Service})
        require.True(t, entered, "Must be entered in state of Service")
        require.Nil(t, err, "No error must be raised when a state has not been browsed")

        for count := 11; count > 0; count-- {
            entered, err = sm.EnterNextState(LCGEvents{DefaultState})
            require.True(t, entered, "Must be entered in state of Service")
            require.Nil(t, err, "No error must be raised here")
            state = sm.GetState()
            require.Equal(t, Run, state, "Current state must be Run now it is the Default after Service")

            entered, err = sm.EnterNextState(LCGEvents{Service})
            require.True(t, entered, "Must be entered in state of Service")
            require.Nil(t, err, "No error must be raised here")
        }

        entered, err = sm.EnterNextState(LCGEvents{Run, Stop}) // Priority test
        require.True(t, entered, "Must be entered in Run state")
        require.Nil(t, err, "No error must be raised here")
        state = sm.GetState()
        require.Equal(t, Stop, state, "Current state must be Stop as Run is not the prior state")

        entered, err = sm.EnterNextState(LCGEvents{DefaultState})
        require.True(t, entered, "Must be entered in  state")
        require.Nil(t, err, "No error must be raised here")
        off = sm.IsOFF()
        require.True(t, off, "The state machine should be OFF after traversing a leaf node in the graph")
        state = sm.GetState()
        require.Equal(t, DefaultState, state, "Current state must undefined/default when machine is OFF")

        path = sm.GetPathInGraph()
        fmt.Println(path)

        // Start>Run>Stop>End
        off = sm.IsOFF()
        require.True(t, off, "The state machine should be OFF after traversing a leaf node in the graph")
        entered = sm.SetState(Start)
        require.True(t, entered, "Should be entered in Start state")

        for sm.IsOFF() == false {
          entered, _ = sm.EnterNextState(LCGEvents{DefaultState})
          require.True(t, entered, "Should be entered in Start state")
        }
        require.True(t, entered, "Should be entered in state")

        path = sm.GetPathInGraph()
        fmt.Println(path)
      }

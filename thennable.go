package thennable

/*
IThennable is an interface for thennable
The main purpose to to be able to chain execution of functions
The chainnable function must be of type FRunnable
Example 1: Normal execution
    result, err := thennable.Start(1). --initial value of 1
        --AddOne is a function that conforms to FRunnable. Adds 1 to supplied value.
        Then(thennable.New(AddOne)). --1 + 1 = 2
        --AddTwo is a function that conforms to FRunnable. Adds 2 to supplied value.
        Then(thennable.New(AddTwo)). --2 + 2 = 4
        --AddSix is a function that conforms to FRunnable. Adds 6 to supplied value.
        Then(thennable.New(AddSix)). --4 + 6 = 10
        End()
    at this point, result = 10, err = nil
Example 2: Break on error
    result, err := thennable.Start(1). --initial value of 1
        Then(thennable.New(AddOne)). --1 + 1 = 2
        Then(thennable.New(AddTwo)). --2 + 2 = 4
        --Exception is a function that conforms to FRunnable and returns an error
        Then(thennable.New(Exception). -- err = error
        Then(thennable.New(AddSix)). --AddSix ignored
        End()
    at this point, result = 4, err = error
*/
type IThennable interface {
    BreakOnError(bool) IThennable
    Supply(interface{}) IThennable
    Then(IThennable) IThennable
    Handle(FErrorHandler) IThennable
    Start(interface{}) IThennable
    End() (interface{}, error)
}

//FRunnable is the actual function we want to chain
type FRunnable func(interface{}) (interface{}, error)

type FErrorHandler func(error)

//private implementation of IThennable
type thennable struct {
    throw        error
    state        interface{}
    runnable     FRunnable
    breakOnError bool
}

//Start creating a new IThennable with and initial value/state
func Start(state interface{}) IThennable {
    return &thennable{state: state, breakOnError: true}
}

//New IThennable with a runnable function
func New(runnable FRunnable) IThennable {
    return &thennable{runnable: runnable, breakOnError: true}
}


//BreakOnError stops function chaining when error occured
func (tnb *thennable) BreakOnError(breakOnError bool) IThennable {
    tnb.breakOnError = breakOnError
    return tnb
}

//Supply next runnable function with a value
func (tnb *thennable) Supply(param interface{}) IThennable {
    return &thennable{state: param, breakOnError: tnb.breakOnError}
}

//Then do next runnable function
func (tnb *thennable) Then(next IThennable) IThennable {
    if tnb.breakOnError && tnb.throw != nil {
        return tnb
    }
    
    return next.BreakOnError(tnb.breakOnError).Start(tnb.state)
}

//Handle error
func (tnb *thennable) Handle(handle FErrorHandler) IThennable {
    handle(tnb.throw)
    return tnb
}

//Start current runnable function
func (tnb *thennable) Start(param interface{}) IThennable {
    // if tnb.throw != nil && tnb.breakOnError {
    //     tnb.state, tnb.throw = param, err
    // }
    tnb.state, tnb.throw = tnb.runnable(param)
    return tnb
}

//End execution and collect the result
func (tnb *thennable) End() (interface{}, error) {
    return tnb.state, tnb.throw
}
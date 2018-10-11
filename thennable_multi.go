package thennable

/*
IThennableMulti is similar to IThennable, but receive & returns multiple values
*/
type IThennableMulti interface {
    BreakOnError(bool) IThennableMulti
    Supply(...interface{}) IThennableMulti
    Then(IThennableMulti) IThennableMulti
    Handle(FErrorHandler) IThennableMulti
    Start(...interface{}) IThennableMulti
    End() ([]interface{}, error)
}

//FRunnableMulti is the actual function we want to chain
type FRunnableMulti func(...interface{}) ([]interface{}, error)

//private implementation of IThennable
type thenmulti struct {
    throw        error
    state        []interface{}
    runnable     FRunnableMulti
    breakOnError bool
}

//StartMulti creating a new IThennableMulti with and initial value/state
func StartMulti(state ...interface{}) IThennableMulti {
    return &thenmulti{state: state, breakOnError: true}
}

//New IThennable with a runnable function
func NewMulti(runnable FRunnableMulti) IThennableMulti {
    return &thenmulti{runnable: runnable, breakOnError: true}
}


//BreakOnError stops function chaining when error occured
func (tnb *thenmulti) BreakOnError(breakOnError bool) IThennableMulti {
    tnb.breakOnError = breakOnError
    return tnb
}

//Supply next function with some values
func (tnb *thenmulti) Supply(param ...interface{}) IThennableMulti {
    return &thenmulti{state: param, breakOnError: tnb.breakOnError}
}

//Then do next runnable function
func (tnb *thenmulti) Then(next IThennableMulti) IThennableMulti {
    if tnb.breakOnError && tnb.throw != nil {
        return tnb
    }
    
    return next.BreakOnError(tnb.breakOnError).Start(tnb.state...)
}

//Handle error
func (tnb *thenmulti) Handle(handle FErrorHandler) IThennableMulti {
    handle(tnb.throw)
    return tnb
}

//Start current runnable function
func (tnb *thenmulti) Start(param ...interface{}) IThennableMulti {
    tnb.state, tnb.throw = tnb.runnable(param...)
    return tnb
}

//End execution and collect the result
func (tnb *thenmulti) End() ([]interface{}, error) {
    return tnb.state, tnb.throw
}
package thennable

import (
	"errors"
	"reflect"
)

//Error collection
var (
	ErrNotFunction     = errors.New("Argument supplied must be a function")
	ErrNoErrorHandling = errors.New("Function last return value must be an error")
)

/*
IThennable is an interface for thennable
The main purpose to to be able to chain execution of functions
*/
type IThennable interface {
	BreakOnError(bool) IThennable
	Supply(...interface{}) IThennable
	Then(interface{}) IThennable
	Handle(FErrorHandler) IThennable
	End() ([]interface{}, error)
}

//FErrorHandler function
type FErrorHandler func(error)

//private implementation of IThennable
type thennable struct {
	throw        error
	state        []reflect.Value
	runnable     reflect.Value
	breakOnError bool
}

//Start creating a new IThennable with and initial value/state
func Start(params ...interface{}) IThennable {
	state := []reflect.Value{}
	for _, param := range params {
		state = append(state, reflect.ValueOf(param))
	}
	return &thennable{state: state, breakOnError: true}
}

//Runnable function is a function that have at least error return value
func newRunnable(fn interface{}) (runnable reflect.Value, err error) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return runnable, ErrNotFunction
	}

	out := t.NumOut()
	if out <= 0 || t.Out(out-1).Name() != "error" {
		return runnable, ErrNoErrorHandling
	}

	return reflect.ValueOf(fn), nil
}

//BreakOnError stops function chaining when error occured
func (tnb *thennable) BreakOnError(breakOnError bool) IThennable {
	tnb.breakOnError = breakOnError
	return tnb
}

//Supply next runnable function with a value
func (tnb *thennable) Supply(params ...interface{}) IThennable {
	state := []reflect.Value{}
	for _, param := range params {
		state = append(state, reflect.ValueOf(param))
	}

	return &thennable{state: state, breakOnError: tnb.breakOnError}
}

func (tnb *thennable) Then(next interface{}) IThennable {
	if tnb.breakOnError && tnb.throw != nil {
		return tnb
	}

	//if not runnable func
	runnable, err := newRunnable(next)
	if err != nil {
		tnb.state, tnb.throw = make([]reflect.Value, 0), err
		return tnb
	}

	//else
	result := runnable.Call(tnb.state)
	retCount := len(result)
	errIdx := retCount - 1

	tnb.state, tnb.throw = result[0:errIdx], nil
	lastOutput := result[errIdx].Interface()
	if lastOutput != nil {
		tnb.throw = lastOutput.(error)
	}

	return tnb
}

//Handle error
func (tnb *thennable) Handle(handle FErrorHandler) IThennable {
	handle(tnb.throw)
	return tnb
}

//End execution and collect the result
func (tnb *thennable) End() ([]interface{}, error) {
	end := make([]interface{}, 0)
	for _, state := range tnb.state {
		end = append(end, state.Interface())
	}
	return end, tnb.throw
}

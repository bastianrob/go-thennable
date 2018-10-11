package thennable_test

import (
    "fmt"
    "thennable"
    "testing"
)

func AddOne(val interface{}) (interface{}, error) {
    intVal := val.(int)
    return intVal + 1, nil
}

func Throwable(val interface{}) (interface{}, error) {
    return nil, fmt.Errorf("Exception occurred with supplied value: %+v", val)
}

func LogError(err error) {
    fmt.Printf("Log Error: %+v\n", err)
}

func Test_IThennable(t *testing.T) {
    fmt.Println("\nStarting thennable_test.go")
    
    addOne := thennable.New(AddOne)
    throw := thennable.New(Throwable)
    
    //Test 1
    fmt.Printf("Start Test 1: Expect result = 5 and err is nil\n")
    result, err := thennable.Start(1).//starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(addOne). //result = 4
        Then(addOne). //result = 5
        End()
    
    if result.(int) != 5 {
        t.Errorf("Test 1 result should be = 5, actual value: %d", result)
    }
    if err != nil {
        t.Errorf("Test 1 error should be = nil, actual value: %+v", err)
    }
    
    //Test 2
    fmt.Printf("\nStart Test 2: Expect result is nil and err is not nil\n")
    result, err = thennable.Start(1).//starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //add 2 with 1 = 3
        Then(throw). //result = nil, err = exception
        Handle(LogError).
        Then(addOne). //skipped
        Then(addOne). //skipped
        End()
        
    if result != nil {
        t.Errorf("Test 2 result should be = nil, actual value: %d", result)
    }
    if err == nil {
        t.Errorf("Test 2 error should not be = nil, actual value: %+v", err)
    }
    
    //Test 3
    fmt.Printf("\nStart Test 3: Expect result = 2 and err is nil\n")
    result, err = thennable.Start(1).//starts with 1
        BreakOnError(false). //keep propagating the function chain
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(throw).  //result = 0, err = exception ignored
        Supply(0).    //resupply with 0
        Then(addOne). //add 0 with 1 = 1, err = nil
        Then(addOne). //result = 2, err = nil
        End()
        
    if result.(int) != 2 {
        t.Errorf("Test 3 result should be = 2, actual value: %d", result)
    }
    if err != nil {
        t.Errorf("Test 3 error should be = nil, actual value: %+v", err)
    }
    
    //Test 4
    fmt.Printf("\nStart Test 4: Expect recover from error, result is 10 and err is nil\n")
    result, err = thennable.Start(1).//starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(throw). //result = nil, err = exception
        Then(addOne). //skipped
        Handle(LogError).    //log the error
        BreakOnError(false). //recover from error
        Supply(8).    //resuply the value with 8
        Then(addOne). //result = 9
        Then(addOne). //result = 10
        Handle(LogError).    //log error should be nil
        End()
        
    if result.(int) != 10 {
        t.Errorf("Test 4 result should be = 10, actual value: %d", result)
    }
    if err != nil {
        t.Errorf("Test 4 error should be nil, actual value: %+v", err)
    }
}
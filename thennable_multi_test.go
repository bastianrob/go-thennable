package thennable_test

import (
    "fmt"
    "thennable"
    "testing"
)

func AddOneMulti(val ...interface{}) ([]interface{}, error) {
    intVal := val[0].(int)
    return []interface{}{intVal + 1}, nil
}

func ThrowableMulti(val ...interface{}) ([]interface{}, error) {
    return nil, fmt.Errorf("Exception occurred with supplied value: %+v", val)
}

func Test_IThennableMulti(t *testing.T) {
    fmt.Println("\nStarting thennable_multi_test.go")
    
    addOne := thennable.NewMulti(AddOneMulti)
    throw := thennable.NewMulti(ThrowableMulti)
    
    //Test 1
    fmt.Printf("Start Test 1: Expect result = 5 and err is nil\n")
    result, err := thennable.StartMulti(1). //starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(addOne). //result = 4
        Then(addOne). //result = 5
        End()
    
    if result[0].(int) != 5 {
        t.Errorf("Test 1 result should be = 5, actual value: %d", result)
    }
    if err != nil {
        t.Errorf("Test 1 error should be = nil, actual value: %+v", err)
    }
    
    //Test 2
    fmt.Printf("\nStart Test 2: Expect result is nil and err is not nil\n")
    result, err = thennable.StartMulti(1). //starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //add 2 with 1 = 3
        Then(throw).  //result = nil, err = exception
        Then(addOne). //skipped
        Then(addOne). //skipped
        End()
        
    if result != nil {
        t.Errorf("Test 2 result should be = nil, actual value: %d", result)
    }
    if err == nil {
        t.Errorf("Test 2 error should not be = nil, actual value: %+v", err)
    }
    
    //Test 4
    fmt.Printf("\nStart Test 4: Expect recover from error, result is 10 and err is nil\n")
    result, err = thennable.StartMulti(1).//starts with 1
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
        
    if result[0].(int) != 10 {
        t.Errorf("Test 4 result should be = 10, actual value: %d", result)
    }
    if err != nil {
        t.Errorf("Test 4 error should be nil, actual value: %+v", err)
    }
}
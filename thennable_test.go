package thennable_test

import (
    "fmt"
    "thennable"
    "testing"
)

func AddOne(val interface{}, err error) (interface{}, error) {
    intVal := val.(int)
    return intVal + 1, err
}

func Throwable(val interface{}, err error) (interface{}, error) {
    return 0, fmt.Errorf("Exception occurred with supplied value: %+v", val)
}

func Test_IRunnable(t *testing.T) {
    fmt.Println("Starting thennable_test.go")
    
    addOne := thennable.New(AddOne)
    throw := thennable.New(Throwable)
    
    //Test 1
    result, err := thennable.Start(1).//starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(addOne). //result = 4
        Then(addOne). //result = 5
        End()
    
    fmt.Printf("Test 1 State: %+v Error:%+v\n", result, err)
    if result.(int) != 5 {
        t.Errorf("Test 1 result should be = 5, actual value: %d", result)
    }
    if err != nil {
        t.Errorf("Test 1 error should be = nil, actual value: %+v", err)
    }
    
    //Test 2
    result, err = thennable.Start(1).//starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //add 2 with 1 = 3
        Then(throw). //result = 0, err = exception
        Then(addOne). //skipped
        Then(addOne). //skipped
        End()
        
    fmt.Printf("Test 2 State: %+v Error:%+v\n", result, err)
    if result.(int) != 0 {
        t.Errorf("Test 2 result should be = 0, actual value: %d", result)
    }
    if err == nil {
        t.Errorf("Test 2 error should not be = nil, actual value: %+v", err)
    }
    
    //Test 3
    result, err = thennable.Start(1).//starts with 1
        BreakOnError(false). //keep propagating the function chain
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(throw). //result = 0, err = exception
        Then(addOne). //add 0 with 1 = 1, err = exception from previous step
        Then(addOne). //result = 2, err = exception from previous step
        End()
        
    fmt.Printf("Test 3 State: %+v Error:%+v\n", result, err)
    if result.(int) != 2 {
        t.Errorf("Test 3 result should be = 2, actual value: %d", result)
    }
    if err == nil {
        t.Errorf("Test 3 error should not be = nil, actual value: %+v", err)
    }
}
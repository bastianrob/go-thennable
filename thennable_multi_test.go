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
    fmt.Println("Starting thennable_multi_test.go")
    
    addOne := thennable.NewMulti(AddOneMulti)
    throw := thennable.NewMulti(ThrowableMulti)
    
    //Test 1
    result, err := thennable.StartMulti(1). //starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(addOne). //result = 4
        Then(addOne). //result = 5
        End()
    
    fmt.Printf("Test 1 State: %+v Error:%+v\n", result, err)
    if result[0].(int) != 5 {
        t.Errorf("Test 1 result should be = 5, actual value: %d", result)
    }
    if err != nil {
        t.Errorf("Test 1 error should be = nil, actual value: %+v", err)
    }
    
    //Test 2
    result, err = thennable.StartMulti(1). //starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //add 2 with 1 = 3
        Then(throw).  //result = nil, err = exception
        Then(addOne). //skipped
        Then(addOne). //skipped
        End()
        
    fmt.Printf("Test 2 State: %+v Error:%+v\n", result, err)
    if result != nil {
        t.Errorf("Test 2 result should be = nil, actual value: %d", result)
    }
    if err == nil {
        t.Errorf("Test 2 error should not be = nil, actual value: %+v", err)
    }
    
    //Test 3
    result, err = thennable.StartMulti(1). //starts with 1
        BreakOnError(false). //keep propagating the function chain
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(throw).  //result = nil, err = exception
        Supply(0).     //resupply value with 0
        Then(addOne). //add 0 with 1 = 1, err = exception from previous step
        Then(addOne). //result = 2, err = exception from previous step
        End()
        
    fmt.Printf("Test 3 State: %+v Error:%+v\n", result, err)
    if result[0].(int) != 2 {
        t.Errorf("Test 3 result should be = 2, actual value: %d", result)
    }
    if err != nil {
        t.Errorf("Test 3 error should be = nil, actual value: %+v", err)
    }
}
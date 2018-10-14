package thennable_test

import (
    "fmt"
    "strconv"
    "thennable"
    "testing"
)

func AddOne(val int) (int, error) {
    return val + 1, nil
}

func Throwable(val interface{}) (error) {
    return fmt.Errorf("Exception occurred with supplied value: %+v", val)
}

func LogError(err error) {
    fmt.Printf("Log Error: %+v\n", err)
}

func Test_IThennable(t *testing.T) {
    fmt.Println("\nStarting Test_IThennable.go")
    
    addOne := AddOne
    throw := Throwable
    
    //Test 1
    fmt.Println("Start Test 1: Expect result = 5 and err is nil")
    result, err := thennable.Start(1).//starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(addOne). //result = 4
        Then(addOne). //result = 5
        End()
    
    if result[0].(int) != 5 {
        t.Errorf("Test 1 result should be = 5, actual value: %d", result[0])
    }
    if err != nil {
        t.Errorf("Test 1 error should be = nil, actual value: %+v", err)
    }
    
    //Test 2
    fmt.Println("\nStart Test 2: Expect result is nil and err is not nil")
    result, err = thennable.Start(1).//starts with 1
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //add 2 with 1 = 3
        Then(throw). //err = exception
        Handle(LogError).
        Then(addOne). //skipped
        Then(addOne). //skipped
        End()
        
    if len(result) != 0 {
        t.Errorf("Test 2 result should be = empty, actual value: %d", len(result))
    }
    if err == nil {
        t.Errorf("Test 2 error should not be = nil, actual value: %+v", err)
    }
    
    //Test 3
    fmt.Println("\nStart Test 3: Expect result = 2 and err is nil")
    result, err = thennable.Start(1).//starts with 1
        BreakOnError(false). //keep propagating the function chain
        Then(addOne). //add 1 with 1 = 2
        Then(addOne). //result = 3
        Then(throw).  //result = 0, err = exception ignored
        Supply(0).    //resupply with 0
        Then(addOne). //add 0 with 1 = 1, err = nil
        Then(addOne). //result = 2, err = nil
        End()
        
    if result[0].(int) != 2 {
        t.Errorf("Test 3 result should be = 2, actual value: %d", result[0])
    }
    if err != nil {
        t.Errorf("Test 3 error should be = nil, actual value: %+v", err)
    }
    
    //Test 4
    fmt.Println("\nStart Test 4: Expect recover from error, result is 10 and err is nil")
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
        
    if result[0].(int) != 10 {
        t.Errorf("Test 4 result should be = 10, actual value: %d", result[0])
    }
    if err != nil {
        t.Errorf("Test 4 error should be nil, actual value: %+v", err)
    }
    
    //Test 5
    fmt.Println("\nStart Test 5: Expect result is '6' and err is nil")
    result, err = thennable.Start(1).
        Then(func (one int) (int, int, error) {
            return one, 2, nil
        }).
        Then(func (one, two int) (int, int, int, error) {
            return one, two, 3, nil
        }).
        Then(func (one, two, three int) (string, error) {
            return strconv.Itoa(one + two + three), nil
        }).
        End()
        
    if result[0].(string) != "6" {
        t.Errorf("Test 5 result should be = '6', actual value: %v\n", result)
    }
    if err != nil {
        t.Errorf("Test 5 error should be nil, actual value: %+v\n", err)
    }
}
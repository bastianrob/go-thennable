# go-thennable

With this package, you can chain a series of functions to build a pipeline

## Using thennable for error handling / control flow

Consider the following example

```go
func main() {
    resultOfDoingThis, err := DoThis()
    if err != nil {
        //log the error
        return
    }

    resultOfDoingThat err := DoThat(resultOfDoingThis)
    if err != nil {
        //log the error
        return
    }

    resultOfDoingAnother err := DoAnother(resultOfDoingThat)
    if err != nil {
        //log the error
        return
    }
}
```

The above example is the idiomatic way of error handling in go.
Some like it, but most hate it because we need to manually check and return from error one by one.

Now consider the following example using thennable

```go
func main() {
    result, err := thennable.Start().
        Then(DoThis).
        Then(DoThat).
        Then(DoAnother).
        End()
}
```

In the above example, `err` will either be from `DoThis`, `DoThat`, or `DoAnother` function.

* `DoThat` and `DoAnother` won't be executed if `DoThis` produce an error.
* `DoAnother` won't be executed if `DoThat` produce an error.
* Or you can bypass the error and keep running next function by setting `BreakOnError(true)` in the thennable pipeline. ()

To handle error can either write extra lines like:

```go
    if err != nil {
        //log the error
        return
    }
```

or better yet, you can use the built in `Handle` function so the code will looks like:

```go
func GlobalErrorHandler(err error) {
    log.Println(err)
}

func main() {
    result, err := thennable.Start().
        Then(DoThis).
        Then(DoThat).
        Then(DoAnother).
        Handle(GlobalErrorHandler)
        End()
}
```

## Example 1

```go
import (
    "fmt"
    thennable "github.com/bastianrob/go-thennable"
)

func AddOne(num int) (int, error) {
    return num + 1, nil
}

func Decide(num int) (string, error) {
    if num == 1 {
        return "one", nil
    else {
        return "not one", nil
    }
}

func main() {
    result, err := thennable.Start(0). //start with zero
        Then(AddOne). //0+1
        Then(Decide). //1 = one
        End()
    str := result[0].(string)
    fmt.Printf("Res: %v, Str: %s, Err: %v", result, str, err)
    //Res: [one], Str: one, Err: <nil>
}
```

## Example 2: Inlining

```go
import (
    "fmt"
    thennable "github.com/bastianrob/go-thennable"
)

func main() {
    result, err := thennable.Start(1). //start with 1
        Then(func (one int) (int, int, error) {
            return one, 2, nil
        }). //return 1, 2
        Then(func (one, two int) (int, int, int, error) {
            return one, two, 3, nil
        }) //return 1, 2, 3
        End()
    fmt.Printf("Res: %v, Str: %s, Err: %v", result, str, err)
    //Res: [1 2 3], Err: <nil>
}
```

## Example 3: Error Occurred

```go
import (
    "fmt"
    "errors"
    thennable "github.com/bastianrob/go-thennable"
)

func AddOne(num int) (int, error) {
    return num + 1, nil
}

func main() {
    result, err := thennable.Start(1). //start with 1
        Then(AddOne). //1 + 1
        Then(func (two int) (int, error) {
            return two, errors.New("Whooops, something happened")
        }). //return 2, error
        Then(AddOne). //skipped
        Then(AddOne). //skipped
        End()

    fmt.Printf("Res: %v, Str: %s, Err: %v", result, str, err)
    //Res: [2], Err: Whooops, something happened
}
```

## Example 4: Recover from error

```go
import (
    "fmt"
    "errors"
    thennable "github.com/bastianrob/go-thennable"
)

func AddOne(num int) (int, error) {
    return num + 1, nil
}

func main() {
    result, err := thennable.Start(1). //start with 1
        BreakOnError(false). //Pipeline keeps going when error happens
        Then(AddOne). //1 + 1
        Then(func (two int) error {
            return errors.New("Whooops, something happened")
        }).           //return error
        Supply(8)      //resuply the pipeline with 8
        Then(AddOne). //8 + 1
        Then(AddOne). //9 + 1
        End()

    fmt.Printf("Res: %v, Str: %s, Err: %v", result, str, err)
    //Res: [10], Err: nil
}
```

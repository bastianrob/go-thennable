# go-thennable

With this package, you can chain a series of functions to build a pipeline

### Example 1

```
import (
	"fmt"
	"thennable"
)

func AddOne(num int) (int, error) {
	return num + 1, nil
}

func Decide(num int) (string, error) {
	if num == 1 {
    	return "one"
    else {
    	return "not one"
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


### Example 2: Inlining

```
import (
	"fmt"
	"thennable"
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

### Example 3: Error Occurred
```
import (
	"fmt"
    "errors"
	"thennable"
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


### Example 4: Recover from error
```
import (
	"fmt"
    "errors"
	"thennable"
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
        }). 		  //return error
        Supply(8)	  //resuply the pipeline with 8
        Then(AddOne). //8 + 1
        Then(AddOne). //9 + 1
        End()
        
    fmt.Printf("Res: %v, Str: %s, Err: %v", result, str, err)
    //Res: [10], Err: nil
}
```






















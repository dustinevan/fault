### Problem: 
Given the discussion layed out [here](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully), there are times when a middle teir function is the best place to know what the system should do with an error. 

### Solution:
fault shows how to wrap errors with new opaque functionality without breaking err.Cause() or stack tracing. For my use case I've added the ability to add HTTP status codes and an alert flag. Because `httpStatus` and `alert` are opaque errors, like those in github.com/pkg/errors, as long as they are rewrapped with an error that implements the causer interface--`errors.Wrap()`--they can be inspected later. The causer interface is described in github.com/pkg/errors 
```golang
type causer interface {
    Cause() error
}
``` 

Suggesting HTTP status codes and adding an Alert flag work like this:  
```golang
func GetData(key string) (*stuff, error) {
    ... database query ...
    if err != nil {
      err = errors.Wrap(err, "db connection gone")
      err = fault.WithHttpStatus(err, http.StatusInternalServerError)
      return nil, err
    ... 
 }
 
 ...
 func GetDataService(key string) (*stuff, error) {
    s, err := GetData(key)
    if err != nil {
      err = fault.WithAlert(err)
      return nil, errors.Wrap(err, "service failed to get the data")
    }
    ...
 }
 
 func main() {
    ...
    s, err GetDataService(key)
    if err != nil {
      if status, ok := fault.HttpStatus(err); ok {
        fmt.Println(status)
      }
      if fault.IsAlert(err) {
        fmt.Println("Alert!!")
      }
    }
 }
 ```
 
 

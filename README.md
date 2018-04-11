### Problem: 
Given the discussion layed out here[https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully], there are times when a middle teir function is the best place to know what the system should do with an error. 

### Solution:
fault shows how to wrap errors with new opaque functionality without breaking errors.Cause() or stack tracing. For my use case I've added the ability to add HTTP status codes and an alert flag. Because httpStatus and alert are opaque errors, like those in github.com/pkg/errors, they can be rewrapped with any error that implements the causer interface described in github.com/pkg/errors. 
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
      
      return nil, err
    ... 
 }
 
 ...
 func GetDataService(key string) (*stuff, error) {
    s, err := GetData(key)
    if err != nil {
      err = fault.WithHttpStatus(err, http.StatusInternalServerError)
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
 
 

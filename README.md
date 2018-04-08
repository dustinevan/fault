# fault
fault builds on the ideas in github.com/pkg/errors to allow systems that use errors.Wrap() to add alert flags or http status codes that can be discovered later on. The means low-level code can suggest http statuses and flag errors for alert, service code can wrap errors as normal, then top-level code can inspect the errors for http statuses or alert flags. 

```golang
func GetData(key string) (*stuff, error) {
    ... database query ...
    if err != nil {
      err = errors.Wrap(err, "db connection gone")
      err = fault.WithHttpStatus(err, http.StatusInternalServerError)
      err = fault.WithAlert(err)
      return nil, err
    ... 
 }
 
 ...
 func GetDataService(key string) (*stuff, error) {
    s, err := GetData(key)
    if err != nil {
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
 
 

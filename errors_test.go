package fault

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestWithHttpStatus(t *testing.T) {
	cause := errors.New("db connection error")
	wrap1 := errors.Wrap(cause, "1st wrap")
	withHttp := WithHttpStatus(wrap1, http.StatusInternalServerError)
	wrap2 := errors.Wrap(withHttp, "2nd wrap")
	wrap3 := errors.Wrap(wrap2, "3rd wrap")

	assert.Equal(t, cause, errors.Cause(wrap3))
	status, ok := HttpStatus(wrap3)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, status)

	assert.Equal(t, cause, errors.Cause(withHttp))
	status, ok = HttpStatus(wrap1)
	assert.False(t, ok)
	assert.Equal(t, 0, status)

	assert.Equal(t, "3rd wrap: 2nd wrap: http status 500: 1st wrap: db connection error", wrap3.Error())

	status, ok = HttpStatus(nil)
	assert.False(t, ok)
	assert.Equal(t, 0, status)


	stack := fmt.Sprintf("%+v", wrap2)
	assert.Contains(t, stack, "db connection error", "db connection error msg not found")
	assert.Contains(t, stack, "fault/errors_test.go:13", "no line trace for db connection error")
	assert.Contains(t, stack, "1st wrap", "1st wrap msg not found")
	assert.Contains(t, stack, "fault/errors_test.go:14", "no line trace for 1st wrap")
	assert.Contains(t, stack, "http status 500", "http status 500 msg not found")
	assert.Contains(t, stack, "fault/errors_test.go:15", "no line trace for http status 500")
	assert.Contains(t, stack, "2nd wrap", "2nd wrap msg not found")
	assert.Contains(t, stack, "fault/errors_test.go:16", "no line trace for 2nd wrap")

}

func TestWithAlert(t *testing.T) {
	cause := errors.New("db connection error")
	wrap1 := errors.Wrap(cause, "1st wrap")
	withAlert := WithAlert(wrap1)
	wrap2 := errors.Wrap(withAlert, "2nd wrap")
	wrap3 := errors.Wrap(wrap2, "3rd wrap")

	assert.Equal(t, cause, errors.Cause(wrap3))
	ok := IsAlert(wrap3)
	assert.True(t, ok)

	assert.Equal(t, cause, errors.Cause(withAlert))
	ok = IsAlert(wrap1)
	assert.False(t, ok)

	assert.Equal(t, "3rd wrap: 2nd wrap: alert: 1st wrap: db connection error", wrap3.Error())

	assert.Equal(t, cause, errors.Cause(withAlert))
	ok = IsAlert(nil)
	assert.False(t, ok)

	stack := fmt.Sprintf("%+v", wrap2)
	assert.Contains(t, stack, "db connection error", "db connection error msg not found")
	assert.Contains(t, stack, "fault/errors_test.go:49", "no line trace for db connection error")
	assert.Contains(t, stack, "1st wrap", "1st wrap msg not found")
	assert.Contains(t, stack, "fault/errors_test.go:50", "no line trace for 1st wrap")
	assert.Contains(t, stack, "alert", "alert msg not found")
	assert.Contains(t, stack, "fault/errors_test.go:51", "no line trace for http status 500")
	assert.Contains(t, stack, "2nd wrap", "2nd wrap msg not found")
	assert.Contains(t, stack, "fault/errors_test.go:52", "no line trace for 2nd wrap")
}

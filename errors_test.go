package fault

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "3rd wrap: 2nd wrap: 1st wrap: db connection error", wrap3.Error())

	status, ok = HttpStatus(nil)
	assert.False(t, ok)
	assert.Equal(t, 0, status)
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

	assert.Equal(t, "3rd wrap: 2nd wrap: 1st wrap: db connection error", wrap3.Error())

	assert.Equal(t, cause, errors.Cause(withAlert))
	ok = IsAlert(nil)
	assert.False(t, ok)

}
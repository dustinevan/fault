package fault

import (
	"fmt"

	"github.com/pkg/errors"
)

// This interface allows use of the Cause method via type assertion. errors created by
// errors.Wrap, errors.Wrapf, errors.WithMessage, or errors.WithStack implement this interface
type causer interface {
	Cause() error
}

// This interface allows use of StackTrace method via type assertion. errors created by
// errors.Wrap, errors.Wrapf, errors.WithStack implement this interface
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// httpStatus is an error that has a recommended http status associated with it
type httpStatuser interface {
	HttpStatus() int
	causer
}

type httpStatus struct {
	status int
	cause  error
}

func WithHttpStatus(err error, status int) error {
	if err == nil {
		return nil
	}

	return &httpStatus{
		status: status,
		cause:  errors.Wrap(err, fmt.Sprintf("http status %v", status)),
	}
}

func (h *httpStatus) HttpStatus() int {
	return h.status
}

func (h *httpStatus) Cause() error { return h.cause }

func (h *httpStatus) Error() string {
	return h.cause.Error()
}

func (h *httpStatus) Format(s fmt.State, verb rune) {
	if err, ok := h.cause.(fmt.Formatter); ok {
		err.Format(s, verb)
		return
	}
	panic(h.cause.Error() + " does not implement the fmt.Formatter interface.")
}

func HttpStatus(err error) (int, bool) {
	for err != nil {
		c, ok := err.(causer)
		if !ok {
			return 0, false
		}
		h, ok := c.(httpStatuser)
		if !ok {
			err = c.Cause()
			continue
		}
		return h.HttpStatus(), true
	}
	return 0, false
}

// if an alerter exists in the error tree, the entire error should be sent as an alert.
type alerter interface {
	Alert()
	causer
}

type alertErr struct {
	cause error
}

func WithAlert(err error) error {
	if err == nil {
		return nil
	}
	err = errors.Wrap(err, "alert")

	return &alertErr{
		cause: err,
	}
}

func (a *alertErr) Error() string {
	return a.cause.Error()
}

func (a *alertErr) Cause() error {
	return a.cause
}

func (a *alertErr) Alert() {
}

func (a *alertErr) Format(s fmt.State, verb rune) {
	if err, ok := a.cause.(fmt.Formatter); ok {
		err.Format(s, verb)
		return
	}
	panic(a.cause.Error() + " does not implement the fmt.Formatter interface.")
}

func IsAlert(err error) bool {
	for err != nil {
		c, ok := err.(causer)
		if !ok {
			return false
		}
		_, ok = c.(alerter)
		if !ok {
			err = c.Cause()
			continue
		}
		return true
	}
	return false
}

type ErrCode interface {
	System() string
	Code() int
	Description() string
}

type errCoder interface {
	ErrCode() ErrCode
}

type errCode struct {
	cause error
	ec    ErrCode
}

func WithErrCode(err error, code ErrCode) *errCode {
	return &errCode{
		cause: err,
		ec:    code,
	}
}

func (c *errCode) Error() string {
	return c.cause.Error()
}

func (c *errCode) Cause() error {
	return c.cause
}

func (c *errCode) ErrCode() ErrCode {
	return c.ec
}

func (c *errCode) Format(s fmt.State, verb rune) {
	if err, ok := c.cause.(fmt.Formatter); ok {
		err.Format(s, verb)
		return
	}
	panic(c.cause.Error() + " does not implement the fmt.Formatter interface.")
}

func HasErrCode(err error) (ErrCode, bool) {
	for err != nil {
		c, ok := err.(causer)
		if !ok {
			return nil, false
		}
		e, ok := c.(errCoder)
		if !ok {
			err = c.Cause()
			continue
		}
		return e.ErrCode(), true
	}
	return nil, false
}
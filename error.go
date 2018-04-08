package fault

import (
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

// if an alerter exists in the error tree, the entire error should be sent as an alert.
type alerter interface {
	Alert()
	causer
}

type httpStatus struct {
	status int
	err    error
}

func (h httpStatus) HttpStatus() int {
	return h.status
}

func (h httpStatus) Cause() error {
	return h.err
}

func (h httpStatus) Error() string {
	return h.err.Error()
}

func WithHttpStatus(err error, status int) error {
	return httpStatus{
		status: status,
		err:    err,
	}
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

type alertErr struct {
	err error
}

func WithAlert(err error) error {
	return &alertErr{
		err: err,
	}
}

func (a alertErr) Error() string {
	return a.err.Error()
}

func (a alertErr) Cause() error {
	return a.err
}

func (a alertErr) Alert() {
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

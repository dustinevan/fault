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
	error
}

type tracePrinter interface {
	Trace() string
}

// Interfaces that inherit cause are discoverable during the cause chain loop. See laag.HttpCode(err)
type httpCodeError interface {
	HttpCode() int
	causer
	error
}

type tagger interface {
	Tag() Tag
	causer
}

func ErrorWithTrace(err error) string {
	return err.Error() + "\n" + Trace(err)
}

func Trace(err error) string {
	trace := ""
	for err != nil {
		c, ok := err.(causer)
		if !ok {
			t, ok := err.(tracePrinter)
			if !ok {
				return trace
			}
			return t.Trace()
		}
		t, ok := err.(tracePrinter)
		if !ok {
			err = c.Cause()
			continue
		}
		trace = t.Trace()
		err = c.Cause()
	}
	return trace
}

func HttpCode(err error) (int, bool) {
	for err != nil {
		c, ok := err.(causer)
		if !ok {
			return 0, false
		}
		h, ok := c.(httpCodeError)
		if !ok {
			err = c.Cause()
			continue
		}
		return h.HttpCode(), true
	}
	return 0, false
}

func LogTag(err error) Tag {
	tag := Err
	for err != nil {
		t, ok := err.(tagger)
		if !ok {
			c, ok := err.(causer)
			if !ok {
				return tag
			}
			err = c.Cause()
			continue
		}
		if t.Tag() != tag {
			tag = t.Tag()
		}
		if tag == NoLog {
			return NoLog
		}
		c, ok := err.(causer)
		if !ok {
			return tag
		}
		err = c.Cause()
	}
	return tag
}

// func SubSystem(err error) (Subsystem, bool)

type Tag int

const (
	NoLog      Tag = iota
	Success
	Info
	Err
	Bug
	SysFailure
)

var tags = [...]string{
	"no log needed",
	"success",
	"info",
	"error",
	"bug",
	"system failure",
}

func (l Tag) String() string {
	return tags[int(l)]
}

type appError struct {
	err       stackTracer
	tag       Tag
	subsystem string
}

func Error(err error, subsystem string, tag Tag) error {
	st, ok := err.(stackTracer)
	if !ok {
		// callers should wrap the error so that the stacktrace works correctly
		e := errors.WithStack(err)
		return &appError{
			err:       e.(stackTracer),
			tag:       tag,
			subsystem: subsystem,
		}

	}
	return &appError{
		err:       st,
		tag:       tag,
		subsystem: subsystem,
	}
}

func (a *appError) Error() string {
	return fmt.Sprintf("%s: %s: %s", a.subsystem, a.tag, a.err.Error())
}

func (a *appError) Trace() string {
	s := ""
	for _, f := range a.err.StackTrace() {
		s += fmt.Sprintf("%+v:", f)
	}
	return s
}

func (a *appError) Cause() error {
	return a.err
}

func (a *appError) Tag() Tag {
	return a.tag
}

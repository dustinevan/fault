package fault

import (
	"fmt"

	"github.com/pkg/errors"
)

type webRequestError struct {
	err      stackTracer
	httpcode int
	tag      Tag
}

func WebRequestErr(err error, httpcode int, tag Tag) error {
	st, ok := err.(stackTracer)
	if !ok {
		// callers should wrap the error so that the stacktrace works correctly
		e := errors.WithStack(err)
		return &webRequestError{
			err:      e.(stackTracer),
			httpcode: httpcode,
			tag:      tag,
		}
	}
	return &webRequestError{
		err:      st,
		httpcode: httpcode,
		tag:      tag,
	}
}

func (w *webRequestError) HttpCode() int {
	return w.httpcode
}

func (w *webRequestError) Error() string {
	return fmt.Sprintf("%s: httpcode: %v ", w.err.Error(), w.httpcode)
}

func (w *webRequestError) Trace() string {
	s := ""
	for _, f := range w.err.StackTrace() {
		s += fmt.Sprintf("%+v:", f)
	}
	return s
}

func (w *webRequestError) Cause() error {
	return w.err
}

func (w *webRequestError) Tag() Tag {
	return w.tag
}

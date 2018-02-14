package fault

import (
	"fmt"

	"github.com/pkg/errors"
)

type webRequestError struct {
	err       stackTracer
	httpcode  int
	tag       LogTag
	subsystem string
}

func WebRequestErr(err error, subsystem string, tag LogTag, httpcode int) error {
	st, ok := err.(stackTracer)
	if !ok {
		// callers should wrap the error so that the stacktrace works correctly
		e := errors.WithStack(err)
		return &webRequestError{
			err:       e.(stackTracer),
			httpcode:  httpcode,
			tag:       tag,
			subsystem: subsystem,
		}
	}
	return &webRequestError{
		err:       st,
		httpcode:  httpcode,
		tag:       tag,
		subsystem: subsystem,
	}
}

func (w *webRequestError) HttpCode() int {
	return w.httpcode
}

func (w *webRequestError) Error() string {
	return fmt.Sprintf("%s: %s: %s httpcode: %v ", w.subsystem, w.tag, w.err.Error(), w.httpcode)
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


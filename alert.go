package fault

import (
	"fmt"

	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type alertError struct {
	err       stackTracer
	Tag       Tag
	Subsystem string
	Msg       string
}

func AlertErr(err error, subsystem string, tag Tag) error {
	st, ok := err.(stackTracer)
	if !ok {
		// callers should wrap the error so that the stacktrace works correctly
		e := errors.WithStack(err)
		return &alertError{
			err:       e.(stackTracer),
			Tag:       tag,
			Subsystem: subsystem,
		}

	}
	return &alertError{
		err:       st,
		Tag:       tag,
		Subsystem: subsystem,
	}
}

const OpeningAlertTag = "\n<alert>\n"
const ClosingAlertTag = "\n</alert>\n"

func (a *alertError) Error() string {
	a.Msg = a.err.Error() + Trace(a.err)
	return OpeningAlertTag + a.jsonstring() + ClosingAlertTag + fmt.Sprintf("%s: %s: %s", a.Subsystem, a.Tag, a.err.Error())
}

func (a *alertError) Trace() string {
	s := ""
	for _, f := range a.err.StackTrace() {
		s += fmt.Sprintf("%+v:", f)
	}
	return s
}

func (a *alertError) jsonstring() string {
	b, err := json.Marshal(a)
	if err != nil {
		return "unable to marshal alert err: " + a.Msg + ": " + err.Error()
	}
	return string(b)
}

func (a *alertError) Cause() error {
	return a.err
}

func (a *alertError) LogTag() Tag {
	return a.Tag
}
package fault

import (
	"fmt"

	"github.com/pkg/errors"
	"strconv"
	"github.com/json-iterator/go"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

// Interfaces that inherit cause are discoverable during the cause chain loop. See laag.HttpCode(err)
type httpCodeError interface {
	HttpCode() (int, bool)
	causer
	error
}

func GetHttpCode(err error) (int, bool) {
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
		return h.HttpCode()
	}
	return 0, false
}

type typ int

const (
	NoLog typ = iota
	Success
	Info
	Err
	Bug
	SysFailure
)

var typs = [...]string{
	"no log needed",
	"success",
	"info",
	"error",
	"bug",
	"system failure",
}

func (t typ) String() string {
	return typs[int(t)]
}

func (t *typ) UnmarshalJSON(b []byte) error {
	for i, p := range typs {
		if p == string(b) {
			*t = typ(i)
			return nil
		}
	}
	*t = NoLog
	return nil
}

func (t *typ) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}
const OpeningAlertTag = "\n<alert>\n"
const ClosingAlertTag = "\n</alert>\n"

type AlertTag struct {
	Msg       string
	Type       typ
	Subsystem string
}

type ErrOption func(*appError)

func WithTrace() ErrOption {
	return func(a *appError) {
		a.includeTrace = true
	}
}

func Alert(subsystem string, tag typ) ErrOption {
	return func(a *appError) {
		a.isAlert = true
		a.subsystem = subsystem
		a.tag = tag
	}
}

func Type(t typ) ErrOption {
	return func(a *appError) {
		a.tag = t
	}
}

func System(s string) ErrOption {
	return func(a *appError) {
		a.subsystem = s
	}
}

func HttpCode(h int) ErrOption {
	return func(a *appError) {
		a.httpcode = h
	}
}

func HttpUnauthorized() ErrOption {
	return HttpCode(http.StatusUnauthorized)
}

func HttpServerError() ErrOption {
	return HttpCode(http.StatusInternalServerError)
}

func HttpBadRequest() ErrOption {
	return HttpCode(http.StatusBadRequest)
}

func HttpNotFound() ErrOption {
	return HttpCode(http.StatusNotFound)
}

func DontLog() ErrOption {
	return Type(NoLog)
}

func LogSuccess() ErrOption {
	return Type(Success)
}
func LogInfo() ErrOption {
	return Type(Info)
}
func LogError() ErrOption {
	return Type(Err)
}
func LogBug() ErrOption {
	return Type(Bug)
}

func LogFailure() ErrOption {
	return Type(SysFailure)
}

type appError struct {
	err          stackTracer
	tag          typ
	subsystem    string
	httpcode     int
	includeTrace bool
	isAlert      bool
}

func Error(err error, opts ...ErrOption) error {
	st, ok := err.(stackTracer)
	if !ok {
		e := errors.WithStack(err)
		a := &appError{
			err: e.(stackTracer),
			tag: Err,
		}
		for _, opt := range opts {
			opt(a)
		}
		return a
	}
	a := &appError{
		err: st,
		tag: Err,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *appError) Error() string {
	e := a.builderror()
	if a.isAlert == true {
		b, err := json.Marshal(a.Alert())
		if err != nil {
			return e + ": unable to marshal alert: " + err.Error()
		}
		return e + OpeningAlertTag + string(b) + ClosingAlertTag
	}
	return e
}

func (a *appError) builderror() string {
	if a.tag == NoLog {
		return ""
	}
	e := fmt.Sprintf("%s: %s", a.tag, a.err.Error())
	if a.subsystem != "" {
		e = a.subsystem + ": " + e
	}
	if a.httpcode != 0 {
		e += ": http code " + strconv.Itoa(a.httpcode)
	}
	if a.includeTrace == true {
		e += "\n" + a.Trace()
	}
	return e
}

func (a *appError) Trace() string {
	s := ""
	for _, f := range a.err.StackTrace() {
		s += fmt.Sprintf("%+v:", f)
	}
	return s
}

func (a *appError) Alert() AlertTag {
	return AlertTag{
		Msg: a.builderror(),
		Type: a.tag,
		Subsystem: a.subsystem,
	}
}

func (a *appError) Cause() error {
	return a.err
}

func (a *appError) HttpCode() (int, bool) {
	if a.httpcode == 0 {
		return 0, false
	}
	return a.httpcode, true
}
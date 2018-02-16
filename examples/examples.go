package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/dustinevan/fault"
	"github.com/pkg/errors"
)

func main() {
	orig := errors.Wrap(fmt.Errorf("db error"), "found this")
	web := fault.WebRequestErr(orig, http.StatusInternalServerError, fault.Bug)
	err := errors.Wrap(web, "saw this error")
	err = errors.Wrap(err, "saw this error again")

	alert := fault.AlertErr(err, "database", fault.SysFailure)

	faulterr := fault.Error(alert, "email", fault.NoLog)

	log.Println(errors.Cause(orig))
	log.Println(errors.Cause(web))
	log.Println(errors.Cause(err))
	log.Println(errors.Cause(alert))
	log.Println(errors.Cause(faulterr))

	log.Println(fault.Trace(orig))
	log.Println(fault.Trace(web))
	log.Println(fault.Trace(err))
	log.Println(fault.Trace(alert))
	log.Println(fault.Trace(faulterr))

	log.Println(fault.HttpCode(orig))
	log.Println(fault.HttpCode(web))
	log.Println(fault.HttpCode(err))
	log.Println(fault.HttpCode(alert))
	log.Println(fault.HttpCode(faulterr))

	log.Println(fault.ErrorWithTrace(orig))
	log.Println(fault.ErrorWithTrace(web))
	log.Println(fault.ErrorWithTrace(err))
	log.Println(fault.ErrorWithTrace(alert))
	log.Println(fault.ErrorWithTrace(faulterr))

	log.Println(orig)
	log.Println(web)
	log.Println(err)
	log.Println(alert)
	log.Println(faulterr)

	e0 := fault.WebRequestErr(fmt.Errorf("something happened"), http.StatusInternalServerError, fault.NoLog)

	fmt.Println(e0)

	log.Println(fault.LogTag(orig))
	log.Println(fault.LogTag(web))
	log.Println(fault.LogTag(err))
	log.Println(fault.LogTag(alert))
	log.Println(fault.LogTag(faulterr))
	log.Println(fault.LogTag(e0))

}

/*
2018/02/15 20:58:36 db error
2018/02/15 20:58:36 db error
2018/02/15 20:58:36 db error
2018/02/15 20:58:36 db error
2018/02/15 20:58:36 db error
2018/02/15 20:58:36
2018/02/15 20:58:36 main.main
        /Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/15 20:58:36 main.main
        /Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/15 20:58:36 main.main
        /Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/15 20:58:36 main.main
        /Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/15 20:58:36 0 false
2018/02/15 20:58:36 500 true
2018/02/15 20:58:36 500 true
2018/02/15 20:58:36 500 true
2018/02/15 20:58:36 500 true
2018/02/15 20:58:36 found this: db error

2018/02/15 20:58:36 found this: db error: httpcode: 500
main.main
        /Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/15 20:58:36 saw this error again: saw this error: found this: db error: httpcode: 500
main.main
        /Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/15 20:58:36
<alert>
{"Tag":5,"Subsystem":"database","Msg":"saw this error again: saw this error: found this: db error: httpcode: 500 main.main\n\t/Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:"}
</alert>
database: system failure: saw this error again: saw this error: found this: db error: httpcode: 500
main.main
        /Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/15 20:58:36 email: no log needed:
<alert>
{"Tag":5,"Subsystem":"database","Msg":"saw this error again: saw this error: found this: db error: httpcode: 500 main.main\n\t/Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:"}
</alert>
database: system failure: saw this error again: saw this error: found this: db error: httpcode: 500
main.main
        /Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
        /usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/15 20:58:36 found this: db error
2018/02/15 20:58:36 found this: db error: httpcode: 500
2018/02/15 20:58:36 saw this error again: saw this error: found this: db error: httpcode: 500
2018/02/15 20:58:36
<alert>
{"Tag":5,"Subsystem":"database","Msg":"saw this error again: saw this error: found this: db error: httpcode: 500 main.main\n\t/Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:"}
</alert>
database: system failure: saw this error again: saw this error: found this: db error: httpcode: 500
2018/02/15 20:58:36 email: no log needed:
<alert>
{"Tag":5,"Subsystem":"database","Msg":"saw this error again: saw this error: found this: db error: httpcode: 500 main.main\n\t/Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:"}
</alert>
database: system failure: saw this error again: saw this error: found this: db error: httpcode: 500
something happened: httpcode: 500
2018/02/15 20:58:36 error
2018/02/15 20:58:36 bug
2018/02/15 20:58:36 bug
2018/02/15 20:58:36 bug
2018/02/15 20:58:36 no log needed
2018/02/15 20:58:36 no log needed
*/

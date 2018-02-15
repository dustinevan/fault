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
	web := fault.WebRequestErr(orig, "auth", fault.Err, http.StatusInternalServerError)
	err := errors.Wrap(web, "saw this error")
	err = errors.Wrap(err, "saw this error again")

	alert := fault.AlertErr(err, "database", fault.SysFailure)

	faulterr := fault.Error(alert, "email", fault.Err)

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

}

/*
2018/02/14 13:39:56 db error
2018/02/14 13:39:56 db error
2018/02/14 13:39:56 db error
2018/02/14 13:39:56 db error
2018/02/14 13:39:56 db error
2018/02/14 13:39:56
2018/02/14 13:39:56 main.main
	/Users/dustinevan/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/14 13:39:56 main.main
	/Users/dustinevan/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/14 13:39:56 main.main
	/Users/dustinevan/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/14 13:39:56 main.main
	/Users/dustinevan/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/14 13:39:56 0 false
2018/02/14 13:39:56 500 true
2018/02/14 13:39:56 500 true
2018/02/14 13:39:56 500 true
2018/02/14 13:39:56 500 true
2018/02/14 13:39:56 found this: db error

2018/02/14 13:39:56 auth: error: found this: db error httpcode: 500
main.main
	/Users/dustinevan/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/14 13:39:56 saw this error again: saw this error: auth: error: found this: db error httpcode: 500
main.main
	/Users/dustinevan/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/14 13:39:56
<alert>
{"Tag":5,"Subsystem":"database","Msg":"saw this error again: saw this error: auth: error: found this: db error httpcode: 500 main.main\n\t/Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:"}
</alert>
database: system failure: saw this error again: saw this error: auth: error: found this: db error httpcode: 500
main.main
	/Users/dustinevan/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/14 13:39:56 email: error:
<alert>
{"Tag":5,"Subsystem":"database","Msg":"saw this error again: saw this error: auth: error: found this: db error httpcode: 500 main.main\n\t/Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:"}
</alert>
database: system failure: saw this error again: saw this error: auth: error: found this: db error httpcode: 500
main.main
	/Users/dustinevan/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit
	/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:
2018/02/14 13:39:56 found this: db error
2018/02/14 13:39:56 auth: error: found this: db error httpcode: 500
2018/02/14 13:39:56 saw this error again: saw this error: auth: error: found this: db error httpcode: 500
2018/02/14 13:39:56
<alert>
{"Tag":5,"Subsystem":"database","Msg":"saw this error again: saw this error: auth: error: found this: db error httpcode: 500 main.main\n\t/Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:"}
</alert>
database: system failure: saw this error again: saw this error: auth: error: found this: db error httpcode: 500
2018/02/14 13:39:56 email: error:
<alert>
{"Tag":5,"Subsystem":"database","Msg":"saw this error again: saw this error: auth: error: found this: db error httpcode: 500 main.main\n\t/Users/dustincurrie/code/go/src/github.com/dustinevan/fault/examples/examples.go:14:runtime.main\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/proc.go:195:runtime.goexit\n\t/usr/local/Cellar/go/1.9.2/libexec/src/runtime/asm_amd64.s:2337:"}
</alert>
database: system failure: saw this error again: saw this error: auth: error: found this: db error httpcode: 500

 */
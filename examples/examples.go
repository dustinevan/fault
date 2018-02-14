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

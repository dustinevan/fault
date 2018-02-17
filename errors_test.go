package fault

import (
	"testing"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func TestHttpCode(t *testing.T) {
	orig := errors.Wrap(fmt.Errorf("db error"), "found this")
	web := Error(orig, HttpCode(http.StatusInternalServerError), WithTrace())
	//fmt.Println(web)
	//fmt.Println(GetHttpCode(web))
	alert := Error(web, Alert("Database", SysFailure))
	//fmt.Println(alert)
	trace := Error(alert, WithTrace())
	msg := trace.Error()
	//fmt.Println(msg)
	i := strings.Index(msg, OpeningAlertTag)
	j := strings.Index(msg, ClosingAlertTag)
	alertjson := strings.Trim(msg[i+len(OpeningAlertTag)-1:j], "\n")
	var a AlertTag
	err := json.Unmarshal([]byte(alertjson), &a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a.Subsystem, "|", a.Type, "|", a.Msg)
}
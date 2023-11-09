package logic

import (
	"devops/utils/tools"
	"fmt"
)

var (
	ReqAssertErr = tools.NewRspError(tools.SystemErr, fmt.Errorf("请求异常"))
	Dms          = &DmsLogic{}
)

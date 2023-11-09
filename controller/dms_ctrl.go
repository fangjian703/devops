package controller

import (
	"bytes"
	"devops/common"
	"devops/logic"
	"devops/model/request"
	"devops/utils/tools"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DmsController struct {
}

func (d DmsController) GetDmsReqBody(c *gin.Context) {
	var DRB request.DmsReqBody
	err := c.ShouldBind(&DRB)
	if err != nil {
		tools.Err(c, tools.NewRspError(500, err), err)
	}
	respData := map[string]interface{}{"msg": "ok", "code": 200}
	if DRB.Event == "待审批" {
		respData = logic.Dms.DealDmsBody(&DRB)
	}
	tools.Success(c, respData)
}

func (d DmsController) OwnerDmsReq(c *gin.Context) {
	var DRB request.DmsReqBody
	err := c.ShouldBind(&DRB)
	if err != nil {
		tools.Err(c, tools.NewRspError(500, err), err)
	}
	respData := map[string]interface{}{"msg": "ok", "code": 200}
	if DRB.Event == "待审批" {
		err = logic.Dms.DealOwnerDmsBody(&DRB)
		if err != nil {
			common.Log.Error(err)
			tools.Err(c, tools.NewRspError(500, err), err)
		}
	}
	tools.Success(c, respData)
}

func (d DmsController) GetFeiShuReqBody(c *gin.Context) {

	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	var FRB request.FsCardBody
	err := c.ShouldBind(&FRB)
	if err != nil {
		common.Log.Error(err)
		tools.Err(c, tools.NewRspError(500, err), err)
	}

	if FRB.Category == "工单" && FRB.Event == "待审批" && FRB.Tag == "button" {
		err := logic.Dms.DmsApprovalOrder(&FRB)

		if err != nil {
			common.Log.Error(err)
			tools.Err(c, tools.NewRspError(500, err), err)
		}
		msg := "success"
		tools.Success(c, msg)
	}

	if FRB.Category == "工单" && FRB.Event == "审批通过" && FRB.Tag == "button" {

	}

	if FRB.Category == "工单" && FRB.Event == "审批拒绝" && FRB.Tag == "button" {

	}

	c.JSON(http.StatusOK, gin.H{
		"challenge": FRB.Challenge,
	})
}

package tools

import (
	"devops/common"
	"devops/config"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dms_enterprise20181101 "github.com/alibabacloud-go/dms-enterprise-20181101/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type AliYunCli struct {
}

type OrderBaseInfo struct {
	OrderId            *int64
	StatusCode         *string
	WorkflowStatusDesc *string
	StatusDesc         *string
}

func (a AliYunCli) CreateClient() (_cli *dms_enterprise20181101.Client, _err error) {
	Conf := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(config.Conf.AliYun.AK),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(config.Conf.AliYun.SK),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/dms-enterprise
	Conf.Endpoint = tea.String("dms-enterprise.cn-hangzhou.aliyuncs.com")
	_cli = &dms_enterprise20181101.Client{}
	_cli, _err = dms_enterprise20181101.NewClient(Conf)
	return _cli, _err
}

func (a AliYunCli) SubmitOrderApprovalRequest(taskID int64) error {
	client, _err := a.CreateClient()
	if _err != nil {
		common.Log.Error(_err)
		return _err
	}
	submitOrderApprovalRequest := &dms_enterprise20181101.SubmitOrderApprovalRequest{
		OrderId: tea.Int64(taskID),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.SubmitOrderApprovalWithOptions(submitOrderApprovalRequest, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()
	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

func (a AliYunCli) getWorkflowInstanceId(taskID int64) (*int64, error) {
	var WorkflowInstanceId *int64
	client, _err := a.CreateClient()
	if _err != nil {
		return WorkflowInstanceId, _err
	}

	getOrderBaseInfoRequest := &dms_enterprise20181101.GetOrderBaseInfoRequest{
		OrderId: tea.Int64(taskID),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		respBody, _err := client.GetOrderBaseInfoWithOptions(getOrderBaseInfoRequest, runtime)
		if _err != nil {
			return _err
		}
		WorkflowInstanceId = respBody.Body.OrderBaseInfo.WorkflowInstanceId
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return WorkflowInstanceId, _err
		}
	}
	return WorkflowInstanceId, _err
}

func (a AliYunCli) GetWorkflowInfo(taskID int64) (OrderBaseInfo, error) {
	var orderBaseInfo OrderBaseInfo
	client, _err := a.CreateClient()
	if _err != nil {
		return orderBaseInfo, _err
	}

	getOrderBaseInfoRequest := &dms_enterprise20181101.GetOrderBaseInfoRequest{
		OrderId: tea.Int64(taskID),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		respBody, _err := client.GetOrderBaseInfoWithOptions(getOrderBaseInfoRequest, runtime)
		if _err != nil {
			return _err
		}
		orderBaseInfo.OrderId = respBody.Body.OrderBaseInfo.OrderId
		orderBaseInfo.StatusCode = respBody.Body.OrderBaseInfo.StatusCode
		orderBaseInfo.StatusDesc = respBody.Body.OrderBaseInfo.StatusDesc
		orderBaseInfo.WorkflowStatusDesc = respBody.Body.OrderBaseInfo.WorkflowStatusDesc
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return orderBaseInfo, _err
		}
	}
	return orderBaseInfo, _err
}

func (a AliYunCli) ApproveOrderRequest(approvalType string, taskID int64) (bool, error) {

	ResApprove := false

	flowInstanceId, err := a.getWorkflowInstanceId(taskID)
	if err != nil {
		return false, err
	}

	client, _err := a.CreateClient()
	if _err != nil {
		return false, _err
	}
	approveOrderRequest := &dms_enterprise20181101.ApproveOrderRequest{
		WorkflowInstanceId: flowInstanceId,
		ApprovalType:       tea.String(approvalType),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		Res, _err := client.ApproveOrderWithOptions(approveOrderRequest, runtime)
		if _err != nil {
			return _err
		}
		common.Log.Info(fmt.Sprintf("order: %v approve %v", taskID, *Res.Body.Success))
		ResApprove = *Res.Body.Success
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return ResApprove, _err
		}
	}
	return ResApprove, nil
}

func (a AliYunCli) GetUserInfoByUid(dmsUid string) (string, error) {
	var UserEmail string
	client, _err := a.CreateClient()
	if _err != nil {
		return UserEmail, _err
	}
	getUserRequest := &dms_enterprise20181101.GetUserRequest{
		Uid: tea.String(dmsUid),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		userResp, _err := client.GetUserWithOptions(getUserRequest, runtime)
		if _err != nil {
			return _err
		}
		UserEmail = *userResp.Body.User.Email

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return UserEmail, _err
		}
	}
	return UserEmail, _err
}

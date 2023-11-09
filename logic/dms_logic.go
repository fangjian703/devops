package logic

import (
	"devops/common"
	"devops/config"
	"devops/model/request"
	"devops/utils/tools"
	"encoding/json"
	"fmt"
	"strconv"
)

var Content map[string]interface{}

type DmsLogic struct {
}

func (d DmsLogic) DealDmsBody(dmsBody *request.DmsReqBody) map[string]interface{} {
	err := tools.SendFeiShuCard(dmsBody)
	if err != nil {
		common.Log.Error(err)
	}
	return Content
}

func (d DmsLogic) DealOwnerDmsBody(dmsBody *request.DmsReqBody) error {
	var UsersEmail []string
	var AliCli tools.AliYunCli
	var TaskValue common.OrderInfo

	taskIdKey := strconv.FormatInt(dmsBody.TaskId, 10)
	ExistsNum, err := tools.RdbClient.Exists(tools.Ctx, taskIdKey).Result()

	if err != nil {
		return err
	}

	if ExistsNum > int64(0) {
		return nil
	}

	for _, user := range dmsBody.Receivers {
		userEmail, err := AliCli.GetUserInfoByUid(user.Uid)
		if err != nil {
			common.Log.Error(err)
		}
		UsersEmail = append(UsersEmail, userEmail)
	}

	TaskValue.TaskId = taskIdKey
	TaskValue.UserEmail = UsersEmail

	dataValue, err := json.Marshal(TaskValue)
	if err != nil {
		return err
	}
	err = tools.RdbClient.SetNX(tools.Ctx, taskIdKey, dataValue, common.ExpiredTimeRdb).Err()
	if err != nil {
		return err
	}

	UserIds, err := tools.GetFeiShuOpenIdByEmail(UsersEmail)
	if err != nil {
		return err
	}
	for _, userId := range UserIds {
		err := tools.SendFeiShuCardByOpen(dmsBody, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d DmsLogic) DmsApprovalOrder(fsBody *request.FsCardBody) (err error) {
	var AliCli tools.AliYunCli
	var approvalType, titleColor, msgContent string
	if fsBody.Name == "approve" {
		approvalType = "AGREE"
		titleColor = "green"
		msgContent = "审批通过"
		fsBody.Event = "审批通过"
	}
	if fsBody.Name == "reject" {
		approvalType = "REJECT"
		titleColor = "red"
		msgContent = "审批已拒绝"
		fsBody.Event = "审批拒绝"
	}
	TaskId, _ := strconv.ParseInt(fsBody.TaskId, 10, 64)

	if fsBody.OpenChatId == config.Conf.FeiShu.ChatId {
		taskLockKey := fsBody.TaskId + "group_lock"
		fsBody.OpenId, err = tools.GetChatMember()
		if err != nil {
			return err
		}
		taskExsit, err := tools.RdbClient.Exists(tools.Ctx, taskLockKey).Result()
		if err != nil {
			return nil
		}
		if taskExsit > int64(0) {
			userOpenID, err := tools.RdbClient.Get(tools.Ctx, taskLockKey).Result()
			if err != nil {
				return err
			}
			userName, err := tools.GetFeiShuUserInfoByOpenId(userOpenID)
			if err != nil {
				return err
			}
			msgContent = fmt.Sprintf("审批已由%s操作", userName)
			fsBody.Event = msgContent

			CardContent, err := tools.BuildUpdateCard(fsBody, titleColor, msgContent)
			if err != nil {
				return err
			}
			err = tools.UpdateFeiShuCard(CardContent)
			if err != nil {
				return err
			}
			return nil
		}
		err = tools.RdbClient.SetNX(tools.Ctx, taskLockKey, fsBody.OpenId, common.ExpiredTimeRdb).Err()
		_, err = AliCli.ApproveOrderRequest(approvalType, TaskId)
		if err != nil {
			return err
		}

		CardContent, err := tools.BuildUpdateCard(fsBody, titleColor, msgContent)
		if err != nil {
			return err
		}

		err = tools.UpdateFeiShuCard(CardContent)
		if err != nil {
			return err
		}
		return nil
	}

	if fsBody.OpenChatId != config.Conf.FeiShu.ChatId {
		taskLockKey := fsBody.TaskId + "_lock"
		taskExsit, err := tools.RdbClient.Exists(tools.Ctx, taskLockKey).Result()
		if err != nil {
			return nil
		}
		if taskExsit > int64(0) {
			userOpenID, err := tools.RdbClient.Get(tools.Ctx, taskLockKey).Result()
			if err != nil {
				return err
			}
			userName, err := tools.GetFeiShuUserInfoByOpenId(userOpenID)
			if err != nil {
				return err
			}

			msgContent = fmt.Sprintf("审批已由%s操作", userName)
			fsBody.Event = msgContent

			CardContent, err := tools.BuildUpdateCard(fsBody, titleColor, msgContent)
			if err != nil {
				return err
			}
			err = tools.UpdateFeiShuCard(CardContent)
			if err != nil {
				return err
			}
			return nil
		}

		if taskExsit == int64(0) {
			err = tools.RdbClient.SetNX(tools.Ctx, taskLockKey, fsBody.OpenId, common.ExpiredTimeRdb).Err()
			if err != nil {
				return err
			}
			Success, err := AliCli.ApproveOrderRequest(approvalType, TaskId)
			if err != nil {
				return err
			}
			if !Success {
				orderBaseInfo, err := AliCli.GetWorkflowInfo(TaskId)
				if err != nil {
					return err
				}
				if *orderBaseInfo.StatusCode == "success" {
					titleColor = "green"
					msgContent = *orderBaseInfo.StatusCode
					fsBody.Event = *orderBaseInfo.WorkflowStatusDesc
				}
				if *orderBaseInfo.StatusCode != "success" {
					titleColor = "red"
					msgContent = *orderBaseInfo.StatusCode
					fsBody.Event = *orderBaseInfo.WorkflowStatusDesc
				}
			}
			CardContent, err := tools.BuildUpdateCard(fsBody, titleColor, msgContent)
			if err != nil {
				return err
			}
			err = tools.UpdateFeiShuCard(CardContent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

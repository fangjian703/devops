package tools

import (
	"bytes"
	"context"
	"devops/common"
	"devops/config"
	"devops/model/request"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var FsClient *lark.Client

func InitFeiShu() {
	FsClient = lark.NewClient(config.Conf.FeiShu.AppId, config.Conf.FeiShu.AppSecret)
}

func SendFeiShuCard(dmsReqBody *request.DmsReqBody) error {
	content, err := buildCard(dmsReqBody)
	if err != nil {
		return err
	}

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("chat_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(config.Conf.FeiShu.ChatId).
			MsgType("interactive").
			Content(content).
			Build()).
		Build()

	resp, err := FsClient.Im.Message.Create(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return resp.CodeError
	}

	return nil
}

func GetFeiShuOpenIdByEmail(usersEmail []string) ([]string, error) {
	var UserIds []string
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType("open_id").
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			Emails(usersEmail).
			Mobiles([]string{}).
			Build()).
		Build()
	resp, err := FsClient.Contact.User.BatchGetId(context.Background(), req)
	if err != nil {
		return UserIds, err
	}

	for _, user := range resp.Data.UserList {
		UserIds = append(UserIds, *user.UserId)
	}
	return UserIds, nil
}

func SendFeiShuCardByOpen(dmsReqBody *request.DmsReqBody, openId string) error {
	content, err := buildCard(dmsReqBody)
	if err != nil {
		return err
	}

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("open_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(openId).
			MsgType("interactive").
			Content(content).
			Build()).
		Build()

	resp, err := FsClient.Im.Message.Create(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return resp.CodeError
	}

	return nil
}

func buildCard(dmsReqBody *request.DmsReqBody) (string, error) {
	var bs []byte
	var err error
	if dmsReqBody.Event == "待审批" {
		bs, err = os.ReadFile("template/card_json/wfa.json")
	}
	if dmsReqBody.Event == "审批通过" {
		bs, err = os.ReadFile("template/card_json/approval.json")
	}
	if dmsReqBody.Event == "审批拒绝" {
		bs, err = os.ReadFile("template/card_json/refuse.json")
	}

	if err != nil {
		return "", err
	}

	card := string(bs)
	card = strings.Replace(card, "{category}", dmsReqBody.Category, -1)
	card = strings.Replace(card, "{event}", dmsReqBody.Event, -1)
	card = strings.Replace(card, "{message}", dmsReqBody.Message, -1)
	card = strings.Replace(card, "{module}", dmsReqBody.Module, -1)
	card = strings.Replace(card, "{submitterName}", dmsReqBody.SubmitterName, -1)
	card = strings.Replace(card, "{submitterUid}", dmsReqBody.SubmitterUid, -1)
	card = strings.Replace(card, "{taskId}", fmt.Sprintf("%v", dmsReqBody.TaskId), -1)
	card = strings.Replace(card, "{webUrl}", dmsReqBody.WebUrl, -1)
	return card, nil
}

func UpdateFeiShuCard(cardText string) error {
	time.Sleep(1 * time.Second)
	var tenantAccessToken string
	url := "https://open.feishu.cn/open-apis/interactive/v1/card/update"
	TAK := "tenantAccessToken"
	KExists, err := RdbClient.Exists(Ctx, TAK).Result()
	if err != nil {
		return nil
	}
	if KExists == int64(0) {
		tenantAccessTokenResp, err := FsClient.GetTenantAccessTokenBySelfBuiltApp(context.Background(), &larkcore.SelfBuiltTenantAccessTokenReq{
			AppID:     config.Conf.FeiShu.AppId,
			AppSecret: config.Conf.FeiShu.AppSecret,
		})
		if err != nil {
			return err
		}
		tenantAccessToken = tenantAccessTokenResp.TenantAccessToken
		err = RdbClient.SetNX(Ctx, TAK, tenantAccessToken, 3600*time.Second).Err()
		if err != nil {
			return err
		}
	}
	if KExists > int64(0) {
		tenantAccessToken, _ = RdbClient.Get(Ctx, TAK).Result()
	}

	var jsonData = []byte(cardText)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	HTenantAccessToken := fmt.Sprintf("Bearer %v", tenantAccessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", HTenantAccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	common.Log.Info(string(respBody))
	return nil
}

func BuildUpdateCard(fsBody *request.FsCardBody, titleColor string, msgContent string) (string, error) {
	bs, err := os.ReadFile("template/card_json/update_wfa.json")
	if err != nil {
		return "", err
	}

	card := string(bs)
	card = strings.Replace(card, "{token}", fsBody.Token, -1)
	card = strings.Replace(card, "{openIds}", fsBody.OpenId, -1)
	card = strings.Replace(card, "{category}", fsBody.Category, -1)
	card = strings.Replace(card, "{event}", fsBody.Event, -1)
	card = strings.Replace(card, "{submitterName}", fsBody.SubmitterName, -1)
	card = strings.Replace(card, "{taskId}", fmt.Sprintf("%v", fsBody.TaskId), -1)
	card = strings.Replace(card, "{webUrl}", fsBody.WebUrl, -1)
	card = strings.Replace(card, "{titleColor}", titleColor, -1)
	card = strings.Replace(card, "{msgContent}", msgContent, -1)
	return card, nil
}

func GetChatMember() (string, error) {
	var openIds string
	req := larkim.NewGetChatMembersReqBuilder().
		ChatId(config.Conf.FeiShu.ChatId).
		MemberIdType("open_id").
		PageSize(10).
		PageToken("").
		Build()
	MemResp, err := FsClient.Im.ChatMembers.Get(context.Background(), req)
	if err != nil {
		return openIds, err
	}
	ItemsNum := len(MemResp.Data.Items)
	for index, item := range MemResp.Data.Items {
		memberId := *item.MemberId
		if index == 0 {
			openIds = fmt.Sprintf("%s\",", memberId)
			continue
		}
		if index == ItemsNum-1 {
			openIds = openIds + "\"" + memberId
			continue
		}
		openIds = openIds + "\"" + memberId + "\","
	}
	return openIds, nil
}

func GetFeiShuUserInfoByOpenId(openId string) (string, error) {
	req := larkcontact.NewGetUserReqBuilder().
		UserId(openId).
		UserIdType("open_id").Build()
	UserResp, err := FsClient.Contact.User.Get(context.Background(), req)
	if err != nil {
		return "", err
	}
	return *UserResp.Data.User.Name, nil
}

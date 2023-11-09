package request

type DmsReqBody struct {
	MessageEvent    `json:"messageEvent"`
	SignatureMethod string `json:"signatureMethod"`
}

type MessageEvent struct {
	Category      string  `json:"category"`
	Event         string  `json:"event"`
	EventTime     int64   `json:"eventTime"`
	Message       string  `json:"message"`
	Module        string  `json:"module"`
	Receivers     []Users `json:"receivers"`
	SubmitterName string  `json:"submitterName"`
	SubmitterUid  string  `json:"submitterUid"`
	TargetUsers   []Users `json:"targetUsers"`
	TaskId        int64   `json:"taskId"`
	WebUrl        string  `json:"webUrl"`
}

type Users struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

type FsCardBody struct {
	Challenge     string `json:"challenge"`
	AppId         string `json:"app_id"`
	OpenId        string `json:"open_id"`
	UserId        string `json:"user_id"`
	OpenMessageId string `json:"open_message_id"`
	OpenChatId    string `json:"open_chat_id"`
	TenantKey     string `json:"tenant_key"`
	Token         string `json:"token"`
	ResAction     `json:"action"`
}

type ResAction struct {
	ResValue `json:"value"`
	Tag      string `json:"tag"`
}

type ResValue struct {
	Category      string `json:"category"`
	Event         string `json:"event"`
	Name          string `json:"name"`
	Uid           string `json:"uid"`
	SubmitterName string `json:"submitterName"`
	SubmitterUid  string `json:"submitterUid"`
	TaskId        string `json:"taskId"`
	WebUrl        string `json:"webUrl"`
}

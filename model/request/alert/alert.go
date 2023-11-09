package alert

// ReqAlert 定义接收JSON数据的结构体
type ReqAlert struct {
	Status       string              `json:"status"`
	StartsAt     string              `json:"startsAt"`
	EndsAt       string              `json:"endsAt"`
	GeneratorURL string              `json:"generatorURL"`
	Fingerprint  string              `json:"fingerprint"`
	Labels       ReqAlertLabel       `json:"labels"`
	Annotations  ReqAlertAnnotations `json:"annotations"`
}

type ReqGroupLabels struct {
	AlertName string `json:"alertname"`
}

type ReqCommonLabels struct {
	AlertName string `json:"alertname"`
	Instance  string `json:"instance"`
	Job       string `json:"job"`
	Severity  string `json:"severity"`
}

type ReqCommonAnnotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type ReqAlertLabel struct {
	AlertName string `json:"alertname"`
	Instance  string `json:"instance"`
	Job       string `json:"job"`
	Severity  string `json:"severity"`
}

type ReqAlertAnnotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type ReqBody struct {
	// alertmanager传来的请求体
	Receiver          string               `json:"receiver"`
	Status            string               `json:"status"`
	ExternalURL       string               `json:"externalURL"`
	Version           string               `json:"version"`
	GroupKey          string               `json:"groupkey"`
	TruncatedAlerts   float64              `json:"truncatedAlerts"`
	Alert             []ReqAlert           `json:"alerts"`
	GroupLabels       ReqGroupLabels       `json:"groupLabels"`
	CommonLabels      ReqCommonLabels      `json:"commonLabels"`
	CommonAnnotations ReqCommonAnnotations `json:"commonAnnotations"`
}

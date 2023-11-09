package common

import (
	"strings"
)

func SpliceString(status, description, summary string) string {
	// 拼接字符串
	var rst strings.Builder
	rst.WriteString("# 告警接收测试 \n\n")
	rst.WriteString("- status: " + status + " \n\n")
	rst.WriteString("- description: " + description + " \n\n")
	rst.WriteString("- summary: " + summary + " \n\n")
	return rst.String()
}

package middleware

import "github.com/gin-gonic/gin"

// CheckSign 定义一个中间件用于验证请求头签名
func CheckSign(c *gin.Context) {
	rSign := c.Request.Header.Get("sign")
	cSign := "640834f8-d6de-4513-8c38-1ec27a211e0c"
	if rSign == cSign {
		c.Next()
	} else {
		data := make(map[string]interface{})
		data["msg"] = "签名错误，请检查"
		c.JSON(403, data)
		c.Abort()
	}
}

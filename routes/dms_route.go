package routes

import (
	"devops/controller"

	"github.com/gin-gonic/gin"
)

// InitJenkinsRoutes 注册路由
func InitDmsRoutes(r *gin.RouterGroup) gin.IRoutes {

	jenkins := r.Group("/dms")
	{
		jenkins.POST("/owner", controller.Dms.OwnerDmsReq)
		jenkins.POST("/devops", controller.Dms.GetDmsReqBody)
		jenkins.POST("/fscard", controller.Dms.GetFeiShuReqBody)
	}
	return r
}

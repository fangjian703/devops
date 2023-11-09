package routes

import (
	"devops/controller"

	"github.com/gin-gonic/gin"
)

// InitBaseRoutes 注册基础路由
func InitBaseRoutes(r *gin.RouterGroup) gin.IRoutes {
	base := r.Group("/base")
	{
		base.GET("ping", controller.Demo)
	}
	return r
}

package routes

import (
	"github.com/gin-gonic/gin"
	"picasso/module/qa"
	"picasso/module/user"
)

func Build(r *gin.Engine) {
	group := r.Group("/api/v1") // *gin.RouterGroup
	user.NewUserController().Build(group)
	qa.NewQaController().Build(group)
}

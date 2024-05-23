package middlewares

import (
	"github.com/gin-gonic/gin"
	. "picasso/pkg/utils"
)

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//context.JSON(500, gin.H{"error": err})
				ret := ResultWrapper(context)(nil, err.(string))(Error)
				context.JSON(500, ret)
			}
		}()
		context.Next()
	}
}

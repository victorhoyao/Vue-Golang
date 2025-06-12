package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, httpStatus, code int, data gin.H, msg interface{}) {
	c.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": msg})
}

// 请求成功
func Success(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, 200, data, msg)
}

// 参数错误
func Fail(c *gin.Context, data gin.H, msg interface{}) {
	Response(c, http.StatusOK, 422, data, msg)
}

// 系统错误
func ServerBad(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, 500, data, msg)
}

// 请求失败
func RequsetBad(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, 400, data, msg)
}

// 权限不足
func Unauthorized(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, 403, data, msg)
}

// 权限不足
func AuthError(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, 401, data, msg)
}

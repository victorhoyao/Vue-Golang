package middleware

import (
	"BTaskServer/common"
	"BTaskServer/model"
	"BTaskServer/util/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware
//
//	@Description: 校验用户是否登录
//	@return gin.HandlerFunc
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取authorization header
		tokenString := c.GetHeader("token")

		if tokenString == "" {
			response.AuthError(c, nil, "用户未登录")
			c.Abort()
			return
		}
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			response.AuthError(c, nil, "用户未登录")
			c.Abort()
			return
		}

		// 获得userID,通过userid查询用户是否存在
		var UserModel model.User
		db := common.GetDB()
		userid := claims.UserId
		if res := db.Where("id = ?", userid).First(&UserModel); res.RowsAffected == 0 {
			response.AuthError(c, nil, "用户未登录")
			c.Abort()
			return
		}

		// 用户如果存在，将用户id存入上下文
		c.Set("user", UserModel)
		c.Next()
	}
}

/**
 * AuthWebsocketMiddleware
 * @Description: websocket权限校验中间件
 * @return gin.HandlerFunc
 */
func AuthWebsocketMiddleware(tokenString string) (bool, uint) {
	if tokenString == "" {
		return false, 0
	}
	token, claims, err := common.ParseToken(tokenString)
	if err != nil || !token.Valid {
		return false, 0
	}

	// 获得userID,通过userid查询用户是否存在
	var UserModel model.User
	db := common.GetDB()
	userid := claims.UserId
	if res := db.Where("id = ?", userid).First(&UserModel); res.RowsAffected == 0 {
		return false, 0
	}
	return true, userid
}

package middleware

// 定义最大并发请求数
//const maxConcurrentRequests = 100
//
//// 计数器
//var currentRequests int32

//// @return gin.HandlerFunc
//func LimitMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 检查当前请求数是否超出最大并发数
//		if atomic.LoadInt32(&currentRequests) >= maxConcurrentRequests {
//			// 返回系统繁忙错误
//			c.JSON(http.StatusTooManyRequests, gin.H{
//				"error": "系统繁忙，请稍后再试",
//			})
//			c.Abort() // 中止请求，不再继续处理
//			return
//		}
//
//		// 增加计数
//		atomic.AddInt32(&currentRequests, 1)
//
//		// 使用 defer 确保请求结束时减少计数
//		defer atomic.AddInt32(&currentRequests, -1)
//
//		// 继续处理请求
//		c.Next()
//	}
//}

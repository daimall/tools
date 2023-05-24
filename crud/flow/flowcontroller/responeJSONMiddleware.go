package flowcontroller

import (
	"net/http"

	"github.com/daimall/tools/crud/common"
	"github.com/daimall/tools/curd/customerror"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 自定义 JSON 封装中间件
func ResponseJSONMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		viper.SetDefault("restful.code", "code")
		viper.SetDefault("restful.message", "message")
		viper.SetDefault("restful.data", "data")
		// 检测 gin errors
		if len(c.Errors) > 0 {
			// 如果有错误，封装错误信息并返回
			c.JSON(http.StatusOK, gin.H{
				viper.GetString("restful.code"):    c.Errors.Last().Type,
				viper.GetString("restful.message"): c.Errors.String(),
			})
			c.Abort()
			return
		}
		// 检查自定义错误
		if err, ok := c.Keys[common.CustomErrKey].(customerror.CustomError); ok {
			c.JSON(http.StatusOK, gin.H{
				viper.GetString("restful.code"):    err.GetCode(),
				viper.GetString("restful.message"): err.GetMessage(),
			})
			c.Abort()
			return
		}
		c.Next() // 执行路由处理函数
		// 检测 gin errors
		if len(c.Errors) > 0 {
			// 如果有错误，封装错误信息并返回
			c.JSON(http.StatusOK, gin.H{
				viper.GetString("restful.code"):    c.Errors.Last().Type,
				viper.GetString("restful.message"): c.Errors.String(),
			})
			c.Abort()
			return
		}
		// 检查自定义错误
		if err, ok := c.Keys[common.CustomErrKey].(customerror.CustomError); ok {
			c.JSON(http.StatusOK, gin.H{
				viper.GetString("restful.code"):    err.GetCode(),
				viper.GetString("restful.message"): err.GetMessage(),
			})
			c.Abort()
			return
		}
		// 如果没有错误，封装成功信息并返回
		c.JSON(http.StatusOK, gin.H{
			viper.GetString("restful.code"):    0,
			viper.GetString("restful.message"): "OK",
			viper.GetString("restful.data"):    c.Keys[common.ResponeDataKey], // 这里的 是请求处理程序设置的响应数据
		})
	}
}

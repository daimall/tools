package flowcontroller

import (
	"github.com/daimall/tools/crud/common"
	"github.com/daimall/tools/crud/customerror"
	"github.com/daimall/tools/crud/flow/flowservice"
	"github.com/daimall/tools/crud/logger"
	"github.com/daimall/tools/functions"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 自定义 JSON 封装中间件
func PrepareMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		var username string
		var err error
		var ok bool
		// 获取 session中的username
		if username, ok = session.Get(common.UserNameSessionKey).(string); !ok {
			// session中没有username，从header中获取token，然后解析username
			if username, err = functions.GetAccountIdFromToken(c.GetHeader(common.TokenKey), []byte("test_group")); err != nil {
				// session中没有，token没有解析出，判断当前环境是不是调试环境
				if viper.GetString("env.runmode") == "dev" {
					viper.SetDefault("env.username", "devuser")
					username = viper.GetString("env.username")
				}
			}
		}
		// 没有找到用户名，设置错误，进入结果处理中间件
		if username == "" {
			c.Set(common.CustomErrKey, customerror.UsernameNotFound)
			c.Next()
			return
		}
		crudContext := flowservice.CRUDContext{}
		idStr := c.Param("id")
		if idStr != "" {
			if crudContext.ServiceId, err = functions.Str2Uint(idStr); err != nil {
				logger.Error("get serviceId[%s] failed, %s", idStr, err.Error())
				c.Set(common.CustomErrKey, customerror.ServiceIdNotInt)
				c.Next()
				return
			}
		}
		crudContext.UserName = username
		if crudContext.Service, err = flowservice.GetService(c.Param("service")); err != nil {
			logger.Error("get service bye name[%s] failed, %s", c.Param("service"), err.Error())
			c.Set(common.CustomErrKey, err)
			c.Next()
			return
		}
		c.Set(common.CRUDContextKey, crudContext)
		c.Next() // 执行下一个插件
	}
}

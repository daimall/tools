package flowcontroller

import (
	"net/http"

	"github.com/daimall/tools/crud/common"
	"github.com/daimall/tools/crud/dbmysql/dbgorm"
	"github.com/daimall/tools/crud/flow/flowservice"
	"github.com/daimall/tools/crud/oplog"
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
		checkGinError(c)
		c.Next() // 执行路由处理函数
		// 检测 gin errors
		checkGinError(c)
		// 如果没有错误，封装成功信息并返回
		c.JSON(http.StatusOK, gin.H{
			viper.GetString("restful.code"):    0,
			viper.GetString("restful.message"): "OK",
			viper.GetString("restful.data"):    c.Keys[common.ResponeDataKey], // 这里的 是请求处理程序设置的响应数据
		})
		// 记录操作日志
		recordOperateLog(c)
	}
}

// 错误检查
func checkGinError(c *gin.Context) {
	// 检测 gin errors
	if len(c.Errors) > 0 {
		// 如果有错误，封装错误信息并返回
		c.JSON(http.StatusBadRequest, gin.H{
			viper.GetString("restful.code"):    c.Errors.Last().Type,
			viper.GetString("restful.message"): c.Errors.String(),
		})
		c.Abort()
		return
	}
	// 检查自定义错误
	if err, ok := c.Keys[common.CustomErrKey].(customerror.CustomError); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			viper.GetString("restful.code"):    err.GetCode(),
			viper.GetString("restful.message"): err.GetMessage(),
		})
		c.Abort()
		return
	}
}

// 记录操作日志
func recordOperateLog(c *gin.Context) {
	if crudContextInf, ok := c.Get(common.CRUDContextKey); ok {
		if crudContext, ok := crudContextInf.(flowservice.CRUDContext); ok {
			if crudContext.OperateLog != "" {
				var logModel interface{}
				// servcie 实现了自定义日志输出函数
				if logService, ok := crudContext.Service.(flowservice.OplogModelInf); ok {
					logModel = logService.OplogModel(c, crudContext)
				}
				// 没有自定义日志输出接口，采用公共的
				if logModel == nil {
					logModel = &oplog.OpLog{User: crudContext.UserName, Action: crudContext.Action,
						FlowId: crudContext.ServiceId, Flow: crudContext.ServiceName, Remark: crudContext.OperateLog}
				}
				// 记录操作日志
				oplog.AddOperationLog(dbgorm.GetDBInst(), logModel)

			}
		}
	}

	// if log == "" {
	// 	// 不记录操作日志
	// 	return
	// }
	// var logModel interface{}
	// if logService, ok := c.Service.(flowservice.OplogModelInf); ok {
	// 	logModel = logService.OplogModel(c.uname, c.ServiceName, serviceId, action, log)
	// }
	// if logModel == nil {
	// 	logModel = &oplog.OpLog{User: c.uname, Action: action,
	// 		FlowId: serviceId, Flow: c.ServiceName, Remark: log}
	// }
	// oplog.AddOperationLog(dbgorm.GetDBInst(), logModel)
}

package crudcontroller

import (
	"runtime"

	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/curd/common"
	"github.com/daimall/tools/curd/customerror"
	dbgorm "github.com/daimall/tools/curd/dbmysql/dbgorm/mysql"
	"github.com/daimall/tools/curd/flow/v1/flowcomm"
	curdservice "github.com/daimall/tools/curd/flow/v1/services"
	oplog "github.com/daimall/tools/curd/oplog"
)

// 继承公共基础
type BaseController struct {
	common.BaseController
	Service     curdservice.CrudService //
	ServiceName string                  //
	uname       string                  // 用户名
}

// 记录操作日志
func (c *BaseController) LogFunc(serviceId int64, action, log string) {
	if log == "" {
		// 不记录操作日志
		return
	}
	var logModel interface{}
	if logService, ok := c.Service.(curdservice.OplogModelInf); ok {
		logModel = logService.OplogModel(c.uname, c.ServiceName, serviceId, action, log)
	}
	if logModel == nil {
		logModel = &oplog.OpLog{User: c.uname, Action: action,
			FlowId: serviceId, Flow: c.ServiceName, Remark: log}
	}
	oplog.AddOperationLog(dbgorm.GetDBInst(), logModel)
}

// 预执行，获取service对象
func (c *BaseController) NestPrepare() {
	var ok bool
	if c.uname, ok = c.GetSession(common.UserNameSessionKey).(string); ok {
		c.ServiceName = c.Ctx.Input.Param(":service")
		c.Service = curdservice.GetService(c.ServiceName)
		return
	}
	logs.Error("username[KEY:%s] does not exist in session.", common.UserNameSessionKey)
	c.JSONResponse(flowcomm.UnameNotFound)
	c.StopRun()
}

//ResponseJSON（重写方法） 返回JSON格式结果
func (c *BaseController) ResponseJSON(err error, ret interface{}, serviceId int64, action, oplog string) {
	var method string
	pc, _, _, _ := runtime.Caller(1)
	method = runtime.FuncForPC(pc).Name()
	if err == nil {
		c.Data["json"] = common.StandRestResult{Code: 0, Data: ret, Message: "OK"}
		// 记录操作日志
		c.LogFunc(serviceId, action, oplog)
	} else if customErr, ok := err.(customerror.CustomError); ok {
		c.Data["json"] = common.StandRestResult{Code: customErr.GetCode(), Data: ret, Message: customErr.GetMessage()}
		logs.Error("FlowController[%s]%s(customErr)", method, err.Error())
	} else {
		c.Data["json"] = common.StandRestResult{Code: -1, Data: ret, Message: err.Error()}
		logs.Error("FlowController[%s]%s", method, err.Error())
	}
	c.ServeJSON()
}

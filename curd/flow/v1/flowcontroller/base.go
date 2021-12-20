package flowcontroller

import (
	"runtime"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/curd/common"
	"github.com/daimall/tools/curd/customerror"
	dbgorm "github.com/daimall/tools/curd/dbmysql/dbgorm/mysql"
	"github.com/daimall/tools/curd/flow/v1/flowservice"
	oplog "github.com/daimall/tools/curd/oplog"
)

// 继承公共基础
type BaseController struct {
	common.BaseController
	Service     flowservice.FlowService //
	ServiceName string                  //
	uname       string                  // 用户名
}

// 记录操作日志
func (c *BaseController) LogFunc(serviceId uint, action, log string) {
	if log == "" {
		// 不记录操作日志
		return
	}
	var logModel interface{}
	if logService, ok := c.Service.(flowservice.OplogModelInf); ok {
		logModel = logService.OplogModel(c.uname, c.ServiceName, serviceId, action, log)
	}
	if logModel == nil {
		logModel = &oplog.OpLog{User: c.uname, Action: action,
			FlowId: serviceId, Flow: c.ServiceName, Remark: log}
	}
	oplog.AddOperationLog(dbgorm.GetDBInst(), logModel)
}

// 预执行，获取service对象
func (c *BaseController) Prepare() {
	var ok bool
	if c.uname, ok = c.GetSession(common.UserNameSessionKey).(string); ok {
		c.ServiceName = c.Ctx.Input.Param(":service")
		c.Service = flowservice.GetService(c.ServiceName)
		if app, ok := c.Service.(flowservice.SetBaseControllerInf); ok {
			app.SetBaseController(c.BaseController)
		}
		return
	}
	logs.Error("username[KEY:%s] does not exist in session.", common.UserNameSessionKey)
	c.JSONResponse(common.UnameNotFound)
	c.StopRun()
}

//ResponseJSON（重写方法） 返回JSON格式结果
func (c *BaseController) ResponseJSON(err error, ret interface{}, serviceId uint, action, oplog string) {
	var method string
	pc, _, _, _ := runtime.Caller(1)
	method = runtime.FuncForPC(pc).Name()
	if err == nil {
		// 记录操作日志
		c.LogFunc(serviceId, action, oplog)
		c.JSONResponse(nil, ret)
	} else if customErr, ok := err.(customerror.CustomError); ok {
		logs.Error("FlowController[%s]%s(customErr)", method, err.Error())
		c.JSONResponse(customErr, nil)
	} else {
		logs.Error("FlowController[%s]%s", method, err.Error())
		c.JSONResponse(customerror.New(-1, err.Error()))
	}
}

// 公共方法
func getUintID(idStr string) (id uint, err error) {
	var v int
	if v, err = strconv.Atoi(idStr); err != nil {
		return
	}
	id = uint(v)
	return
}

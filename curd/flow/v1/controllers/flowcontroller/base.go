package flowcontroller

import (
	"net/http"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/curd/common"
	"github.com/daimall/tools/curd/customtypes/customerror"
	"github.com/daimall/tools/flow/v3/dbutil"
	"github.com/daimall/tools/flow/v3/flowcomm"
	"github.com/daimall/tools/flow/v3/services/flowservice"
	oplog "github.com/daimall/tools/oplog/v2"
	"github.com/daimall/tools/sign"
)

// 继承公共基础
type BaseController struct {
	common.BaseController
	Service     flowservice.FlowService //
	ServiceName string                  //
	uname       string                  // 用户名
	appid       string                  // 请求的appid
	timestamp   string                  // 时间戳
	sign        string                  // 签名字符串
}

// 记录操作日志
func (c *BaseController) LogFunc(serviceId int64, action, log string) {
	if log == "" {
		// 不记录操作日志
		return
	}
	logModel := c.Service.OplogModel(c.uname, c.ServiceName, serviceId, action, log)
	if logModel == nil {
		logModel = &oplog.OpLog{User: c.uname, Action: action,
			FlowId: serviceId, Flow: c.ServiceName, Remark: log}
	}
	oplog.AddOperationLog(dbutil.GetDBInst(), logModel)
}

// 预执行，获取service对象
func (c *BaseController) NestPrepare() {
	var ok bool
	if c.uname, ok = c.GetSession(common.UserNameSessionKey).(string); ok {
		c.ServiceName = c.Ctx.Input.Param(":service")
		c.Service = flowservice.GetService(c.ServiceName)
		return
	}
	logs.Error("username[KEY:%s] does not exist in session.", common.UserNameSessionKey)
	c.ResponseError(flowcomm.UnameNotFound)
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

func (c *BaseController) CheckSign() *common.StandRestResult {
	if ret := c.checkHeader(); ret != nil {
		return ret
	}
	header := c.Ctx.Request.Header
	body := c.Ctx.Input.RequestBody
	if ok := c.checkSign(header, string(body)); !ok {
		logs.Error("Sign error")
		return c.NewResponse(common.SignError)
	}
	return nil
}

func (c *BaseController) checkHeader() *common.StandRestResult {
	h := c.Ctx.Request.Header
	if c.appid = h.Get("appid"); c.appid == "" {
		logs.Error("appid in header is nil")
		return c.NewResponse(common.AppIDError)
	}
	if c.timestamp = h.Get("timestamp"); c.timestamp == "" {
		logs.Error("timestamp in header is nil")
		return c.NewResponse(common.TimestampError)
	}
	if c.sign = h.Get("sign"); c.sign == "" {
		logs.Error("sign in header is nil")
		return c.NewResponse(common.SignError)
	}
	return nil
}

func (c *BaseController) checkSign(h http.Header, body string) bool {
	signMap := make(map[string]string)
	signMap["timestamp"] = c.timestamp
	signMap["body"] = body
	appkey := beego.AppConfig.String("Apps::" + c.appid)
	if appkey == "" {
		logs.Error("appid error(can not get app key),", c.appid)
		return false
	}
	if !sign.New().SetType(sign.FROM_AUTHOR).
		SetAppId(c.appid).
		SetKey(appkey).
		SetIsSgin(beego.AppConfig.DefaultBool("SECURITY::NEED_SIGN", false)).
		ToLower(true).VerifyMapSign(c.sign, signMap) {
		logs.Error("Sign error, %+v", signMap)
		return false
	}
	return true
}

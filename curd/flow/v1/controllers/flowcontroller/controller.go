package flowcontroller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/curd/common"
	"github.com/daimall/tools/flow/v3/flowcomm"
	"github.com/daimall/tools/flow/v3/services/flowservice"
)

// FlowController operations for model
type FlowController struct {
	BaseController
}

// URLMapping ...
func (c *FlowController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Next", c.Next)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Configs", c.Configs)
	c.Mapping("Action", c.Action)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("DeleteList", c.DeleteList)
	c.Mapping("GetHistory", c.GetHistory)
}

// Post ...
// @Title 创建一条流程
// @Description create Service
// @Param	service	string 	Service	true		"body for Service content"
// @Success 201 {int} Service
// @Failure 403 body is empty
// @router /:service [post]
func (c *FlowController) Post() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	var service = flowservice.GetService(c.Ctx.Input.Param(":service"))
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionCreate, oplog)
	}()
	serviceId, ret, oplog, err = service.New(c.uname, c.BaseController.BaseController)
}

// Next ...
// @Title 进入流程下一步
// @Description handle Service to next step
// @Param	body		body 	params	true		"body for Service content"
// @Success 201 {int} OK(step info)
// @Failure 403 body is empty
// @router /:service/:id/:handlerId [post]
func (c *FlowController) Next() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	var handler flowservice.FlowHandler
	var handlerId int64
	var action string
	defer func() {
		c.ResponseJSON(err, ret, serviceId, action, oplog)
	}()

	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	idStr = c.Ctx.Input.Param(":handlerId")
	if handlerId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	if serviceId != 0 {
		if c.Service, err = c.Service.LoadInst(serviceId); err != nil {
			logs.Error("service.LoadInst failed,", err.Error())
			return
		}
	}
	var step flowservice.FlowStep
	if step, err = c.Service.GetCurStep(); err != nil {
		logs.Error("GetCurStep failed, %s", err.Error())
		return
	}
	action = step.Key()
	if dataFlow, ok := c.Service.(flowservice.DataCacheInf); ok {
		dataFlow.SetData("BaseController", c.BaseController.BaseController)
	}
	if handler, err = step.Hander().LoadInst(c.Service, c.uname, handlerId); err != nil {
		logs.Error("LoadInst failed, %s", err.Error())
		return
	}
	ret, oplog, err = handler.Do(c.uname, c.BaseController.BaseController)
}

// Action ...
// @Title 独立动作
// @Description handle a action
// @Param	body		body 	params	true		"body for Service content"
// @Success 201 {int} OK(step info)
// @Failure 403 body is empty
// @router /:service/:id/:action [post]
func (c *FlowController) Action() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	var action flowservice.Action
	var actionType = c.Ctx.Input.Param(":action")
	defer func() {
		//日志中包含serviceId 信息
		oplogSep := strings.Split(oplog, flowservice.LogSep)
		if len(oplogSep) == 2 {
			if serviceId, err = strconv.ParseInt(oplogSep[0], 0, 64); err == nil {
				oplog = oplogSep[1]
			}
		}
		c.ResponseJSON(err, ret, serviceId, actionType, oplog)
	}()

	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	if actionType == "" {
		err = flowcomm.ActionNotFound
		return
	}
	if serviceId != 0 {
		if c.Service, err = c.Service.LoadInst(serviceId); err != nil {
			logs.Error("service.LoadInst failed,", err.Error())
			return
		}
	}
	if actionApp, ok := c.Service.(flowservice.ActionInf); ok {
		if action, err = actionApp.GetAction(serviceId, actionType); err != nil {
			logs.Error("GetAction[%s] failed, %s", actionType, err.Error())
			return
		}
	} else {
		err = fmt.Errorf("GetAction method is not implement")
		return
	}
	ret, oplog, err = action.Do(c.uname, serviceId, actionType, c.BaseController.BaseController)
}

// Configs ...
// @Title Get 获取参数对象结构
// @Description get configs for next
// @Param id path string true "The key for staticblock"
// @Success 200 {object} Service
// @Failure 403 :id is empty
// @router /:service/:id/config [get]
func (c *FlowController) Configs() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionGetConfigs, oplog)
	}()
	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	ret, oplog, err = c.Service.GetConfigs(c.uname, serviceId, c.BaseController.BaseController)
}

// GetOne ...
// @Title Get One
// @Description get Service by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} Service
// @Failure 403 :id is empty
// @router /:service/:id [get]
func (c *FlowController) GetOne() {
	var err error
	var oplog string
	var serviceId int64
	var ret interface{}
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionGetOne, oplog)
	}()

	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		logs.Error(":id to int64 failed", err.Error())
		return
	}
	ret, oplog, err = c.Service.GetOne(c.uname, serviceId)
}

// GetAll ...
// @Title Get All
// @Description get Service
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} Service
// @Failure 403
// @router /:service [get]
func (c *FlowController) GetAll() {
	var err error
	var oplog string
	var serviceId int64
	var ret struct {
		Items interface{} `json:"items"`
		Total int64       `json:"total"`
	}
	var l interface{}
	var count int64
	defer func() {
		ret.Items = l
		ret.Total = count
		c.ResponseJSON(err, ret, serviceId, ServiceActionGetAll, oplog)
	}()

	var fields []string
	var sortby []string
	var order []string
	var query []*common.QueryConditon
	var limit int64 = 10
	var offset int64
	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		v = strings.Replace(v, ".", "__", -1)
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k|type:v|v|v,k|type:v|v|v  其中Type可以没有,默认值是 MultiText
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") { // 分割多个查询key
			qcondtion := new(common.QueryConditon)
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				logs.Error("query condtion format error:%s, need key:value", kv)
				c.ResponseError(flowcomm.QueryCondErr)
				return
			}
			k_init, v_init := kv[0], kv[1]         // 初始分割查询key和value（备注，value是多个用|分割）
			key_type := strings.Split(k_init, "|") // 解析key中的type信息
			if len(key_type) == 2 {
				qcondtion.QueryKey = key_type[0]
				qcondtion.QueryType = key_type[1]
			} else if len(key_type) == 1 {
				qcondtion.QueryKey = key_type[0]
				qcondtion.QueryType = common.MultiText
			} else {
				logs.Error("Error: invalid query key|type format," + k_init)
				c.ResponseError(flowcomm.QueryCondErr)
				return
			}
			qcondtion.QueryValues = strings.Split(v_init, "|") // 解析出values信息
			qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
			query = append(query, qcondtion)
		}
	}
	l, count, oplog, err = c.Service.GetAll(c.uname, query, fields, sortby, order, offset, limit)
}

// Put ...
// @Title Put
// @Description update the Service
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	Service	true		"body for Service content"
// @Success 200 {object} Service
// @Failure 403 :id is not int
// @router /:service/:id [put]
func (c *FlowController) Put() {
	var err error
	var ret interface{}
	var oplog string
	var serviceId int64
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionPut, oplog)
	}()

	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	if serviceId != 0 {
		if c.Service, err = c.Service.LoadInst(serviceId); err != nil {
			logs.Error("service.LoadInst failed,", err.Error())
			return
		}
	}
	if updateApp, ok := c.Service.(flowservice.UpdateInf); ok {
		ret, oplog, err = updateApp.Update(serviceId, c.BaseController.BaseController)
	} else {
		err = errors.New("update interface not impl")
	}
}

// Delete ...
// @Title Delete
// @Description delete the Service
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:service/:id [delete]
func (c *FlowController) Delete() {
	var err error
	var ret interface{}
	var oplog string
	var serviceId int64
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionDelete, oplog)
	}()

	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	ret, oplog, err = c.Service.Delete(serviceId)
}

// DeleteList ...
// @Title multi-Delete
// @Description delete multi Services
// @Param	ids	 	string	true		"The ids you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:service/deletelist [delete]
func (c *FlowController) DeleteList() {
	var err error
	var ret interface{}
	var oplog string
	var serviceId int64
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionDeleteList, oplog)
	}()

	ret, oplog, err = c.Service.MultiDelete(strings.Split(c.GetString("ids"), ","))
}

// GetHistory ...
// @Title Get 获取流程的操作记录（历史）
// @Description get get operation list history
// @Param id path string true "The key for staticblock"
// @Success 200 {object} Service
// @Failure 403 :id is empty
// @router /:service/:id/oplist [get]
func (c *FlowController) GetHistory() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceOpList, oplog)
	}()
	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	if c.Service, err = c.Service.LoadInst(serviceId); err != nil {
		logs.Error("service.LoadInst failed,", err.Error())
		return
	}
	if historyApp, ok := c.Service.(flowservice.OpHistoryInf); ok {
		ret, oplog, err = historyApp.GetOpHistory()
	} else {
		err = errors.New("GetOpHistory interface not impl")
	}
}

// GetOpLogHistory ...
// @Title Get 获取流程等操作日志记录
// @Description get get operation list history
// @Param id path string true "The key for staticblock"
// @Success 200 {object} Service
// @Failure 403 :id is empty
// @router /:service/:id/oploglist [get]
func (c *FlowController) GetOpLogHistory() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceOpList, oplog)
	}()
	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	if c.Service, err = c.Service.LoadInst(serviceId); err != nil {
		logs.Error("service.LoadInst failed,", err.Error())
		return
	}
	if historyApp, ok := c.Service.(flowservice.OpLogHistoryInf); ok {
		ret, oplog, err = historyApp.GetOpLogHistory()
	} else {
		err = errors.New("GetOpLogHistory interface not impl")
	}
}

// GetPreHandlers ...
// @Title Get 获取上一步处理人（退回流程使用）
// @Description get get operation list history
// @Param id path string true "The key for staticblock"
// @Success 200 {object} Service
// @Failure 403 :id is empty
// @router /:service/:id/prehandlers [get]
func (c *FlowController) GetPreHandlers() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceOpList, oplog)
	}()
	idStr := c.Ctx.Input.Param(":id")
	if serviceId, err = strconv.ParseInt(idStr, 0, 64); err != nil {
		return
	}
	if c.Service, err = c.Service.LoadInst(serviceId); err != nil {
		logs.Error("service.LoadInst failed,", err.Error())
		return
	}
	if preHandlersApp, ok := c.Service.(flowservice.PreHandlersInf); ok {
		ret, oplog, err = preHandlersApp.GetPreHandlers()
	} else {
		err = errors.New("GetOpHistory interface not impl")
	}
}

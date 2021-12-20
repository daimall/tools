package crudcontroller

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/curd/common"
	"github.com/daimall/tools/curd/customerror"
	"github.com/daimall/tools/curd/flow/v1/flowcomm"
	"github.com/daimall/tools/curd/flow/v1/services/curdservice"
)

// CURDController operations for model
type CURDController struct {
	BaseController
}

// URLMapping ...
func (c *CURDController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Action", c.Action)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("UISetting", c.UISetting)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("DeleteList", c.DeleteList)
}

// Post ...
// @Title 创建一条流程
// @Description create Service
// @Param	service	string 	Service	true		"body for Service content"
// @Success 201 {int} Service
// @Failure 403 body is empty
// @router /:service [post]
func (c *CURDController) Post() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionCreate, oplog)
	}()
	serviceId, ret, oplog, err = c.Service.New(c.uname, c.BaseController.BaseController)
}

// Action ...
// @Title 独立动作
// @Description handle a action
// @Param	body		body 	params	true		"body for Service content"
// @Success 201 {int} OK(step info)
// @Failure 403 body is empty
// @router /:service/:id/:action [post]
func (c *CURDController) Action() {
	var err error
	var ret interface{}
	var serviceId int64
	var oplog string
	var action curdservice.Action
	var actionType = c.Ctx.Input.Param(":action")
	defer func() {
		oplogSep := strings.Split(oplog, curdservice.LogSep)
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
	if actionApp, ok := c.Service.(curdservice.ActionInf); ok {
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

// GetOne ...
// @Title Get One
// @Description get Service by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} Service
// @Failure 403 :id is empty
// @router /:service/:id [get]
func (c *CURDController) GetOne() {
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

// UISetting ...
// @Title Get CRUD UI setting data
// @Description 初始化增删改查的界面设置功能，用于设置BSTable组件
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} Service
// @Failure 403 setting is empty
// @router /:service/uisetting [get]
func (c *CURDController) UISetting() {
	var err error
	var oplog string
	var serviceId int64
	var ret interface{}
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionGetOne, oplog)
	}()
	if uiapp, ok := c.Service.(curdservice.UIInf); ok {
		ret, oplog, err = uiapp.GetUISetting(c.uname, c.BaseController.BaseController)
		return
	}
	err = fmt.Errorf("GetUISetting method is not implement")
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
func (c *CURDController) GetAll() {
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
	//amis orderBy=id&orderDir=desc
	// orderBy = sortBy
	// orderDir = order
	// limit = perPage
	// offset = (page-1) * perPage
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
	// 适配amis
	if v, err := c.GetInt64("perPage"); err == nil {
		limit = v
		if v, err := c.GetInt64("page"); err == nil {
			offset = (v - 1) * limit
		}
	}

	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		v = strings.Replace(v, ".", "__", -1)
		sortby = strings.Split(v, ",")
	}

	// 适配amis
	if v := c.GetString("orderBy"); v != "" {
		sortby = []string{v}
	}

	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}

	// 适配amis
	if v := c.GetString("orderDir"); v != "" {
		order = []string{v}
	}

	var keepMap = map[string]struct{}{
		"orderDir": {},
		"orderBy":  {},
		"page":     {},
		"perPage":  {},
	}
	if beego.AppConfig.DefaultString("webKind", "BS") == "AMIS" {
		// query: k|type=v,v,v  k|type:v|v|v  其中Type可以没有,默认值是 MultiText
		kv := c.Ctx.Request.URL.Query()
		for kInit, v1 := range kv {
			if _, ok := keepMap[kInit]; ok {
				continue
			}
			vInit := v1[0]
			qcondtion := new(common.QueryConditon)
			key_type := strings.Split(kInit, "|") // 解析key中的type信息
			if len(key_type) == 2 {
				qcondtion.QueryKey = key_type[0]
				qcondtion.QueryType = key_type[1]
			} else if len(key_type) == 1 {
				qcondtion.QueryKey = key_type[0]
				qcondtion.QueryType = common.MultiText
			} else {
				logs.Error("Error: invalid query key|type format," + kInit)
				c.JSONResponse(flowcomm.QueryCondErr)
				return
			}
			qcondtion.QueryValues = strings.Split(vInit, ",") // 解析出values信息
			//qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
			query = append(query, qcondtion)
		}
	} else {
		// query: k|type:v|v|v,k|type:v|v|v  其中Type可以没有,默认值是 MultiText
		if v := c.GetString("query"); v != "" {
			for _, cond := range strings.Split(v, ",") { // 分割多个查询key
				qcondtion := new(common.QueryConditon)
				kv := strings.SplitN(cond, ":", 2)
				if len(kv) != 2 {
					logs.Error("query condtion format error:%s, need key:value", kv)
					c.JSONResponse(flowcomm.QueryCondErr)
					return
				}
				kInit, vInit := kv[0], kv[1]          // 初始分割查询key和value（备注，value是多个用|分割）
				key_type := strings.Split(kInit, "|") // 解析key中的type信息
				if len(key_type) == 2 {
					qcondtion.QueryKey = key_type[0]
					qcondtion.QueryType = key_type[1]
				} else if len(key_type) == 1 {
					qcondtion.QueryKey = key_type[0]
					qcondtion.QueryType = common.MultiText
				} else {
					logs.Error("Error: invalid query key|type format," + kInit)
					c.JSONResponse(flowcomm.QueryCondErr)
					return
				}
				qcondtion.QueryValues = strings.Split(vInit, "|") // 解析出values信息
				//qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
				query = append(query, qcondtion)
			}
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
func (c *CURDController) Put() {
	var err error
	var ret interface{}
	var oplog string
	var serviceId int64
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionPut, oplog)
	}()

	// 指定更新字段
	var fields []string
	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
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
	ret, oplog, err = c.Service.Update(serviceId, fields, c.BaseController.BaseController)
}

// Delete ...
// @Title Delete
// @Description delete the Service
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:service/:id [delete]
func (c *CURDController) Delete() {
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
func (c *CURDController) DeleteList() {
	var err error
	var ret interface{}
	var oplog string
	var serviceId int64
	defer func() {
		c.ResponseJSON(err, ret, serviceId, ServiceActionDeleteList, oplog)
	}()
	ret, oplog, err = c.Service.MultiDelete(strings.Split(c.GetString("ids"), ","))
}

// Import ...
// @Title 导出excel，批量创建对象
// @Description batch create Service
// @Param	service	string 	Service	true		"body for Service content"
// @Success 201 {int} Service
// @Failure 403 body is empty
// @router /:service/import [post]
func (c *CURDController) Import() {
	var err error
	var ret interface{}
	var oplog string
	defer func() {
		c.ResponseJSON(err, ret, 0, ServiceActionDeleteList, oplog)
	}()
	if importApp, ok := c.Service.(curdservice.Import); ok {
		var importFile io.Reader
		if importFile, _, err = c.GetFile("importFile"); err != nil {
			logs.Error("get GetFile:importFile failed", err.Error())
			return
		}
		ret, oplog, err = importApp.Import(c.uname, importFile, c.BaseController.BaseController)
		return
	}
	err = fmt.Errorf("import interface not implement")
}

// Export ...
// @Title export
// @Description get Service
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} Service
// @Failure 403
// @router /:service/export [get]
func (c *CURDController) Export() {
	var err error
	var oplog string

	var fields []string
	var sortby []string
	var order []string
	var query []*common.QueryConditon

	defer func() {
		var method string
		pc, _, _, _ := runtime.Caller(1)
		method = runtime.FuncForPC(pc).Name()
		if err == nil {
			// 记录操作日志
			c.LogFunc(0, "export", oplog)
		} else if customErr, ok := err.(customerror.CustomError); ok {
			c.Data["json"] = common.StandRestResult{Code: customErr.GetCode(), Message: customErr.GetMessage()}
			logs.Error("FlowController[%s]%s(customErr)", method, err.Error())
			c.ServeJSON()
		} else {
			c.Data["json"] = common.StandRestResult{Code: -1, Message: err.Error()}
			logs.Error("FlowController[%s]%s", method, err.Error())
			c.ServeJSON()
		}
	}()
	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		v = strings.Replace(v, ".", "__", -1)
		sortby = strings.Split(v, ",")
	}

	// 适配amis
	if v := c.GetString("orderBy"); v != "" {
		sortby = []string{v}
	}

	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}

	// 适配amis
	if v := c.GetString("orderDir"); v != "" {
		order = []string{v}
	}

	var keepMap = map[string]struct{}{
		"orderDir": {},
		"orderBy":  {},
		"page":     {},
		"perPage":  {},
	}
	if beego.AppConfig.DefaultString("webKind", "AMIS") == "AMIS" {
		// query: k|type=v,v,v  k|type:v|v|v  其中Type可以没有,默认值是 MultiText
		kv := c.Ctx.Request.URL.Query()
		for kInit, v1 := range kv {
			if _, ok := keepMap[kInit]; ok {
				continue
			}
			vInit := v1[0]
			qcondtion := new(common.QueryConditon)
			key_type := strings.Split(kInit, "|") // 解析key中的type信息
			if len(key_type) == 2 {
				qcondtion.QueryKey = key_type[0]
				qcondtion.QueryType = key_type[1]
			} else if len(key_type) == 1 {
				qcondtion.QueryKey = key_type[0]
				qcondtion.QueryType = common.MultiText
			} else {
				logs.Error("Error: invalid query key|type format," + kInit)
				c.JSONResponse(flowcomm.QueryCondErr)
				return
			}
			qcondtion.QueryValues = strings.Split(vInit, ",") // 解析出values信息
			//qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
			query = append(query, qcondtion)
		}
	} else {
		// query: k|type:v|v|v,k|type:v|v|v  其中Type可以没有,默认值是 MultiText
		if v := c.GetString("query"); v != "" {
			for _, cond := range strings.Split(v, ",") { // 分割多个查询key
				qcondtion := new(common.QueryConditon)
				kv := strings.SplitN(cond, ":", 2)
				if len(kv) != 2 {
					logs.Error("query condtion format error:%s, need key:value", kv)
					c.JSONResponse(flowcomm.QueryCondErr)
					return
				}
				kInit, vInit := kv[0], kv[1]          // 初始分割查询key和value（备注，value是多个用|分割）
				key_type := strings.Split(kInit, "|") // 解析key中的type信息
				if len(key_type) == 2 {
					qcondtion.QueryKey = key_type[0]
					qcondtion.QueryType = key_type[1]
				} else if len(key_type) == 1 {
					qcondtion.QueryKey = key_type[0]
					qcondtion.QueryType = common.MultiText
				} else {
					logs.Error("Error: invalid query key|type format," + kInit)
					c.JSONResponse(flowcomm.QueryCondErr)
					return
				}
				qcondtion.QueryValues = strings.Split(vInit, "|") // 解析出values信息
				//qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
				query = append(query, qcondtion)
			}
		}
	}

	if exportApp, ok := c.Service.(curdservice.Export); ok {
		var content io.ReadSeeker
		if content, oplog, err = exportApp.Export(c.uname, query, fields, sortby, order); err != nil {
			logs.Error("export failed,", err.Error())
			return
		}
		c.Ctx.ResponseWriter.Header().Add("Content-Disposition", "attachment")
		c.Ctx.ResponseWriter.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		http.ServeContent(c.Ctx.ResponseWriter, c.Ctx.Request, "export", time.Now(), content)
		return
	}
	err = fmt.Errorf("export interface not implement")
}

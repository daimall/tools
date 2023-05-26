package flowcontroller

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/crud/common"
	"github.com/daimall/tools/crud/customerror"
	"github.com/daimall/tools/crud/flow/flowservice"
	"github.com/daimall/tools/crud/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// FlowController operations for model
type FlowController struct {
	// BaseController
	methodMapping map[string]gin.HandlerFunc //method:routertree
}

// Mapping the method to function
func (c *FlowController) Mapping(method string, fn gin.HandlerFunc) {
	if c.methodMapping == nil {
		c.methodMapping = make(map[string]gin.HandlerFunc, 1)
	}
	c.methodMapping[method] = fn
}

// according the method name to get the function
func (c *FlowController) GetMapMethod(method string) gin.HandlerFunc {
	if m, ok := c.methodMapping[method]; ok {
		return m
	}
	panic("method[" + method + "] not in mapping")
}

// Post ...
// @Title 创建一条流程
// @Description create Service
// @Param	service	string 	Service	true		"body for Service content"
// @Success 201 {int} Service
// @Failure 403 body is empty
// @router /:service [post]
func (f *FlowController) Post(c *gin.Context) {
	var err customerror.CustomError
	var ret interface{}
	if crudContextInf, ok := c.Get(common.CRUDContextKey); ok {
		if crudContext, ok := crudContextInf.(flowservice.CRUDContext); ok {
			defer func() {
				c.Set(common.CRUDContextKey, crudContext)
			}()
			crudContext.Action = common.ServiceActionCreate
			if crudContext.ServiceId, ret, crudContext.OperateLog, err = crudContext.Service.New(c); err != nil {
				c.Set(common.CustomErrKey, err)
				return
			}
			// create successfull
			c.Set(common.ResponeDataKey, ret)
			return
		}
	}
	logger.Error("failed to find crud context")
	c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
}

// Action ...
// @Title 独立动作
// @Description handle a action
// @Param	body		body 	params	true		"body for Service content"
// @Success 201 {int} OK(step info)
// @Failure 403 body is empty
// @router /:service/:id/:action [post]
func (f *FlowController) Action(c *gin.Context) {
	var err customerror.CustomError
	var ret interface{}
	var action flowservice.Action
	if crudContextInf, ok := c.Get(common.CRUDContextKey); ok {
		if crudContext, ok := crudContextInf.(flowservice.CRUDContext); ok {
			defer func() {
				c.Set(common.CRUDContextKey, crudContext)
			}()
			crudContext.Action = c.Param("action")
			if crudContext.Action == "" {
				c.Set(common.CustomErrKey, customerror.ActionNotFound)
				logger.Error("action is nil")
				return
			}
			if crudContext.ServiceId != 0 {
				if crudContext.Service, err = crudContext.Service.LoadInst(crudContext, c); err != nil {
					logger.Error("service.LoadInst failed,", err.Error())
					c.Set(common.CustomErrKey, customerror.ServiceLoadFailed)
					return
				}
			}
			if actionApp, ok := crudContext.Service.(flowservice.ActionInf); ok {
				if action, err = actionApp.GetAction(crudContext.ServiceId, crudContext.Action); err != nil {
					logs.Error("GetAction[%s] failed, %s", crudContext.Action, err.Error())
					c.Set(common.CustomErrKey, customerror.ActionNotFound)
					return
				}
			} else {
				c.Set(common.CustomErrKey, customerror.MethodNotImplement)
				logger.Error("GetAction method is not implement")
				return
			}
			if ret, crudContext.ServiceId, crudContext.OperateLog, err = action.Do(c, crudContext); err != nil {
				logger.Error("action.do failed,", err.Error())
				c.Set(common.CustomErrKey, customerror.InternalServerError)
				return
			}
			// action do successfull
			c.Set(common.ResponeDataKey, ret)
			return
		}
	}
	logger.Error("failed to find crud context")
	c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
}

// GetOne ...
// @Title Get One
// @Description get Service by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} Service
// @Failure 403 :id is empty
// @router /:service/:id [get]
func (f *FlowController) GetOne(c *gin.Context) {
	var err customerror.CustomError
	var ret interface{}
	if crudContextInf, ok := c.Get(common.CRUDContextKey); ok {
		if crudContext, ok := crudContextInf.(flowservice.CRUDContext); ok {
			defer func() {
				c.Set(common.CRUDContextKey, crudContext)
			}()
			crudContext.Action = common.ServiceActionGetOne
			if getOneApp, ok := crudContext.Service.(flowservice.GetOneInf); ok {
				if ret, crudContext.OperateLog, err = getOneApp.GetOne(crudContext, c); err != nil {
					c.Set(common.CustomErrKey, err)
					return
				}
				// get on successfull
				c.Set(common.ResponeDataKey, ret)
				return
			}
			logger.Error("get one method is not implemented")
			c.Set(common.CustomErrKey, customerror.MethodNotImplement)
			return
		}
	}
	logger.Error("failed to find crud context")
	c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
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
func (f *FlowController) GetAll(c *gin.Context) {
	var err customerror.CustomError
	var ret struct {
		Items interface{} `json:"items"`
		Total int64       `json:"total"`
	}

	//amis orderBy=id&orderDir=desc
	// orderBy = sortBy
	// orderDir = order
	// limit = perPage
	// offset = (page-1) * perPage
	var crudContextInf interface{}
	var ok bool
	var crudContext flowservice.CRUDContext
	if crudContextInf, ok = c.Get(common.CRUDContextKey); !ok {
		logger.Error("get crudContextInf failed")
		c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
		return
	}
	if crudContext, ok = crudContextInf.(flowservice.CRUDContext); !ok {
		logger.Error("get CRUDContext failed")
		c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
		return
	}
	defer func() {
		crudContext.Action = common.ServiceActionGetAll
		c.Set(common.CRUDContextKey, crudContext)
	}()
	// fields: col1,col2,entity.col3
	if v := c.Query("fields"); v != "" {
		crudContext.Fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := strconv.Atoi(c.Query("limit")); err == nil {
		crudContext.Limit = v
	}
	// offset: 0 (default is 0)
	if v, err := strconv.Atoi(c.Query("offset")); err == nil {
		crudContext.Offset = v
	}
	// 适配amis
	if v, err := strconv.Atoi(c.Query("perPage")); err == nil {
		crudContext.Limit = v
		if v, err := strconv.Atoi(c.Query("page")); err == nil {
			crudContext.Offset = (v - 1) * crudContext.Limit
		}
	}

	// sortby: col1,col2
	if v := c.Query("sortby"); v != "" {
		v = strings.Replace(v, ".", "__", -1)
		crudContext.Sortby = strings.Split(v, ",")
	}

	// 适配amis
	if v := c.Query("orderBy"); v != "" {
		crudContext.Sortby = []string{v}
	}

	// order: desc,asc
	if v := c.Query("order"); v != "" {
		crudContext.Order = strings.Split(v, ",")
	}

	// 适配amis
	if v := c.Query("orderDir"); v != "" {
		crudContext.Order = []string{v}
	}

	var keepMap = map[string]struct{}{
		"orderDir": {},
		"orderBy":  {},
		"page":     {},
		"perPage":  {},
	}
	if viper.GetString("web.webKind") == "AMIS" {
		// query: k|type=v,v,v  k|type:v|v|v  其中Type可以没有,默认值是 MultiText
		kv := c.Request.URL.Query()
		for kInit, v1 := range kv {
			if _, ok := keepMap[kInit]; ok {
				continue
			}
			vInit := v1[0]
			if len(strings.TrimSpace(vInit)) == 0 {
				continue
			}
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
				c.Set(common.CustomErrKey, customerror.QueryCondErr)
				return
			}
			qcondtion.QueryValues = strings.Split(vInit, ",") // 解析出values信息
			//qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
			crudContext.QueryConditon = append(crudContext.QueryConditon, qcondtion)
		}
	} else {
		// query: k|type:v|v|v,k|type:v|v|v  其中Type可以没有,默认值是 MultiText
		if v := c.GetString("query"); v != "" {
			for _, cond := range strings.Split(v, ",") { // 分割多个查询key
				qcondtion := new(common.QueryConditon)
				kv := strings.SplitN(cond, ":", 2)
				if len(kv) != 2 {
					logs.Error("query condtion format error:%s, need key:value", kv)
					c.Set(common.CustomErrKey, customerror.QueryCondErr)
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
					c.Set(common.CustomErrKey, customerror.QueryCondErr)
					return
				}
				qcondtion.QueryValues = strings.Split(vInit, "|") // 解析出values信息
				//qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
				crudContext.QueryConditon = append(crudContext.QueryConditon, qcondtion)
			}
		}
	}

	crudContext.Action = common.ServiceActionGetOne
	if getAllApp, ok := crudContext.Service.(flowservice.GetAllInf); ok {
		if ret.Items, ret.Total, crudContext.OperateLog, err = getAllApp.GetAll(crudContext, c); err != nil {
			logs.Error("getall failed,", err.Error())
			c.Set(common.CustomErrKey, customerror.InternalServerError)
			return
		}
		// success
		c.Set(common.ResponeDataKey, ret)
		return
	}
	logger.Error("get all method is not implemented")
	c.Set(common.CustomErrKey, customerror.MethodNotImplement)
}

// Put ...
// @Title Put
// @Description update the Service
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	Service	true		"body for Service content"
// @Success 200 {object} Service
// @Failure 403 :id is not int
// @router /:service/:id [put]
func (f *FlowController) Put(c *gin.Context) {
	var err customerror.CustomError
	var ret interface{}
	if crudContextInf, ok := c.Get(common.CRUDContextKey); ok {
		if crudContext, ok := crudContextInf.(flowservice.CRUDContext); ok {
			defer func() {
				c.Set(common.CRUDContextKey, crudContext)
			}()
			crudContext.Action = common.ServiceActionPut
			if crudContext.ServiceId != 0 {
				if crudContext.Service, err = crudContext.Service.LoadInst(crudContext, c); err != nil {
					c.Set(common.CustomErrKey, err)
					return
				}
			}
			if v := c.Query("fields"); v != "" {
				crudContext.Fields = strings.Split(v, ",")
			}
			if updateApp, ok := crudContext.Service.(flowservice.UpdateInf); ok {
				if ret, crudContext.OperateLog, err = updateApp.Update(crudContext, c); err != nil {
					c.Set(common.CustomErrKey, err)
					return
				}
				// update successfull
				c.Set(common.ResponeDataKey, ret)
				return
			}
		}
	}
	logger.Error("failed to find crud context")
	c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
}

// Delete ...
// @Title Delete
// @Description delete the Service
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:service/:id [delete]
func (fs *FlowController) Delete(c *gin.Context) {
	var err customerror.CustomError
	var ret interface{}
	if crudContextInf, ok := c.Get(common.CRUDContextKey); ok {
		if crudContext, ok := crudContextInf.(flowservice.CRUDContext); ok {
			defer func() {
				c.Set(common.CRUDContextKey, crudContext)
			}()
			crudContext.Action = common.ServiceActionDelete
			if deleteApp, ok := crudContext.Service.(flowservice.DeleteInf); ok {
				if ret, crudContext.OperateLog, err = deleteApp.Delete(crudContext, c); err != nil {
					c.Set(common.CustomErrKey, err)
					return
				}
				// delete successfull
				c.Set(common.ResponeDataKey, ret)
				return
			}
		}
	}
	logger.Error("failed to find crud context")
	c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
}

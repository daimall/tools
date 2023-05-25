package flowcontroller

import (
	"github.com/daimall/tools/crud/common"
	"github.com/daimall/tools/crud/customerror"
	"github.com/daimall/tools/crud/flow/flowservice"
	"github.com/daimall/tools/crud/logger"
	"github.com/gin-gonic/gin"
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
	return c.methodMapping[method]
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
		if crudContext, ok := crudContextInf.(CRUDContext); ok {
			if crudContext.ServiceId, ret, crudContext.OperateLog, err = crudContext.Service.New(c); err != nil {
				c.Set(common.CustomErrKey, err)
				c.Next()
				return
			}
			// create successfull
			c.Set(common.ResponeDataKey, ret)
			c.Set(common.CRUDContextKey, crudContext)
			c.Next()
			return
		}
	}
	logger.Error("failed to find crud context")
	c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
	c.Next()
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
		if crudContext, ok := crudContextInf.(CRUDContext); ok {
			if getOneApp, ok := crudContext.Service.(flowservice.GetOneInf); ok {
				if ret, crudContext.OperateLog, err = getOneApp.GetOne(c); err != nil {
					c.Set(common.CustomErrKey, err)
					c.Next()
					return
				}
				// get on successfull
				c.Set(common.ResponeDataKey, ret)
				c.Set(common.CRUDContextKey, crudContext)
				c.Next()
				return
			}
			logger.Error("get one method is not implemented")
			c.Set(common.CustomErrKey, customerror.MethodNotImplement)
			c.Next()
			return
		}
	}
	logger.Error("failed to find crud context")
	c.Set(common.CustomErrKey, customerror.CRUDContextNotFound)
	c.Next()
}

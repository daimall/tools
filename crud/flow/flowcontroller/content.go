package flowcontroller

import "github.com/daimall/tools/crud/flow/flowservice"

type CRUDContext struct {
	UserName    string
	ServiceName string
	ServiceId   uint
	Service     flowservice.FlowService
	OperateLog  string
}

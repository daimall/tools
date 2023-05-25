package flowservice

import "github.com/daimall/tools/crud/common"

type CRUDContext struct {
	UserName    string
	ServiceName string
	ServiceId   uint
	Service     FlowService
	OperateLog  string
	Action      string

	// getall 参数
	Fields        []string
	Sortby        []string
	Order         []string
	QueryConditon []*common.QueryConditon
	Limit         int
	Offset        int
}

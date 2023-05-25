package flowservice

type CRUDContext struct {
	UserName    string
	ServiceName string
	ServiceId   uint
	Service     FlowService
	OperateLog  string
	Action      string
}

package flowservice

import (
	"reflect"

	"github.com/daimall/tools/curd/common"
	"gorm.io/gorm"
)

// flow step 接口
type FlowHandler interface {
	// 保存用户的处理结果
	Do(uname string, c common.BaseController) (ret interface{}, oplog string, err error)
	// 从数据库中加载数据，返回一个新实例
	LoadInst(flow FlowService, uname string, id int64) (handler FlowHandler, err error)
	LoadStepHandlers(tx *gorm.DB, flowId int64, stepKey string) (handlers []FlowHandler, err error)
	GetConclusion() bool // 是否通过
	IsFinish() bool      // 是否完成审核
	TableName() string   // 表名
}

// flow step 接口
type FlowStep interface {
	// 唯一标识
	Key() (stepKey string)
	// 通过标准，100 表示所有人结论都要是通过
	PassRate() (rate int)
	// 找到注册handle
	Hander() (handler FlowHandler)
	// 流程步骤等操作信息
	LoadHandlers(tx *gorm.DB, flowId int64) (handlers []FlowHandler, err error)
	// handlers 入库
	AddHandlers(tx *gorm.DB, handlers []FlowHandler) (err error)
	// 退回上一步，清理handlers
	ClearHandlers(tx *gorm.DB, flowId int64, flow string, steps []string, remainingIds []int64) (err error)
	// 从flow属性中获取责任人信息，多个逗号分割，作为步骤初始化时创建handers
	GetDefaultHandlers(flow FlowService) (handers string)
	// 获取处理当前步骤所需要对参数对象（供前端使用）
	GetConfigs(flow FlowService) (ret interface{})
}

//FlowService 流程接口，支持各种流程定制
type FlowService interface {
	// 获取ID
	GetID() int64
	// 获取flowname
	GetFlowName() string
	// 获取处理当前流程所需要对参数对象（供前端使用）
	GetConfigs(uname string, id int64, c common.BaseController) (ret interface{}, oplog string, err error)
	// 新建流程
	New(uname string, c common.BaseController) (flowId int64, ret interface{}, oplog string, err error)
	// make新实例
	NewInst() (flowService FlowService)
	// 获取新实例（从数据库中加载初始值）
	LoadInst(flowId int64) (flowService FlowService, err error)
	// 设置流程状态, 1: 流程草稿   2: 流程流转中  3: 流程完成（被拒绝）， 4: 流程完成（超时） 5： 流程完成（正常）
	SetState(tx *gorm.DB, state int) (err error)
	// 流程步骤（下一步/上一步）， remainingIds：上一步需要重新处理的，默认全部需要
	GoNext(tx *gorm.DB, remainingIds []int64) (handlers []FlowHandler, err error)
	// 获取当前步骤
	GetCurStep() (step FlowStep, err error)
	// 获取下一步
	GetNextStep() (step FlowStep, err error)
	// 获取上一步
	GetPreStep() (step FlowStep, err error)
	// 获取某个属性值（供step使用）
	GetValue(attr string) reflect.Value
	// 更新属性信息（供step使用）
	SetValue(tx *gorm.DB, attr string, refvalue reflect.Value) (err error)
	// 获取一条流程的详情
	GetOne(uname string, id int64) (ret interface{}, oplog string, err error)
	// 获取所有流程
	GetAll(uname string, query []*common.QueryConditon, fields []string,
		sortby []string, order []string, offset int64,
		limit int64) (ret interface{}, count int64, oplog string, err error)
	// 删除一个流程
	Delete(int64) (ret interface{}, oplog string, err error)
	// 删除多个流程
	MultiDelete([]string) (ret interface{}, oplog string, err error)
	OplogModel(uname, flow string, flowid int64, action, remark string) interface{}
}

// 更新字段接口
type Action interface {
	Do(uname string, serviceId int64, actionType string, c common.BaseController) (ret interface{}, oplog string, err error)
}

type ActionInf interface {
	// 获取一个自定义动作
	GetAction(serviceId int64, actionType string) (action Action, err error)
}

type OpHistoryInf interface {
	// 获取操作状态
	GetOpHistory() (ret interface{}, oplog string, err error)
}

type OpLogHistoryInf interface {
	// 获取操作日志历表
	GetOpLogHistory() (ret interface{}, oplog string, err error)
}

type UpdateInf interface {
	// 更新流程
	Update(flowid int64, c common.BaseController) (ret interface{}, oplog string, err error)
}

type PreHandlersInf interface {
	// 获取上一步操作者
	GetPreHandlers() (ret interface{}, oplog string, err error)
}

// 缓存数据的接口
type DataCacheInf interface {
	SetData(string, interface{})
	GetData(string) interface{}
}

// flows 各种类型的流程集合
var services = make(map[string]FlowService)

//Register 注册新类型的服务
func Register(serviceType string, service FlowService) {
	if service == nil {
		panic("service: Register error, service is nil")
	}
	if _, ok := services[serviceType]; ok {
		panic("service: Register called twice for flow " + serviceType)
	}
	services[serviceType] = service
}

// GetService 获取Servie 对象
func GetService(serviceType string) FlowService {
	if v, ok := services[serviceType]; ok {
		return v.NewInst()
	}
	panic("service does not exist: " + serviceType)
}

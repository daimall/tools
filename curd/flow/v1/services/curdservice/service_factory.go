package curdservice

import (
	"io"

	"github.com/daimall/tools/curd/common"
)

// 更新字段接口
type Action interface {
	Do(uname string, serviceId int64, actionType string, c common.BaseController) (ret interface{}, oplog string, err error)
}

//Service 流程接口，支持各种Service定制
type CrudService interface {
	// 新建对象
	New(uname string, c common.BaseController) (flowId int64, ret interface{}, oplog string, err error)
	// make新实例
	NewInst() (crudService CrudService)
	// 获取新实例（从数据库中加载初始值）
	LoadInst(flowId int64) (crudService CrudService, err error)

	// 获取一个对象详情
	GetOne(uname string, id int64) (ret interface{}, oplog string, err error)
	// 获取所有对象列表
	GetAll(uname string, query []*common.QueryConditon, fields []string,
		sortby []string, order []string, offset int64,
		limit int64) (ret interface{}, count int64, oplog string, err error)

	// 更新对象
	Update(flowid int64, fields []string, c common.BaseController) (ret interface{}, oplog string, err error) // 刷新流程基础信息
	// 删除一个对象
	Delete(int64) (ret interface{}, oplog string, err error)
	// 删除多个对象
	MultiDelete([]string) (ret interface{}, oplog string, err error)
}

type UIInf interface {
	// 获取界面初始化功能
	GetUISetting(uname string, c common.BaseController) (ret interface{}, oplog string, err error)
}

type ActionInf interface {
	// 获取一个自定义动作
	GetAction(serviceId int64, actionType string) (action Action, err error)
}

// 日志表自定义接口
type OplogModelInf interface {
	// 返回操作日志记录对象（主要是确定表名）
	OplogModel(uname, flow string, flowid int64, action, remark string) interface{}
}

// 导入接口
type Import interface {
	// 导入操作
	Import(uname string, importFile io.Reader, c common.BaseController) (ret interface{}, oplog string, err error)
}

// 导出接口
type Export interface {
	// 返回excel文件连接
	Export(uname string, query []*common.QueryConditon, fields []string,
		sortby []string, order []string) (content io.ReadSeeker, oplog string, err error)
}

// flows 各种类型的流程集合
var services = make(map[string]CrudService)

//Register 注册新类型的服务
func Register(serviceType string, service CrudService) {
	if service == nil {
		panic("service: Register error, service is nil")
	}
	if _, ok := services[serviceType]; ok {
		panic("service: Register called twice for flow " + serviceType)
	}
	services[serviceType] = service
}

// GetService 获取Servie 对象
func GetService(serviceType string) CrudService {
	if v, ok := services[serviceType]; ok {
		return v
	}
	panic("service does not exist: " + serviceType)
}

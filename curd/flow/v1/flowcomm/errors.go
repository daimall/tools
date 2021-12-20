package flowcomm

import "github.com/daimall/tools/curd/customerror"

var (
	NotFound             = customerror.New(404, "NOT_FOUND")
	ServerErr            = customerror.New(500, "服务器错误")
	UnameNotFound        = customerror.New(1000, "用户认证信息不存在")
	QueryCondErr         = customerror.New(1001, "查询条件不正确")
	StepTypeNotFound     = customerror.New(1002, "Flow Step 不存在")
	ServiceIdNotInt      = customerror.New(1003, "Service Id格式不正确")
	ParamsErr            = customerror.New(1004, "参数错误")
	UploadErr            = customerror.New(1005, "上传文件失败")
	UpdateActionNotFound = customerror.New(1006, "Update Acton 不存在")
	ActionNotFound       = customerror.New(1007, "Action 不存在")
)

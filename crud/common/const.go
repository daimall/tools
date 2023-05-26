package common

const (
	// Token 常量
	TokenKey           = "Authorization"
	UserNameSessionKey = "X-UserName"
	CodeKey            = "X-Code"
	AppTag             = "APPID"

	// CRUD 常量(存储到gin keys中)
	CustomErrKey   = "customerr"
	ResponeDataKey = "responeData"
	UserNameKey    = "username"
	CRUDContextKey = "CRUDContext"
	// ServiceNameKey = "servicename"
	// ServiceKey      = "service"
	// ServiceIdKey    = "serviceId"
	// OperationLogKey = "operationlog"
)

// Service action 常量
const (
	ServiceActionGetAll     = "GetAll"
	ServiceActionGetOne     = "GetOne"
	ServiceActionCreate     = "Create"
	ServiceActionPut        = "Update"
	ServiceActionDelete     = "Delete"
	ServiceActionDeleteList = "DeleteList"
	ServiceActionImport     = "Import"
	ServiceActionExport     = "Export"

	ServiceActionGetConfigs = "GetConfigs"
	ServiceOpList           = "OperationList"
)

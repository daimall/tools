package flowservice

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/crud/customerror"
	"github.com/daimall/tools/crud/dbmysql/dbgorm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CheckHandler struct {
	ID          uint   `gorm:"primary_key"  json:"id"`                                        // 自增主键
	User        string `gorm:"size:100;column:user;unique_index:onerecord"  json:"user"`      // 操作者的账户id
	Step        string `gorm:"size:20;column:step;unique_index:onerecord"  json:"step"`       // 步骤
	ServiceId   uint   `gorm:"column:service_id;unique_index:onerecord"  json:"serviceId"`    // FlowId
	ServiceName string `gorm:"size:50;column:service;unique_index:onerecord"  json:"service"` // FlowType
	Conclusion  int    `gorm:"column:conclusion"  json:"conclusion"`                          // 评审结论,0未评审，1 通过，2 风险通过 3 拒绝, 100 转他人处理
	Remark      string `gorm:"column:remark"  json:"remark"`                                  // 操作详情
	//Attach      string      `gorm:"size:1000;column:attach;index"  json:"attach"` // 附件
	CreatedAt     time.Time   `json:"created_at"` // 创建时间
	UpdatedAt     time.Time   `json:"updated_at"` // 最后更新时间
	TblName       string      `gorm:"-" json:"-"` // 表名
	Flow          FlowService `gorm:"-" json:"-"` // 关联的流程
	AttachKeys    []string    `gorm:"-" json:"-"` // 附件key列表
	IsPublicCheck bool        `gorm:"-" json:"-"` // 是否为公共审核，如果是不需要鉴权

	Pretreatment   func(CRUDContext, *gin.Context, *gorm.DB, *CheckHandler) (oplog string, err error)                `gorm:"-" json:"-"` // 预处理
	Aftertreatment func(CRUDContext, *gin.Context, *gorm.DB, *CheckHandler, []FlowHandler) (oplog string, err error) `gorm:"-" json:"-"` // 后处理
	// Attachs []attach.Attach `gorm:"-" json:"attachs"` // 附件列表

	PreStepHandlerIds []uint `gorm:"-" json:"preStepHandlerIds"` // 上一步需要重新审批的用户
}

func (h *CheckHandler) Do(crudContext CRUDContext, c *gin.Context) (ret interface{}, oplog string, err customerror.CustomError) {

	return
}

// 从数据库中加载数据，返回一个新实例
func (h *CheckHandler) LoadInst(flow FlowService, uname string, id uint) (ret FlowHandler, err error) {
	handler := &CheckHandler{ID: id}
	if id == 0 {
		// 查询默认handle
		handler.ServiceId = flow.GetID()
		handler.ServiceName = flow.GetFlowName()
		var step FlowStep
		if curApp, ok := flow.(GetCurStepInf); ok {
			if step, err = curApp.GetCurStep(); err != nil {
				logs.Error("get cur step failed,", err.Error(), flow)
				return
			}
		} else {
			err = fmt.Errorf("GetCurStepInf is not implement")
			return
		}
		handler.Step = step.Key()
	}
	dbInst := dbgorm.GetDBInst()
	if err = dbInst.Table(h.TableName()).Where(handler).First(handler).Error; err != nil {
		logs.Error("get handler failed,", err.Error(), handler)
		return
	}
	handler.TblName = h.TblName
	handler.Flow = flow
	handler.Pretreatment = h.Pretreatment
	handler.AttachKeys = h.AttachKeys
	handler.Aftertreatment = h.Aftertreatment
	handler.Conclusion = h.Conclusion
	handler.PreStepHandlerIds = h.PreStepHandlerIds
	handler.IsPublicCheck = h.IsPublicCheck

	return handler, nil
}

// 是否通过
func (h *CheckHandler) GetConclusion() bool {
	// 通过或者风险通过视为通过
	return h.Conclusion == ConclusionGo || h.Conclusion == ConclusionGoWithRisk
}

// 是否完成审核
func (h *CheckHandler) IsFinish() bool {
	return h.Conclusion != 0
}

// 表名
func (h *CheckHandler) TableName() string {
	if h.TblName != "" {
		return h.TblName
	}
	return "step_handlers"
}

// 加载handler
func (h *CheckHandler) LoadStepHandlers(tx *gorm.DB, flowId uint, stepKey string) (handlers []FlowHandler, err error) {
	l := []CheckHandler{}
	err = tx.Table(h.TableName()).
		Where(&CheckHandler{ServiceId: flowId, TblName: h.TableName(), Step: stepKey}).
		Find(&l).Error
	if len(l) == 0 {
		return nil, fmt.Errorf("no hander found for step[%s],flowId:[%d]", stepKey, flowId)
	}
	handlers = make([]FlowHandler, len(l))
	for i := range l {
		handlers[i] = &l[i]
	}
	return
}

// 非接口方法
// 设置附件keys，用于保存附件使用
func (h *CheckHandler) SetAttachKeys(keys []string) {
	h.AttachKeys = keys
}

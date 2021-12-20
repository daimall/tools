package common

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/curd/customerror"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

const ControllerKindJSON string = "JSON"

//BaseController ...
//所有控制controller的基础struct
type BaseController struct {
	REST             StandRestResultInf
	Kind             string // 类型 JSON 或非 JSON
	beego.Controller        // beego 基础控制器
}

func (b *BaseController) GetStandRestResult() StandRestResultInf {
	if b.REST == nil {
		return StandRestResult{}
	}
	return b.REST
}

// JSONResponse 返回JSON格式结果
func (c *BaseController) JSONResponse(err customerror.CustomError, data ...interface{}) {
	if err != nil {
		c.Data["json"] = c.GetStandRestResult().GetStandRestResult(err.GetCode(), err.GetMessage(), nil)
		logs.Error("JSONResponse:", err.Error())
	} else {
		if len(data) == 1 {
			c.Data["json"] = c.GetStandRestResult().GetStandRestResult(0, "OK", data[0])
		} else {
			c.Data["json"] = c.GetStandRestResult().GetStandRestResult(0, "OK", data)
		}
	}
	c.ServeJSON()
}

//ValidateParameters obj must pointer, json Unmarshal object and require parameter validate
func (c *BaseController) ValidateParameters(obj interface{}) customerror.CustomError {
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, obj); err != nil {
		logs.Error("UnmarshalRequestBody Err", err.Error())
		return ParamsError
	}
	if err := validate.Struct(obj); err != nil {
		logs.Error("validate params err:", err.Error())
		return ParamsValidateError
	}
	return nil

}

// SaveToFileCustomName saves uploaded file to new path with custom name.
// it only operates the first one of mutil-upload form file field.
func (c *BaseController) SaveToFileCustomName(fromfile string, fc func(*multipart.FileHeader) string) error {
	file, h, err := c.Ctx.Request.FormFile(fromfile)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(fc(h), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

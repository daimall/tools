package flowservice

import (
	"reflect"

	"github.com/daimall/tools/crud/customerror"
)

// 获取可读取/可修改的反射值
func GetRefValue(v reflect.Value) reflect.Value {
	if v.CanInterface() {
		switch v.Kind() {
		case reflect.Ptr:
			return GetRefValue(v.Elem())
		case reflect.Struct:
			return v
		case reflect.Array, reflect.Slice:
			return GetRefValue(v.Index(0))
		}
	}
	return reflect.Value{}
}

func Err2CustomErr(err error) customerror.CustomError {
	if err == nil {
		return nil
	}
	return customerror.New(customerror.InternalServerErrorCODE, err.Error())
}

package dbutil

import (
	"github.com/astaxie/beego/logs"
	dbgorm "github.com/daimall/tools/curd/dbmysql/dbgorm/mysql"
	"gorm.io/gorm"
)

var (
	DBInst *gorm.DB
)

func init() {
	// 初始化DB链接
	if DBInst == nil {
		DBInst = dbgorm.GetDBInst()
	}
}

// GetDBInst 获取DB实例
func GetDBInst() *gorm.DB {
	if DBInst == nil {
		DBInst = dbgorm.GetDBInst()
	}
	return DBInst
}

func CheckTable(table interface{}) {
	if !DBInst.HasTable(table) {
		if err := DBInst.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(table).Error; err != nil {
			logs.Error(err.Error())
			panic(err)
		}
	}
}

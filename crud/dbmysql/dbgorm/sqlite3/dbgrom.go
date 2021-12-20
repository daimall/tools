package dbgorm

import (
	// "lottery/app/sqlite"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "modernc.org/sqlite"
)

var db *gorm.DB

// 获取
func GetDBInst() *gorm.DB {
	dbfile := beego.AppConfig.String("DB::SourceName")
	var err error
	if db == nil {
		if db, err = gorm.Open(sqlite.Dialector{DSN: dbfile, DriverName: "sqlite"}, &gorm.Config{}); err != nil {
			logs.Error("open lottery.db failed,", err.Error())
			panic(err.Error())
		}
	}
	return db
}

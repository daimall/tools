package dbgorm

import (
	// "lottery/app/sqlite"

	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "modernc.org/sqlite"
)

var db *gorm.DB

// 获取
func GetDBInst() *gorm.DB {
	dbfile := beego.AppConfig.String("DB::SourceName")

	newLogger := logger.New(
		logs.GetLogger(), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)

	var err error
	if db == nil {
		if db, err = gorm.Open(sqlite.Dialector{DSN: dbfile, DriverName: "sqlite"}, &gorm.Config{
			Logger: newLogger,
		}); err != nil {
			logs.Error("open lottery.db failed,", err.Error())
			panic(err.Error())
		}
	}
	return db
}

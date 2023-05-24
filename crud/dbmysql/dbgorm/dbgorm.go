package dbgorm

import (
	"fmt"
	"time"

	"github.com/daimall/tools/aes/cbc"
	"github.com/daimall/tools/crud/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

var db *gorm.DB

// 获取 gorm数据库链接实例
func GetDBInst() *gorm.DB {
	if db == nil {
		db = NewDBInst()
	}
	return db
}

func NewDBInst() *gorm.DB {
	return newDBInst("DB")
}

func NewDBInstBySection(section string) *gorm.DB {
	return newDBInst(section)
}

func newDBInst(section string) *gorm.DB {
	var db *gorm.DB
	var err error
	newLogger := gormlog.New(
		&zapLogger{}, // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		gormlog.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  gormlog.Info, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	)
	if viper.GetString(section+".driver") == "mysql" {
		passwd := viper.GetString(section + ".password")
		if pwdEncryptKey := viper.GetString(section + ".pwdEncryptKey"); pwdEncryptKey != "" {
			// 密码是加密形态，需要解密
			if passwd, err = cbc.New(pwdEncryptKey).Decrypt(passwd); err != nil {
				logger.Error("Decrypt db passwd failed, pwdKey: %s, ciphertext: %s, err:%s",
					pwdEncryptKey, passwd, err.Error())
				panic(err)
			}
		}
		dsn := fmt.Sprintf(viper.GetString(section+".sourceName"), passwd)
		if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		}); err != nil {
			logger.Error(err.Error())
			panic(err)
		}
		return db
	} else if viper.GetString(section+".driver") == "sqlite3" {
		dbfile := viper.GetString(section + ".sourceName")
		if db == nil {
			if db, err = gorm.Open(sqlite.Dialector{DSN: dbfile, DriverName: "sqlite"}, &gorm.Config{
				Logger: newLogger,
			}); err != nil {
				logger.Error("open %s.db failed, %s", dbfile, err.Error())
				panic(err.Error())
			}
		}
		return db
	} else {
		logger.Error("unkown Driver")
	}
	return db
}

type zapLogger struct{}

func (l *zapLogger) Printf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

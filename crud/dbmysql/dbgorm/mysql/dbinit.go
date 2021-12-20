package dbgorm

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/daimall/tools/aes/cbc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	var db *gorm.DB
	passwd := beego.AppConfig.String("DB::Passwd")
	if pwdEncryptKey := beego.AppConfig.String("DB::PwdEncryptKey"); pwdEncryptKey != "" {
		// 密码是加密形态，需要解密
		var err error
		if passwd, err = cbc.New(pwdEncryptKey).Decrypt(passwd); err != nil {
			logs.Error("Decrypt db passwd failed, pwdKey: %s, ciphertext: %s, err:%s",
				pwdEncryptKey, passwd, err.Error())
			panic(err)
		}
	}
	dsn := fmt.Sprintf(beego.AppConfig.String("DB::SourceName"), passwd)
	var err error
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		logs.Error(err.Error())
		panic(err)
	}
	db.SetLogger(logs.GetLogger())
	db.LogMode(beego.AppConfig.DefaultBool("dblog", true))
	return db
}

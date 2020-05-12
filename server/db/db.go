package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kyour-cn/gourd/common"
	"log"
)

var db *gorm.DB

var config Config

type Config struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
	User string `toml:"user"`
	Pwd  string `toml:"pwd"`
	Db   string `toml:"dbname"`
}

func InitDb() {

	var err error

	err = common.ReadConfig("database", &config)
	if err != nil {
		log.Printf("Db 配置错误：%v\n", err)
	}

	//连接数据库
	db, err = gorm.Open("mysql", config.User+":"+config.Pwd+
		"@("+config.Host+":"+config.Port+")/"+config.Db+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		//panic(err)
		log.Printf("Db 连接错误：%v\n", err)

	}

	//user := model.User{}
	//
	//Db.First(&user)
	//
	//log.Printf("查询到User：%v", user)

}

func Conn() *gorm.DB {
	return db
}

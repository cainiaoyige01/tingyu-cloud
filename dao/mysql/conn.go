package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"tingyu-cloud/lib"
	"tingyu-cloud/log"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/26 16:17
 * @Description: 数据的初始化
 */

// Db 初始化全局的Db对象
var Db *gorm.DB
var err error

// InitDb 初始化数据库
func InitDb(conf lib.ServerConfig) {
	//拼接完成的MySQL链接路径 使用fmt.Sprintf()
	dbParams := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		conf.DbUser,
		conf.DbPassword,
		conf.DbHost,
		conf.DbPort,
		conf.DbName,
	)
	//返回数据对象 参数1 是数据库类型 参数2 是数据库连接路径
	Db, err = gorm.Open("mysql", dbParams)
	if err != nil {
		log.Logger.Errorln("open", err)
	}

	//全局禁用表明复数
	Db.SingularTable(true)
	//这个函数可以处理创建表的问题
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}
	//这是数据库的最大空闲链接数
	Db.DB().SetMaxIdleConns(10)
	//这是数据库的最大链接数量
	Db.DB().SetMaxOpenConns(100)
	log.Logger.Infoln("Setting up database")
}

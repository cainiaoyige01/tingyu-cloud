package main

import (
	"fmt"
	"github.com/spf13/viper"
	"tingyu-cloud/dao/mysql"
	"tingyu-cloud/dao/redis"
	"tingyu-cloud/lib"
	"tingyu-cloud/log"
	"tingyu-cloud/router"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/14 16:12
 * @Description: 启动处
 */
func main() {
	v := viper.New()
	// 加载配置文件 使用vip配置
	v.SetConfigFile("conf/config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("err is %v\n", err)
	}
	//读取配置信息 绑定到model去
	config := lib.LoadServerConfig(v)
	log.InitLLog()
	// 初始化数据库的信息
	mysql.InitDb(config)
	//记得关闭数据库连接
	defer mysql.Db.Close()
	//初始化redis的链接
	redis.InitRedis(config)
	// 路由的对象
	r := router.SetupRouter()
	// 加载静态资源
	r.LoadHTMLGlob("view/*")
	// 映射静态资源
	r.Static("/static", "./static")
	// 开启监听端口
	err = r.Run(":9090")
	if err != nil {
		log.Logger.Errorln("Server start is fails", err)
	}
}

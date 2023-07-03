package lib

import "github.com/spf13/viper"

/**
 * @Author: _niuzai
 * @Date:   2023/6/26 15:45
 * @Description:
 */

// ServerConfig 服务端配置数据结构
type ServerConfig struct {
	//运行级别
	RunMode string
	//Location
	Location string
	//redis
	RedisHost string
	RedisPort string
	//database
	DbType     string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
	//QQ 配置信息
	AppId          string
	AppKey         string
	AppRedirectUri string
	//阿里云配置
	AccessKeyId     string
	AccessKeySecret string
	EndPoint        string
	BucketName      string
}

// LoadServerConfig 加载服务端的配置
func LoadServerConfig(v *viper.Viper) ServerConfig {
	//这是简单的写法！严谨的写法的是把每一个字段的先写出来 然后判断是否有错误 才能继续的
	return ServerConfig{
		RunMode:         v.GetString("runMode"),
		Location:        v.GetString("app.location"),
		RedisHost:       v.GetString("redis.host"),
		RedisPort:       v.GetString("redis.port"),
		DbType:          v.GetString("database.type"),
		DbUser:          v.GetString("database.user"),
		DbPassword:      v.GetString("database.password"),
		DbHost:          v.GetString("database.host"),
		DbPort:          v.GetString("database.port"),
		DbName:          v.GetString("database.name"),
		AppId:           v.GetString("qq.app_id"),
		AppKey:          v.GetString("qq.redirect_uri"),
		AppRedirectUri:  v.GetString("qq.app_key"),
		AccessKeyId:     v.GetString("oos.access_key_id"),
		AccessKeySecret: v.GetString("oos.access_key_secret"),
		EndPoint:        v.GetString("oos.end_point"),
		BucketName:      v.GetString("oos.bucket_name"),
	}
}

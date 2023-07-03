package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"tingyu-cloud/lib"
	"tingyu-cloud/log"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/26 21:40
 * @Description:redis
 */

// 定义全局redis池
var RedisPool *redis.Pool

func InitRedis(conf lib.ServerConfig) {
	//配置redis连接池
	RedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			dial, err := redis.Dial("tcp", conf.RedisHost+":"+conf.RedisPort)
			if err != nil {
				log.Logger.Errorln("redis dial error: ", err)
				return nil, err
			}
			//选择redis存储盘
			_, err = dial.Do("SELECT", "1")
			if err != nil {
				_ = dial.Close()
				return nil, err
			}
			return dial, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { //一个测试链接可用性
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         5,                 //idle的列表长度，空闲的线程数
		MaxActive:       0,                 //线程池的最大连接数，0表示没有限制
		IdleTimeout:     200 * time.Second, //最大的空闲的等待时间，超过这个时间，空闲的链接将会被关闭
		Wait:            true,              //当连接池已满，是否要阻塞等待获取连接。fasle表示不必等待，直接返回错误
		MaxConnLifetime: 0,
	}
	log.Logger.Infoln("redis init on port", conf.RedisHost+":"+conf.RedisPort)
}

// SetKey 可以设置过期时间key-value 0 表示永久性的存储
func SetKey(key, value interface{}, expires int) error {
	//获取连接
	conn := RedisPool.Get()
	defer conn.Close()
	if expires == 0 {
		_, err := conn.Do("SET", key, value)
		return err
	} else {
		_, err := conn.Do("SETEX", key, expires, value)
		return err
	}
}

// GetKey 根据key获取值
func GetKey(key string) (string, error) {
	//获取连接
	conn := RedisPool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

// DeleteKey 根据key删除value
func DeleteKey(key string) error {
	//获取连接
	conn := RedisPool.Get()
	//操作命令
	_, err := conn.Do("DEL", key)
	defer conn.Close()
	return err
}

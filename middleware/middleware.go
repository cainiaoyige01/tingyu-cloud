package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tingyu-cloud/dao/mysql"
	"tingyu-cloud/dao/redis"
	"tingyu-cloud/log"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/27 11:43
 * @Description: 中间件
 */

// CheckLogin 验证是否登录了
func CheckLogin(ctx *gin.Context) {
	//去cookie中去数据
	cookie, err := ctx.Cookie("Token")
	if err != nil {
		log.Logger.Errorln("Error getting cookie token is error", err)
		//重定向到登录处 302
		ctx.Redirect(http.StatusFound, "/")
		ctx.Abort()
	}
	//与redis进行校验 是否存在数据
	openId, err := redis.GetKey(cookie)
	if err != nil {
		log.Logger.Errorln("Error getting cookie key is error", err)
		//重定向到登录处 302
		ctx.Redirect(http.StatusFound, "/")
		ctx.Abort()
	}
	//去数据库查询openId是否存在  这里其实可以优化不必去数据进行查询
	ok := mysql.QueryUserExists(openId)
	if ok {
		//使用context 传递openId
		ctx.Set("openId", openId)
		ctx.Next()
	} else {
		//校验失败，返回登录处
		ctx.Redirect(http.StatusFound, "/")
	}

}

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tingyu-cloud/dao/mysql"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/30 16:47
 * @Description:帮助
 */

func Help(ctx *gin.Context) {
	openId := ctx.GetString("openId")
	user := mysql.GetUserInfoById(openId)

	//获取用户文件使用明细数量
	fileDetailUse := mysql.GetFileDetailUser(user.FileStoreId)

	ctx.HTML(http.StatusOK, "help.html", gin.H{
		"currHelp":      "active",
		"user":          user,
		"fileDetailUse": fileDetailUse,
	})
}

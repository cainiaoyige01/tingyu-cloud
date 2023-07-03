package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tingyu-cloud/dao/mysql"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/30 11:34
 * @Description: 其他分类
 */

func OtherFiles(ctx *gin.Context) {
	// 获取openId
	openId := ctx.GetString("openId")
	// openId 获取用户数据
	user := mysql.GetUserInfoById(openId)
	// 获取用户使用明细数量
	detailUser := mysql.GetFileDetailUser(user.FileStoreId)
	//count := detailUser["docCount"]
	// 根据文件类型 查出文件
	imageFiles := mysql.GetFileByTypeId(5, user.FileStoreId)
	// 响应会前端
	ctx.HTML(http.StatusOK, "other-files.html", gin.H{
		"user":          user,
		"fileDetailUse": detailUser,
		"otherFiles":    imageFiles,
		"otherCount":    len(imageFiles),
		"currOther":     "active",
		"currClass":     "active",
	})
}

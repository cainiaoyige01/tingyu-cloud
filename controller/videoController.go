package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tingyu-cloud/dao/mysql"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/30 11:32
 * @Description: 视频分类
 */

func VideoFiles(ctx *gin.Context) {
	// 获取openId
	openId := ctx.GetString("openId")
	// openId 获取用户数据
	user := mysql.GetUserInfoById(openId)
	// 获取用户使用明细数量
	detailUser := mysql.GetFileDetailUser(user.FileStoreId)
	//count := detailUser["docCount"]
	// 根据文件类型 查出文件
	imageFiles := mysql.GetFileByTypeId(3, user.FileStoreId)
	// 响应会前端
	ctx.HTML(http.StatusOK, "video-files.html", gin.H{
		"user":          user,
		"fileDetailUse": detailUser,
		"videoFiles":    imageFiles,
		"videoCount":    len(imageFiles),
		"currVideo":     "active",
		"currClass":     "active",
	})
}

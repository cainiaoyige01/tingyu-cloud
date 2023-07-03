package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tingyu-cloud/dao/mysql"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/30 10:47
 * @Description:docFile 文件分类
 */

// DocFiles 文件分类：我的文档
func DocFiles(ctx *gin.Context) {
	// 获取openId
	openId := ctx.GetString("openId")
	// openId 获取用户数据
	user := mysql.GetUserInfoById(openId)
	// 获取用户使用明细数量
	detailUser := mysql.GetFileDetailUser(user.FileStoreId)
	//count := detailUser["docCount"]
	// 根据文件类型 查出文件
	docFiles := mysql.GetFileByTypeId(1, user.FileStoreId)
	// 响应会前端
	ctx.HTML(http.StatusOK, "doc-files.html", gin.H{
		"user":          user,
		"fileDetailUse": detailUser,
		"docFiles":      docFiles,
		"docCount":      len(docFiles),
		"currDoc":       "active",
		"currClass":     "active",
	})
}

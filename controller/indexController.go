package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tingyu-cloud/dao/mysql"
	"tingyu-cloud/log"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/27 12:03
 * @Description:页面加载信息
 */

// Index 页面展示
func Index(ctx *gin.Context) {
	//获取用户openId
	openId, ok := ctx.Get("openId")
	if ok {
		log.Logger.Errorln("openId parameter not set", ok)
	}
	//根据openId 获取用户信息 头像、name
	user := mysql.GetUserInfoById(openId.(string))
	//获取用户仓库的信息
	userFileStore := mysql.GetUserFileStore(user.Id)
	//获取用户文件数量 去表my_file查询
	fileCount := mysql.GetUserFileCount(user.FileStoreId)
	//获取用户文件夹数量
	fileFolderCount := mysql.GetFileFolderCount(user.FileStoreId)
	//获取用户文件使用明细数量 也就是图片、视频、音乐等占比情况
	fileDetailUse := mysql.GetFileDetailUser(user.FileStoreId)
	//响应数据
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"user":            user,
		"currIndex":       "active",
		"userFileStore":   userFileStore,
		"fileCount":       fileCount,
		"fileFolderCount": fileFolderCount,
		"fileDetailUse":   fileDetailUse,
	})
}

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tingyu-cloud/dao/mysql"
	"tingyu-cloud/log"
	"tingyu-cloud/utils"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/27 15:05
 * @Description: 文件夹的展示
 */

// Files 展示文件夹
func Files(ctx *gin.Context) {
	//从context中获取openId
	openId, _ := ctx.Get("openId")
	//获取页面传过来的参数 没有传过来就是用默认的
	fId := ctx.DefaultQuery("fId", "0")
	//通过openId获取用户信息 进行后续查询
	user := mysql.GetUserInfoById(openId.(string))
	//获取当前目录所有文件 也就是根据去查询fId 利用的my_file表
	getUserFile := mysql.GetUserFile(user.FileStoreId, fId)
	//获取当前目录所有文件夹 0级所拥有的文件夹
	getFileFolder := mysql.GetFileFolder(user.FileStoreId, fId)
	//获取父级的文件夹信息 id
	getParentFolderId := mysql.GetParentFileFolder(fId)
	//获取当前目录所有父级
	currentAllParent := mysql.GetCurrentAllParent(getParentFolderId, make([]mysql.FileFolder, 0))
	//获取当前的目录信息
	currentFolder := mysql.GetCurrentFolder(fId)
	//获取用户文件使用明细数量
	detailUser := mysql.GetFileDetailUser(user.FileStoreId)
	ctx.HTML(http.StatusOK, "files.html", gin.H{
		"currAll":          "active",
		"user":             user,
		"fId":              currentFolder.Id,
		"fName":            currentFolder.FileFolderName,
		"files":            getUserFile,
		"fileFolder":       getFileFolder,
		"parentFolder":     getParentFolderId,
		"currentAllParent": currentAllParent,
		"fileDetailUse":    detailUser,
	})
}

// DownloadFile 下载文件
func DownloadFile(ctx *gin.Context) {

	// 获取fId 文件的ID
	fId := ctx.Query("fId")
	// 取数据中回去文件
	file := mysql.GetFileInfo(fId)
	if file.FileHash == "" {
		log.Logger.Errorln("文件不存在：")
		return
	}
	// 从oss获取文件
	fileData := utils.DownloadOss(file.FileHash, file.Postfix)
	// 下载次数加一
	mysql.DownloadNumAdd(fId)

	// 响应给前端
	ctx.Header("Content-disposition", "attachment;filename=\""+file.FileName+file.Postfix+"\"")
	ctx.Data(http.StatusOK, "application/octect-stream", fileData)

}

// DeleteFile 删除文件
func DeleteFile(ctx *gin.Context) {
	// 获取openId user信息
	openId := ctx.GetString("openId")
	user := mysql.GetUserInfoById(openId)
	// 获取文件id
	fId := ctx.DefaultQuery("fId", "")
	// 获取父级目录
	folderId := ctx.Query("folder")
	if fId == "" {
		return
	}
	// 删除数据库文件数据
	mysql.Delete(fId, folderId, user.FileStoreId)

	// 响应数据给前端
	ctx.Redirect(http.StatusMovedPermanently, "/cloud/files?fid="+folderId)
}

//

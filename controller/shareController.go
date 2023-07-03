package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lifei6671/gocaptcha"
	"net/http"
	"strconv"
	"strings"
	"tingyu-cloud/dao/mysql"
	"tingyu-cloud/log"
	"tingyu-cloud/utils"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/30 17:29
 * @Description: 分享页面
 */

// ShareFile 创建分线文件
func ShareFile(ctx *gin.Context) {
	// 去到openId
	openId := ctx.GetString("openId")
	// 用户信息
	user := mysql.GetUserInfoById(openId)

	// 获取参数
	id := ctx.Query("id")
	url := ctx.Query("url")
	// 获取内容 设置验证码
	code := gocaptcha.RandText(4)

	fileId, _ := strconv.Atoi(id)

	// 获取分享的内容
	hash := mysql.CreateShare(code, user.UserName, fileId)

	// 响应给前端
	ctx.JSON(http.StatusOK, gin.H{
		"url":  url + "?f=" + hash,
		"code": code,
	})
}

// SharePass 分享页面
func SharePass(ctx *gin.Context) {
	// 获取url中的参数
	hash := ctx.Query("f")

	// 获取分享信息
	shareInfo := mysql.GetShareInfoByHash(hash)
	// 获取文件信息
	file := mysql.GetFileInfo(strconv.Itoa(shareInfo.FileId))
	ctx.HTML(http.StatusOK, "share.html", gin.H{
		"id":       shareInfo.FileId,
		"username": shareInfo.Username,
		"fileType": file.Type,
		"filename": file.FileName + file.Postfix,
		"hash":     shareInfo.Hash,
	})
}

// DownloadShareFile 下载分享文件
func DownloadShareFile(ctx *gin.Context) {
	// 获取参数中id code hash
	fileId := ctx.Query("id")
	code := ctx.Query("code")
	hash := ctx.Query("hash")

	// 获取文件信息
	fileInfo := mysql.GetFileInfo(fileId)

	// 校验验证码
	ok := mysql.VerifyShareCode(fileId, strings.ToLower(code))
	if !ok {
		log.Logger.Infoln("验证码信息错误")
		ctx.Redirect(http.StatusMovedPermanently, "/file/share?f="+hash)
		return
	}
	// 从oss获取文件
	fileData := utils.DownloadOss(fileInfo.FileHash, fileInfo.Postfix)
	// 下次次数加一
	mysql.DownloadNumAdd(fileId)
	ctx.Header("Content-disposition", "attachment;filename=\""+fileInfo.FileName+fileInfo.Postfix+"\"")
	ctx.Data(http.StatusOK, "application/octect-stream", fileData)
}

package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"tingyu-cloud/dao/mysql"
	"tingyu-cloud/log"
	"tingyu-cloud/utils"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/28 12:57
 * @Description: 上传页面
 */

// Upload 文件上传页面
func Upload(ctx *gin.Context) {
	//从context中获取openId
	openId, _ := ctx.Get("openId")
	//获取页面传过来的参数 没有传过来就是用默认的
	fId := ctx.DefaultQuery("fId", "0")
	//通过openId获取用户信息 进行后续查询
	user := mysql.GetUserInfoById(openId.(string))
	//获取当前目录所有文件信息
	fileFolders := mysql.GetFileFolder(user.FileStoreId, fId)
	//获取父级的文件夹信息 id
	getParentFolderId := mysql.GetParentFileFolder(fId)
	//获取当前目录所有父级
	currentAllParent := mysql.GetCurrentAllParent(getParentFolderId, make([]mysql.FileFolder, 0))
	//获取当前的目录信息
	currentFolder := mysql.GetCurrentFolder(fId)
	//获取用户文件使用明细数量
	detailUser := mysql.GetFileDetailUser(user.FileStoreId)
	ctx.HTML(http.StatusOK, "upload.html", gin.H{
		"user":             user,
		"currUpload":       "active",
		"fId":              currentFolder.Id,
		"fName":            currentFolder.FileFolderName,
		"fileFolders":      fileFolders,
		"parentFolder":     getParentFolderId,
		"currentAllParent": currentAllParent,
		"fileDetailUse":    detailUser,
	})
}

// UploadFile 上传文件
func UploadFile(ctx *gin.Context) {
	// 获取openId 用户信息
	openId, _ := ctx.Get("openId")
	user := mysql.GetUserInfoById(openId.(string))
	// 获取当前目录id
	Fid := ctx.GetHeader("id")
	// 从前端中接受上传文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Logger.Errorln("文件上传错误:", err)
	}

	// 判断当前文件夹是否有同名的存在
	ok := mysql.CurrFileExists(Fid, header.Filename)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 501,
		})
		return
	}
	// 判断用户容量是否充足
	ok = mysql.CapacityIsEnough(header.Size, user.FileStoreId)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 503,
		})
		return
	}
	defer file.Close()

	// 文件保存到本地路径
	location := "D:\\download\\" + header.Filename

	// 在本地创建一个新的文件
	create, err := os.Create(location)
	if err != nil {
		log.Logger.Errorln("创建文件失败：", err)
	}
	defer create.Close()
	// 将上传的文件拷贝至新创建的文件中
	fileSize, err := io.Copy(create, file)
	if err != nil {
		log.Logger.Errorln("文件拷贝失败", err)
	}
	//将光标一直开头
	_, _ = create.Seek(0, 0)
	fileHash := utils.GetSHA256HashCode(create)
	// 通过hash判断数据是否存在 其实可以为图片去唯一的名字就不需要考虑这种了
	ok = mysql.FileOssExists(fileHash)
	if ok {
		// 上传至阿里云oss
		utils.UploadOss(header.Filename, fileHash)
	}

	// 新建文件信息 存到数据库中
	mysql.CreateFile(header.Filename, fileHash, fileSize, Fid, user.FileStoreId)

	// 上传成功减去相应剩余容量
	mysql.SubtractSize(fileSize/1024, user.FileStoreId)
	// 响应数据给前端
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

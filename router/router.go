package router

import (
	"github.com/gin-gonic/gin"
	"tingyu-cloud/controller"
	"tingyu-cloud/middleware"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/26 19:58
 * @Description: router
 */

// SetupRouter 路由管理
func SetupRouter() *gin.Engine {
	router := gin.Default()
	//登录拦截
	router.GET("/", controller.Login)
	//扫描跳转登录
	router.GET("/toLogin", controller.HandlerLogin)
	//登录成功 获取token
	router.GET("/qqLogin", controller.GetQQToken)
	// 获取分享内容
	router.GET("/file/share", controller.SharePass)
	// 现在分享内容
	router.GET("/file/shareDownload", controller.DownloadShareFile)
	//建立一个路由组
	cloud := router.Group("cloud")
	//中间件判断是否登录了
	cloud.Use(middleware.CheckLogin)
	{
		//首页页面需要加载的信息
		cloud.GET("index", controller.Index)
		//文件的页面
		cloud.GET("/files", controller.Files)
		//文件上传页面
		cloud.GET("/upload", controller.Upload)

		//对文件进行分类了
		cloud.GET("/doc-files", controller.DocFiles)
		cloud.GET("/image-files", controller.Images)
		cloud.GET("/video-files", controller.VideoFiles)
		cloud.GET("/music-files", controller.MusicFiles)
		cloud.GET("/other-files", controller.OtherFiles)

		// 下载文件
		cloud.GET("/downloadFile", controller.DownloadFile)
		// 删除文件
		cloud.GET("/deleteFile", controller.DeleteFile)
		// 删除目录
		cloud.GET("/deleteFolder", controller.DeleteFileFolder)
		// 帮助
		cloud.GET("/help", controller.Help)
		// 退出登录
		cloud.GET("/logout", controller.Logout)
	}
	// post 请求
	{
		// 上传文件
		cloud.POST("/uploadFile", controller.UploadFile)
		// 添加文件夹
		cloud.POST("/addFolder", controller.AddFolder)
		// 更新目录名
		cloud.POST("/updateFolder", controller.UpdateFileFolder)
		// 分享文件
		cloud.POST("/getQrCode", controller.ShareFile)
	}

	return router
}

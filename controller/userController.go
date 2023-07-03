package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"tingyu-cloud/dao/mysql"
	"tingyu-cloud/dao/redis"
	"tingyu-cloud/log"
	"tingyu-cloud/utils"
)

// PrivateInfo 绑定token
type PrivateInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"open_id"`
}
type UserInfo struct {
	NickName    string `json:"nickname"`
	FigureUrlQQ string `json:"figureurl_qq"`
}

// Login 登录页面
func Login(ctx *gin.Context) {
	//跳转到登录页面
	ctx.HTML(http.StatusOK, "login.html", nil)
}

// HandlerLogin 扫描跳转登录
func HandlerLogin(ctx *gin.Context) {
	//触发函数进行跳转 也就是获取二维码！注意后面的回调函数！这个很重要的
	url := "https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=101827468&state=test&redirect_uri=http://127.0.0.1:9090/qqLogin"
	log.Logger.Info(url)
	//状态码301
	ctx.Redirect(http.StatusMovedPermanently, url)
}

// GetQQToken 获取QQ传过来的access token！手机端确认后服务端自动
func GetQQToken(ctx *gin.Context) {
	//获取参数中code "/qqLogin?code=55F9804AFB09D9610738310A09C0804F&state=test"
	code := ctx.Query("code")
	log.Logger.Infoln(code)
	//拼接完整的路径 请求到QQ中去 二维码
	loginUrl := "https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=101827468&client_secret=0d2d856e48e0ebf6b98e0d0c879fe74d&redirect_uri=http://127.0.0.1:9090/callbackQQ&code=" + code
	log.Logger.Infoln(loginUrl)
	//转发loginUrl获取响应
	resp, err := http.Get(loginUrl)
	if err != nil {
		log.Logger.Errorln("login url is not valid", err)
	}
	defer resp.Body.Close()
	//转化成字节切片
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Errorln("err reading response body", err)
	}
	//转化成字符串类型
	body := string(bs)
	//resultMap包含access token 、token的刷新时间以及token过期时间
	resultMap := utils.ConverToMap(body)
	//获取access token 、token的刷新时间以及token过期时间 使用struct来数据绑定传输
	info := &PrivateInfo{}
	info.AccessToken = resultMap["access_token"]
	info.RefreshToken = resultMap["refresh_token"]
	info.ExpiresIn = resultMap["expires_in"]
	//获取OpenId
	GetOpenId(info, ctx)
}

// GetOpenId 凭借access_token去取openId等用户信息
func GetOpenId(info *PrivateInfo, ctx *gin.Context) {
	//发送请求去QQ获取信息参数
	resp, err := http.Get(fmt.Sprintf("%s?access_token=%s", "https://graph.qq.com/oauth2.0/me", info.AccessToken))
	if err != nil {
		log.Logger.Errorln("GetOpenId error", err)
	}
	defer resp.Body.Close()
	all, _ := ioutil.ReadAll(resp.Body)
	body := string(all)
	//获取到OpenId
	info.OpenId = body[45:77]
	//这里就可以获取QQ的用户信息
	GetUserInfo(info, ctx)
}

// GetUserInfo 获取QQ用户信息
func GetUserInfo(info *PrivateInfo, ctx *gin.Context) {
	//添加参数到url中去
	params := url.Values{}
	params.Add("access_token", info.AccessToken)
	params.Add("openid", info.OpenId)
	params.Add("oauth_consumer_key", "101827468")

	uri := fmt.Sprintf("https://graph.qq.com/user/get_user_info?%s", params.Encode())
	//获取信息内容
	resp, err := http.Get(uri)
	if err != nil {
		log.Logger.Errorln("err resp is ", err)
	}
	defer resp.Body.Close()
	//resp包含头像、名字等信息 先转成切片先
	body, _ := ioutil.ReadAll(resp.Body)

	LoginSucceed(string(body), info.OpenId, ctx)
}

// LoginSucceed 登录成功了，处理登录
func LoginSucceed(str string, openId string, ctx *gin.Context) {
	// 获取QQ name、QQ号
	var userInfo UserInfo
	//转成json格式
	err := json.Unmarshal([]byte(str), &userInfo)
	if err != nil {
		log.Logger.Errorln("装换成json失败", err)
	}
	//生成token存储于http头部 用于判断是否登录了!
	//使用md5加密 有：token+时间戳+openId
	hashToken := utils.EncodeMd5("token" + string(time.Now().Unix()) + openId)
	//存到redis中去 需要加上一个过期的时间
	err = redis.SetKey(hashToken, openId, 3600*24)
	if err != nil {
		log.Logger.Errorln("error setting key for redis key", err)
	}
	//设置cookie中去
	ctx.SetCookie("Token", hashToken, 3600*24, "/", "127.0.0.1:9090", false, true)
	//查询数据库中 如果存在openId就直接跳转到页面 否则存到数据再跳转
	ok := mysql.QueryUserExists(openId)
	//响应数据
	if ok {
		//登录成功定向到首页
		ctx.Redirect(http.StatusMovedPermanently, "/cloud/index")
	} else {
		mysql.CreateUser(openId, userInfo.NickName, userInfo.FigureUrlQQ)
		ctx.Redirect(http.StatusMovedPermanently, "/cloud/index")
	}

}

// DeleteFileFolder 删除目录
func DeleteFileFolder(ctx *gin.Context) {
	// 前段传来目录id
	id := ctx.DefaultQuery("fId", "")
	if id == "" {
		return
	}
	// 获取文件夹的信息 返回数据给前端去到父级目录的重定向
	folderInfo := mysql.GetCurrentFolder(id)
	// 关键的是如果文件夹里面有内容则怎么办？
	// 删除文件夹并删除文件夹中文件信息
	mysql.DeleteFileFolder(id)
	ctx.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+strconv.Itoa(folderInfo.ParentFolderId))
}

// Logout 退出登录
func Logout(ctx *gin.Context) {
	// 获取cookie
	token, err := ctx.Cookie("Token")
	if err != nil {
		log.Logger.Errorln("无法获取到Cookie", err)
	}
	// 去redis 中删除key
	err = redis.DeleteKey(token)
	if err != nil {
		log.Logger.Errorln("在redis删除key错误", err)
	}
	// 修改cookie的值
	ctx.SetCookie("Token", "", 0, "/", "127.0.0.1:9090", false, false)

	// 重定向到登录页面
	ctx.Redirect(http.StatusFound, "/")
}

// AddFolder 创建文件夹
func AddFolder(ctx *gin.Context) {
	// 获取openId
	openId := ctx.GetString("openId")
	user := mysql.GetUserInfoById(openId)
	// 获取父级目录
	parentId := ctx.DefaultPostForm("parentFolderId", "0")
	// 获取新建目录名字
	folderName := ctx.PostForm("fileFolderName")
	// 新建文件夹数据
	mysql.CreateFolder(folderName, parentId, user.FileStoreId)

	// 获取父文件夹信息
	parentFolder := mysql.GetParentFileFolder(parentId)
	// 刷新目录
	ctx.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+parentId+"&fName="+parentFolder.FileFolderName)
}

// UpdateFileFolder 跟新目录
func UpdateFileFolder(ctx *gin.Context) {
	// 获取目录id
	fileFolderId := ctx.PostForm("fileFolderId")
	// 需要更新的名字
	fileFolderName := ctx.PostForm("fileFolderName")
	fileFolder := mysql.GetCurrentFolder(fileFolderId)
	// 更新
	mysql.UpdateFolderName(fileFolderId, fileFolderName)
	// 刷新
	ctx.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+strconv.Itoa(fileFolder.ParentFolderId))
}

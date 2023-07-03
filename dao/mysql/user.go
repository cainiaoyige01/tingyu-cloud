package mysql

import (
	"time"
	"tingyu-cloud/log"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/26 22:07
 * @Description:
 */

// User 用户表
type User struct {
	Id           int
	OpenId       string
	FileStoreId  int
	UserName     string
	RegisterTime time.Time
	ImagePath    string
}

// QueryUserExists 根据openId查询用户是否存在
func QueryUserExists(openId string) bool {
	var user User
	Db.Find(&user, "open_id=?", openId)
	if user.Id == 0 {
		return false
	}
	return true
}

// CreateUser 创建用户信息
func CreateUser(openId string, qqName string, qqImage string) {
	user := User{
		OpenId:       openId,
		FileStoreId:  0,
		UserName:     qqName,
		RegisterTime: time.Now(),
		ImagePath:    qqImage,
	}
	// 创建表
	Db.Create(&user)
	log.Logger.Infoln("openId", openId, "username", qqName)
	fileStore := FileStore{
		UserId:      user.Id,
		CurrentSize: 0,
		MaxSize:     1048576,
	}
	Db.Create(&fileStore)
	user.FileStoreId = fileStore.Id
	Db.Save(&user)
}

// GetUserInfoById 通过id 获取用户信息
func GetUserInfoById(openId string) (user User) {
	find := Db.Find(&user, "open_id=?", openId)
	log.Logger.Infoln(find)
	return
}

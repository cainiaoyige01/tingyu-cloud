package mysql

import (
	"strings"
	"time"
	"tingyu-cloud/utils"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/30 17:36
 * @Description: 分享数据
 */

// Share 分享表
type Share struct {
	Id       int
	Code     string
	FileId   int
	Username string
	Hash     string
}

// CreateShare 创建分享内容
func CreateShare(code string, name string, fileId int) string {
	share := Share{
		Code:     strings.ToLower(code),
		FileId:   fileId,
		Username: name,
		Hash:     utils.EncodeMd5(code + string(time.Now().Unix())),
	}
	Db.Create(&share)
	return share.Hash
}

// GetShareInfoByHash 根据hash获取分享的内容
func GetShareInfoByHash(hash string) (share Share) {
	Db.Find(&share, "hash=?", hash)
	return
}

// VerifyShareCode 校验验证码
func VerifyShareCode(fileId string, code string) bool {
	var share Share
	Db.Find(&share, "file_id=? and code = ?", fileId, code)
	if share.Id == 0 {
		return false
	}
	return true
}

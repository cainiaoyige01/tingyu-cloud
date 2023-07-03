package mysql

import "tingyu-cloud/log"

/**
 * @Author: _niuzai
 * @Date:   2023/6/27 14:03
 * @Description:file_store表的操作
 */

// FileStore 表
type FileStore struct {
	Id          int
	UserId      int
	CurrentSize int64
	MaxSize     int64
}

// GetUserFileStore 获取用户仓库容量等信息
func GetUserFileStore(userId int) (fileStore FileStore) {
	find := Db.Find(&fileStore, "user_id=?", userId)
	log.Logger.Infoln(find)
	return
}

// CapacityIsEnough 判断容量是否充足
func CapacityIsEnough(fileSize int64, fileStoreId int) bool {
	var fileStore FileStore
	Db.First(&fileStore, fileStoreId)
	if fileStore.MaxSize-(fileSize/1024) < 0 {
		return false
	}

	return true
}

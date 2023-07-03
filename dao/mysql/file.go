package mysql

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"
	"tingyu-cloud/log"
	"tingyu-cloud/utils"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/27 14:02
 * @Description: my_file表数据库操作
 */

// MyFile 表其实严格来说需要加上`db:xxx`的
type MyFile struct {
	Id             int
	FileName       string //文件名
	FileHash       string //文件哈希值
	FileStoreId    int    //文件仓库id
	FilePath       string //文件存储路径
	DownloadNum    int    //下载次数
	UploadTime     string //上传时间
	ParentFolderId int    //父文件夹id
	Size           int64  //文件大小
	SizeStr        string //文件大小单位
	Type           int    //文件类型
	Postfix        string //文件后缀
}

// GetUserFileCount 获取用户上传的文件数量
func GetUserFileCount(fileStoreId int) (fileCount int64) {
	//初始化一个切片
	var myFile []MyFile
	//根据fileStoreId查询数据
	count := Db.Find(&myFile, "file_store_id=?", fileStoreId).Count(&fileCount)
	log.Logger.Infoln("用户文件数量为count:", count)
	//返回
	return
}

// GetFileDetailUser 查看每一种类中文件的数量占比情况
func GetFileDetailUser(fileStoreId int) map[string]int64 {
	// 需要切片去接受
	var myFile []MyFile
	//定义每一种类型的长度
	var (
		docCount   int64 //文件类型
		imgCount   int64 //图片类型jpg
		videoCount int64 //视频类型mp4
		musicCount int64 //音乐类型
		otherCount int64 //其他文件类型的
	)
	//初始化map存储每一种长度用于返回使用
	fileDetailUserMap := make(map[string]int64, 0)
	//去数据查询
	docCount = Db.Find(&myFile, "file_store_id=? AND type=?", fileStoreId, 1).RowsAffected
	imgCount = Db.Find(&myFile, "file_store_id=? AND type=?", fileStoreId, 2).RowsAffected
	videoCount = Db.Find(&myFile, "file_store_id=? AND type=?", fileStoreId, 3).RowsAffected
	musicCount = Db.Find(&myFile, "file_store_id=? AND type=?", fileStoreId, 4).RowsAffected
	otherCount = Db.Find(&myFile, "file_store_id=? AND type=?", fileStoreId, 5).RowsAffected

	fileDetailUserMap["docCount"] = docCount
	fileDetailUserMap["imgCount"] = imgCount
	fileDetailUserMap["videoCount"] = videoCount
	fileDetailUserMap["musicCount"] = musicCount
	fileDetailUserMap["otherCount"] = otherCount

	return fileDetailUserMap
}

// GetUserFile 根据file_store_id和parent_folder_id去查询
func GetUserFile(fileStoreId int, parentFolderId string) (myFile []MyFile) {
	//
	Db.Find(&myFile, "file_store_id = ? and parent_folder_id = ?", fileStoreId, parentFolderId)
	return
}

// GetFileByTypeId 根据文件类型获取文件
func GetFileByTypeId(typeId int, fileStoreId int) (myFile []MyFile) {
	//
	Db.Find(&myFile, "type = ? AND file_store_id = ?", typeId, fileStoreId)
	fmt.Println("--------", Db.Find(&myFile, "type = ? AND file_store_id = ?", typeId, fileStoreId))
	return
}

// CurrFileExists 判断当前文件是否有同名的
func CurrFileExists(fId string, fileName string) bool {
	var file MyFile
	// 获取文件后缀名 path.Ext(name)
	fileSuffix := strings.ToLower(path.Ext(fileName))
	// 获取文件名
	filePrefix := fileName[0 : len(fileName)-len(fileSuffix)]
	Db.Find(&file, "parent_folder_id = ? and file_name = ? and postfix = ?", fId, filePrefix, fileSuffix)
	if file.Size > 0 {
		return false
	}
	return true
}

// FileOssExists 判断文件是否存在
func FileOssExists(fileHash string) bool {
	var file MyFile
	Db.Find(&file, "file_hash = ?", fileHash)
	if file.FileHash != "" {
		return false
	}
	return true
}

// CreateFile 上传文件的信息添加数据库中
func CreateFile(fileName string, fileHash string, fileSize int64, fId string, fileStoreId int) {
	// 需要的数据库：文件名、文件后缀、size_str 是KB还是MB？
	var sizeStr string
	// 获取文件后缀名
	fileSuffix := path.Ext(fileName)
	// 获取文件名
	filePrefix := fileName[0 : len(fileName)-len(fileSuffix)]
	// 注意要fId 转成int类型! 我们使fId查询的时候可以忽略string或int 因为mysql可以进行隐式转换！但是创建新数据就可以了 必须与数据库的字段类型一致
	fid, _ := strconv.Atoi(fId)
	if fileSize > 1048576 {
		sizeStr = strconv.FormatInt(fileSize/1024, 10) + "KB"
	} else {
		sizeStr = strconv.FormatInt(fileSize/102400, 10) + "MB"
	}
	myFile := MyFile{
		FileName:       filePrefix,
		FileHash:       fileHash,
		FileStoreId:    fileStoreId,
		FilePath:       "",
		DownloadNum:    0,
		UploadTime:     time.Now().Format("2006-01-02T15:04:05"),
		ParentFolderId: fid,
		Size:           fileSize / 1024,
		SizeStr:        sizeStr,
		Type:           utils.GetFileTypeInt(fileSuffix), // 注意这个类型在数据是int类型的 需要进行转换一下
		Postfix:        strings.ToLower(fileSuffix),
	}
	// 向数据库中插入数据
	Db.Create(&myFile)
}

// SubtractSize 减去剩余容量
func SubtractSize(fileSize int64, fileStoreId int) {
	var fileStore FileStore
	Db.First(&fileStore, fileStoreId)
	fileStore.CurrentSize = fileStore.CurrentSize + fileSize/1024
	fileStore.MaxSize = fileStore.MaxSize - fileSize/1024
	// 根据数据库save 记得& 地址符号 要不然就会错的
	Db.Save(&fileStore)
}

// GetFileInfo 获取要下载文件
func GetFileInfo(fId string) (myFile MyFile) {
	Db.First(&myFile, fId)
	return
}

// DownloadNumAdd 下载次数加一
func DownloadNumAdd(fId string) {
	file := GetFileInfo(fId)
	file.DownloadNum += 1
	Db.Save(&file)
}

// Delete 删除文件
func Delete(id string, parentFileId string, fileStoreId int) {
	Db.Where("id = ? and file_store_id=? and parent_folder_id=?", id, fileStoreId, parentFileId).Delete(MyFile{})
}

package mysql

import (
	"strconv"
	"time"
	"tingyu-cloud/log"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/27 14:04
 * @Description:file_folder表的操作
 */

// FileFolder 文件夹表
type FileFolder struct {
	Id             int
	FileFolderName string
	ParentFolderId int
	FileStoreId    int
	Time           string
}

// GetFileFolderCount 根据file_store_id获取文件夹的数量
func GetFileFolderCount(fileFolderId int) (count int) {
	var fileFolder []FileFolder
	Db.Find(&fileFolder, "file_store_id=?", fileFolder).Count(&count)
	return
}

// GetFileFolder 获取当前文件夹
func GetFileFolder(fileStoreId int, parentFolderId string) (file []FileFolder) {
	Db.Order("time desc").Find(&file, "file_store_id=? AND parent_folder_id=?", fileStoreId, parentFolderId)
	log.Logger.Infoln(Db.Order("time desc").Find(&file, "file_store_id=? AND parent_folder_id=?", fileStoreId, parentFolderId))
	return
}

// GetParentFileFolder 获取父级的文件夹信息 也就是id
func GetParentFileFolder(parentFolderId string) (file FileFolder) {
	Db.Find(&file, "id=?", parentFolderId)
	log.Logger.Infoln(Db.Find(&file, "parent_folder_id=?", parentFolderId))
	return
}

// GetCurrentAllParent 获取当前所有父级的文件夹
func GetCurrentAllParent(folder FileFolder, folders []FileFolder) []FileFolder {
	//定义一个
	var parentFolder FileFolder
	//parentFolderId 不等于0 说明其还有父级
	if folder.ParentFolderId != 0 {
		Db.Find(&parentFolder, "id=?", folder.ParentFolderId)
		folders = append(folders, parentFolder)
		//递归查询所有父级
		return GetCurrentAllParent(parentFolder, folders)
	}
	//翻转切片
	for i := 0; i < len(folders)/2; i++ {
		folders[i], folders[len(folders)-1-i] = folders[len(folders)-1-i], folders[i]
	}
	return folders
}

// GetCurrentFolder 获取当前目录信息
func GetCurrentFolder(fId string) (fileFolder FileFolder) {
	Db.Find(&fileFolder, "id = ?", fId)
	return
}

// DeleteFileFolder 删除文件夹
func DeleteFileFolder(id string) bool {
	var fileFolder1 FileFolder
	var fileFolder2 FileFolder
	// 删除文件夹信息
	Db.Where("id = ?", id).Delete(FileFolder{})
	// 删除文件夹中的信息
	Db.Where("parent_folder_id = ?", id).Delete(MyFile{})
	// 删除文件中文件夹的信息
	Db.Find(&fileFolder1, "parent_folder_id=?", id)
	Db.Where("parent_folder_id = ?", id).Delete(FileFolder{})
	Db.Find(&fileFolder2, "parent_folder_id=?", fileFolder1.Id)
	// 使用递归进行删除
	if fileFolder2.Id != 0 {
		return DeleteFileFolder(strconv.Itoa(fileFolder1.Id))
	}
	return true
}

// CreateFolder 新创建文件夹
func CreateFolder(folderName string, parentId string, fileStoreId int) {
	// parent转成int
	parentIdInt, _ := strconv.Atoi(parentId)

	fileFolder := FileFolder{
		FileFolderName: folderName,
		ParentFolderId: parentIdInt,
		FileStoreId:    fileStoreId,
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	Db.Create(&fileFolder)
}

// UpdateFolderName 更新文件夹的名字
func UpdateFolderName(fileFolderId string, fileFolderName string) {
	var fileFolder FileFolder
	// 更新使用model
	Db.Model(&fileFolder).Where("id=?", fileFolderId).Update("file_folder_name", fileFolderName)

}

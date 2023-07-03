package utils

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"path"
	"tingyu-cloud/log"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/30 13:47
 * @Description: 上传文件到oss
 */

func UploadOss(fileName, fileHash string) {
	// 获取文件后缀
	fileSuffix := path.Ext(fileName)
	// 创建OSS的对象
	client, err := oss.New("https://oss-cn-beijing.aliyuncs.com", "LTAI5tHAYMTqbF8XRVQEHVRj",
		"zmiQB9MzRzD81XEz0LnewpebnsidCs")
	if err != nil {
		log.Logger.Errorln("创建OSS实例错误", err)
	}
	// 获取存储空间
	bucket, err := client.Bucket("banshantingyu")
	if err != nil {
		log.Logger.Errorln("获取存储空间实例错误", err)
	}
	// 上传本地文件
	err = bucket.PutObjectFromFile("files/"+fileHash+fileSuffix, "D:\\download\\"+fileName)
	if err != nil {
		log.Logger.Errorln("本地文件上传失败：", err)
	}
}

// DownloadOss 下载文件
func DownloadOss(fileHash, fileType string) []byte {
	// 创建oss实例
	client, err := oss.New("https://oss-cn-beijing.aliyuncs.com", "LTAI5tHAYMTqbF8XRVQEHVRj", "zmiQB9MzRzD81XEz0LnewpebnsidCs")
	if err != nil {
		log.Logger.Errorln("创建下载文件实例错误：", err)
	}
	// 获取存储空间
	bucket, err := client.Bucket("banshantingyu")
	if err != nil {
		log.Logger.Errorln("获取下载文件的存储空间错误", err)
	}
	// 下载文件到流
	body, err := bucket.GetObject("files/" + fileHash + fileType)
	if err != nil {
		log.Logger.Errorln("下载文件失败：", err)
	}
	// 流需要关闭的
	defer body.Close()

	// 把流转换成byte
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Logger.Errorln("流转换字节切片出现错误", err)
	}
	return data

}

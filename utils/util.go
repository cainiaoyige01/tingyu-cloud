package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/26 20:51
 * @Description: 工具包
 */

// ConverToMap 将body的=号格式字符串为map
func ConverToMap(str string) map[string]string {
	//初始化切片
	resultMap := make(map[string]string)
	//分割切片
	split := strings.Split(str, "&")
	//遍历循环
	for _, value := range split {
		vl := strings.Split(value, "=")
		resultMap[vl[0]] = vl[1]
	}
	return resultMap
}

// EncodeMd5 md5加密
func EncodeMd5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// GetSHA256HashCode SHA256生成哈希值
func GetSHA256HashCode(file *os.File) string {
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	_, _ = io.Copy(hash, file)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode

}

// GetFileTypeInt 判断文件属于那种类型
func GetFileTypeInt(fileSuffix string) int {
	fileSuffix = strings.ToLower(fileSuffix)
	if fileSuffix == ".doc" || fileSuffix == ".docx" || fileSuffix == ".txt" || fileSuffix == ".pdf" {
		return 1
	}
	if fileSuffix == ".jpg" || fileSuffix == ".png" || fileSuffix == ".gif" || fileSuffix == ".jpeg" {
		return 2
	}
	if fileSuffix == ".mp4" || fileSuffix == ".avi" || fileSuffix == ".mov" || fileSuffix == ".rmvb" || fileSuffix == ".rm" {
		return 3
	}
	if fileSuffix == ".mp3" || fileSuffix == ".cda" || fileSuffix == ".wav" || fileSuffix == ".wma" || fileSuffix == ".ogg" {
		return 4
	}

	return 5
}

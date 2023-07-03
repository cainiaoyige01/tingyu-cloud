package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

/**
 * @Author: _niuzai
 * @Date:   2023/6/26 16:36
 * @Description: 日记信息配置
 */

// logger 定义全局日记对象
var Logger *logrus.Logger

// InitLLog 初始化日记
func InitLLog() {
	//建立一个新的日记对象
	Logger = logrus.New()
	//设置日记的级别为debug
	//logger.SetLevel(logrus.DebugLevel)
	//输出到多个方向
	w1 := os.Stdout
	//日记文件写到哪里去
	w2, _ := os.OpenFile("log/demo.log", os.O_CREATE|os.O_WRONLY, 0644)

	Logger.SetOutput(io.MultiWriter(w1, w2))
	//使用默认的text文本的吧 就进行设置了
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Logger.WithField("cloud", "tingyu")

}

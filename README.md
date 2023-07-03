# tingyu-cloud
一个仿百度网盘的分布式存储！仿百度云，支持断点续传和秒传。
主要使用的技术：
  Gin + Gorm + Redis + OSS + Viper + Logrus + docker

### 功能支持：
* QQ扫码登录
* 文件上传。支持单文件、多文件、文件夹上传
* 展示上传的文件在页面
* 可以管理文件
* 文件下载
* 文件上传类型统计展示
* 分享文件给别人
* 在线浏览
* 退出登录
### 页面展示
![image](https://github.com/cainiaoyige01/tingyu-cloud/blob/main/static/img/1.png)
![image](https://github.com/cainiaoyige01/tingyu-cloud/blob/main/static/img/2.png)
![image](https://github.com/cainiaoyige01/tingyu-cloud/blob/main/static/img/3.png)
![image](https://github.com/cainiaoyige01/tingyu-cloud/blob/main/static/img/5.png)

### 分片上传、断点续传
文件过大时，上传文件需要很长时间，且中途退出将导致文件重传。

分片上传: 在上传文件的时候，在本地文件超过5M的大小文件进行分片。在服务端将文件进行组合。
断点续传: 如果文件没有上传完，关闭客户端。再一次上传文件时，对比服务器已经上传的分片，只需要上传没有的分片。

程序重启时，由于不保存目录结构和上传进度。会删除已经上传的文件分片，再次上传从头开始。

### 秒传
每一个文件都会对应的MD5码。当检测上传文件时，本地已经存在相同的MD5的文件，则不需用户上传了！

### 链接分享
后端会生成一条可访问的链接以及提取码！就是类似于百度网盘那种：
链接：http://127.0.0.1:9090/file/share?f=3ce70d5347e0b44ded78e20c20659585 提取码：gTDN

### 项目启动
只需要修改yaml文件中配置：
改一下OSS的key and token！为mysql创建一个数据库！QQ 的权限码可以使用我的Github上的！这也是我从一篇文章章获取的！
感兴趣的可以去看一下我的QQ登录原理文章：https://blog.csdn.net/niuzai_/article/details/131426888

启动：http://localhost:9090  就可以访问了！

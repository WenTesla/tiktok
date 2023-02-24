# 简易版抖音项目

## 项目答辩文档与App地址
[飞书文档](https://xp8kgipb5a.feishu.cn/docx/RPAvdVcqpoc6DzxRPlnc5f3VnMb)
## 项目演示视频
[青训营演示视频](https://www.bilibili.com/video/BV1uT411S79U/?share_source=copy_web&vd_source=fe55b12bbf1a3c973a095834d9f2ba6d)

## 项目启动
### 建表
resources/initial.sql  
resources/insertData.sql  

### 替换redis地址
config/redis.go  


### 直接启动
```shell
cd go 
```
```shell
go run main.go 
```  

## APP 操作
设置服务端地址
为方便测试登录和注册，及修改网络请求的服务器地址，提供了退出登录和高级设置两个能力。
1. 点击退出登录会自动重启
2. 在高级设置中可以配置自己的服务端项目的前缀地址，如下配置的http://192.168.1.7:8080
   在app中访问上述某个接口时就会拼接该前缀地址，例如访问 http://192.168.1.7:8080/douyin/feed/ 拉取视频列表
![](https://xingqiu-tuchuang-1256524210.cos.ap-shanghai.myqcloud.com/12640/20230224155544.png)

## 表数据

## 使用mvc分层结构
参考文献  （不懂的可以参考这篇文章）
https://juejin.cn/post/7152299022017888286
## 项目层次 (-todo)
### controller
### model
### service
### controller


## 版本  
* go版本 1.19  
* mysql 8.0+
* redis驱动
## 使用到的框架与依赖  
+ gin框架
+ gorm框架
+ mysql驱动
+ golang的jwt框架
+ 腾讯云的oss存储（设置了工作流用于截取视频的第一帧(.jpg)并储存在相同的桶中）文献: https://juejin.cn/post/7195857732846567485
+ redis驱动
## 已知错误    
 

### 日志(resources/gin.log)
### 注意
* 启动服务会自动生成日志文件  
* 每次重启会覆盖日志  
* 同时封装了log日志  


## 待优化地方:  
* 建表为了省事使用自增Id，安全性缺乏 （懒得优化了）
* 上传文件相同文件名称的处理(目前将文件名改为时间戳后处理，好像也可以)
* 未设置读写分离
* 一些地方可以用到指针(javer 的问题)
* 服务更加细致，只返回对应的必要的json数据
* 定时任务


## 注意
* **目前上传文件接口只支持mp4格式**

## to-do
* 使用docker部署
* 自动执行sql语句

### 作者:  
bowen https://www.github.com/WenTesla
### 最后修改时间
2023/2/24

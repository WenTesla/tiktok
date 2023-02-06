# 简易版抖音项目  

## 目前状态:    
**开发中**  
## 项目启动
直接启动
```shell
go run main.go
```

## 表数据(todo)
## 使用mvc分层结构
参考文献  （不懂的可以参考这篇文章）
https://juejin.cn/post/7152299022017888286
## 项目层次 (-todo)
### controller
### model
### service
### controller



**统一使用goland开发**，学生可以用学生邮箱去官网认证学生，免费使用

## 版本  
* go版本 1.19  
* mysql 8.0
## 使用到的框架与依赖  
+ gin框架
+ gorm框架
+ mysql驱动
+ golang的jwt框架
+ 腾讯云的oss存储（设置了工作流用于截取视频的第一帧(.jpg)并储存在相同的桶中）文献: https://juejin.cn/post/7195857732846567485

## 已知错误    
* 可能和数据库的创建连接有关  
Error 1040: Too many connections (datasource)数据库连接过多  

### 错误日志（上传后点击才能响应接口）  
[GIN] 2023/02/01 - 17:26:57 | 400 |     11.1264ms |  192.168.31.236 | POST     "/douyin/publish/action/"

## 待优化地方:  
* 数据库查询请求过多
* 建表为了省事使用自增Id，安全性缺乏
* 数据库访问过多，造成数据库压力大,后续使用**redis**优化  
* 上传文件相同文件名称的处理(目前将文件名改为时间戳后处理)
* 未设置读写分离

### 作者:  
bowen https://www.github.com/WenTesla
### 最后修改时间
2023/2/4

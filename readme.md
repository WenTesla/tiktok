# 简易版抖音项目  

## 目前状态:    
**开发中**  
## 项目启动
直接启动
```shell
go run main.go
```

## 使用mvc分层结构
参考文献  （不懂的可以参考这篇文章）
https://juejin.cn/post/7152299022017888286
## 项目层次
### controller
### model
### service
### controller



**统一使用goland开发**，学生可以用学生邮箱去官网认证学生，免费使用

版本  
go版本 1.19

使用到的框架与依赖  
+ gin框架
+ gorm框架
+ mysql驱动
+ golang的jwt
+ 七牛云的oss存储  

已知错误    
可能和数据库的创建连接有关  
Error 1040: Too many connections (datasource)数据库连接过多  

错误日志（上传后点击才能响应接口）  
[GIN] 2023/02/01 - 17:26:57 | 400 |     11.1264ms |  192.168.31.236 | POST     "/douyin/publish/action/"

优化地方:
数据库访问过多，造成数据库压力大,后续使用redis优化  


作者:  
bowen 


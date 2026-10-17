# 简易版抖音项目

## 项目答辩文档与App地址
[飞书文档](https://xp8kgipb5a.feishu.cn/docx/RPAvdVcqpoc6DzxRPlnc5f3VnMb)
## 项目演示视频
[青训营演示视频](https://www.bilibili.com/video/BV1uT411S79U/?share_source=copy_web&vd_source=fe55b12bbf1a3c973a095834d9f2ba6d)

## 项目启动

### 方式一：本地直接启动（开发）
```shell
cd go
cp .env.example .env   # 按需填写 MySQL / Redis / JWT / COS 配置
go run main.go
```

### 方式二：容器化部署（推荐，Linux 一键启动，见 docs/DEPLOY.md）
```shell
# Linux 主机装好 Docker 即可（无需 Go，app 二进制在容器内多阶段编译）
bash scripts/up-linux.sh     # docker compose 拉起 MySQL + Redis + app，自动引导 .env / 准备 ffmpeg / 健康检查
bash scripts/down-linux.sh   # 关停清理（保留数据；--purge 连 MySQL 卷一起清）
# 裸命令等价：docker compose up -d --build
# Windows + WSL 且本机无 Docker 时的旧路径：bash scripts/up.sh（wslc）
```
- 启动自动 `AutoMigrate` 建表，无需手跑 `initial.sql`；健康检查 `/healthz`；app 等 MySQL/Redis 健康检查通过后启动。
- 视频/封面走**本地文件系统存储**（`./data` bind mount 持久化，不再依赖腾讯云 COS）；封面由镜像内静态 ffmpeg 自动抽首帧生成，ffmpeg 缺失时仅封面留空、发布播放不受影响。

## APP 操作
设置服务端地址
为方便测试登录和注册，及修改网络请求的服务器地址，提供了退出登录和高级设置两个能力。
1. 点击退出登录会自动重启
2. 在高级设置中可以配置自己的服务端项目的前缀地址，如下配置的http://192.168.1.7:8080
   在app中访问上述某个接口时就会拼接该前缀地址，例如访问 http://192.168.1.7:8080/douyin/feed/ 拉取视频列表
![](https://image-bed-1313520634.cos.ap-beijing.myqcloud.com/%E5%B1%8F%E5%B9%95%E6%88%AA%E5%9B%BE%202023-03-07%20162902.png)

## 表数据

## 使用mvc分层结构
参考文献  （不懂的可以参考这篇文章）
https://juejin.cn/post/7152299022017888286



## 版本  
* go 版本 1.26（已升级，原 1.19）
* mysql 8.4
* redis 7
* gin + gorm + golang-jwt/v4 + go-redis v9 + 腾讯云 COS

## 使用到的框架与依赖  
+ gin 框架（路由 / HTTP）
+ gorm 框架（ORM，MySQL 驱动）
+ golang-jwt/v4（HS256 鉴权，带 `exp` 过期）
+ bcrypt（密码哈希，替代原 MD5+盐）
+ 腾讯云 COS（视频与封面存储，工作流截取首帧 .jpg）
+ redis（限速计数器、点赞集合缓存）
+ viper + godotenv（配置外部化）
+ go-sqlmock / testify（单元测试）

## 工程化升级（相对原作业）
- 配置外部化（viper + 环境变量），密钥/连接不再写死
- JWT 带 `exp` 标准过期 + 可配置密钥，启动弱密钥告警
- 密码 bcrypt 哈希
- 上传按文件头魔数校验仅允许 mp4 且 ≤100MB
- feed / 列表接口消除 N+1（批量查询 + map 回填）
- MySQL 连接重试、AutoMigrate、/healthz、优雅关闭
- `model.db` 可注入，单测用 go-sqlmock，无需真实 MySQL

详见 `docs/ARCHITECTURE.md` / `docs/DEPLOY.md` / `docs/API.md`。

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

## to-do（原作业遗留，已部分完成）
* ~~使用 docker 部署~~ → 已完成（Docker Compose 多阶段编排，FROM scratch 运行镜像，Linux 一键启动，见 docs/DEPLOY.md）
* ~~自动执行 sql 语句~~ → 已完成（启动时 AutoMigrate）
* 读写分离 / 水平扩展（未做）
* 关注关系、消息分库分表（未做）

### 获奖证书
![Vfit13.jpeg](https://i.328888.xyz/2023/05/18/Vfit13.jpeg)

一个人完成，只拿到了结营证书，拿到优秀学员证书都是使用微服务

### 作者:  
bowen https://www.github.com/WenTesla
### 最后修改时间
2023/5/18

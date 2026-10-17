# 架构说明（ARCHITECTURE）

简易版抖音后端：Go + Gin + GORM + MySQL + Redis + 本地文件存储。本项目基于 2023 青训营作业做了工程化升级，目标是作为**简历作品集**展示可落地的后端工程能力。

## 1. 技术栈

| 层 | 选型 |
| --- | --- |
| Web 框架 | Gin |
| ORM | GORM（MySQL 驱动） |
| 缓存 / 限速 | Redis（go-redis v9） |
| 鉴权 | golang-jwt/v4（HS256，带 `exp`） |
| 对象存储 | 本地文件系统（`/data` bind mount；腾讯云 COS 已弃用，代码保留非默认） |
| 配置 | viper + godotenv（外部化，不再写死） |
| 部署 | Docker Compose 多阶段构建（`go/Dockerfile`：golang 编译 → `FROM scratch` 运行）+ wslc（Windows 旧路径） |

## 2. 目录结构（核心）

```
go/
├── main.go              # 启动：DB/Redis 探活、AutoMigrate、路由、优雅关闭
├── config/              # 配置加载（viper）、数据源、Redis、JWT、COS
├── route/               # 路由表（/douyin 分组）
├── controller/          # HTTP 入参校验、调用 service、组装响应
├── service/             # 业务逻辑（注册/登录/视频/点赞/关注/消息）
├── model/               # GORM 模型与数据访问（db 可注入，便于测试）
├── middle/jwt/          # 中间件：签发/校验 token、限速
├── util/                # 日志、敏感词、文件类型识别
resources/               # initial.sql、gin.log
scripts/                 # up-linux.sh / down-linux.sh（compose 编排）、up.sh / down.sh（wslc 旧路径）
```

## 3. 请求数据流

```
HTTP → route(gin) → controller(校验+组包) → service(业务) → model(DAO)
                                             ↓
                                    MySQL / Redis / 本地文件存储
```

- **鉴权**：`/douyin/user/register/`、`/douyin/user/login/` 免鉴权；其余写操作经 `jwt.VerifyToken` / `jwt.VerifyTokenByPost` 校验并写入 `c.Set("Id", ...)`。
- **限速**：未带 token 的 `/feed/` 与已登录用户分别按 IP / userId 走 Redis 计数器（阈值 60/min），超阈值直接 400。

## 4. 工程化升级点（相对原始作业）

### 4.1 配置外部化（Phase 1）
- `AppConfig` 经 viper 读取环境变量 / `.env`，覆盖 MySQL、Redis、JWT（密钥 + 有效期）、COS、Server 端口。
- 启动若使用代码内弱默认 JWT 密钥，会在日志给出生产环境告警。

### 4.2 安全（Phase 3）
- **JWT `exp`**：签发的 token 携带标准 `exp` 声明，由 golang-jwt 在解析时自动完成过期校验（移除了原本手写、且写死 30 天的过期判断）。
- **密码哈希**：注册/登录由 `MD5 + 静态盐` 升级为 **bcrypt**（`HashPassword` / `CheckPassword`），哈希逻辑统一收口在 service 层。
- **上传校验**：`/publish/action/` 上传前按**文件头魔数**校验仅允许 mp4，并限制单文件 ≤ 100MB（原实现几乎不校验，存在风险）。

### 4.3 性能：消除 N+1（Phase 4）
feed / 发布列表 / 点赞列表等接口，原本对每个视频分别查作者信息、点赞数、评论数、是否点赞（登录态），视频数 N 时产生 ~4N 次 DB 往返。

重构后（`service.videoService.packageVideos` + `model` 批量查询）：
1. 一次收集去重 `authorIds` 与 `videoIds`；
2. 批量取：作者信息（`UserServiceBatch`）、点赞数、评论数、是否点赞（仅登录态）；
3. 循环内仅做 O(N) 的 `map` 查找回填，无额外 DB/Redis 往返。

| 场景 | 重构前 | 重构后 |
| --- | --- | --- |
| feed 返回 N 条 | ~4N+ 次查询 | 固定 ~7 次（作者 1 + 关注 2 + 点赞 2 + 作品 1 + 评论 1） |

批量查询函数：`GetUserByIds`、`GetFollowCountsByUserIds`、`QueryIsFollowBatch`、`QueryLikeCountsByVideoIds`、`QueryIsLikeBatch`、`GetLikeCountsByUserIds`、`QueryCommentCountsByVideoIds`、`GetWorkCountsByUserIds`。

> 顺手修复的正确性 bug：原 `packageVideos` 循环内取 `&tableVideo`（循环变量地址），在旧 Go 下会导致所有元素共享同一地址、结果错位；改为索引取址 `&tableVideos[i]`。

### 4.4 健壮性（Phase 1/2）
- `model.db` 改为**可注入**（`InitDB`），启动时由 `main` 注入真实连接，单测注入 go-sqlmock，避免测试依赖真实 MySQL。
- MySQL 连接失败有 10×2s 重试，避免依赖未就绪时容器直接崩溃。
- 启动 `AutoMigrateAll` 替代手跑 `initial.sql`；暴露 `/healthz` 健康检查；支持 SIGINT/SIGTERM 优雅关闭。

## 5. 测试（Phase 5）
- **纯函数单测**：bcrypt 哈希/校验、邮箱/手机号格式、JWT 签发/解析/过期、文件头识别。
- **DAO/Service 单测**：用 `go-sqlmock` 注入 mock `*gorm.DB`，覆盖 N+1 批量查询（`UserServiceBatch` 等）。
- 依赖真实 MySQL/Redis 的集成测试通过 `TestMain` 在 `TIKTOK_INTEGRATION=1` 时才运行，默认 `go test ./...` 无需基础设施即全绿。

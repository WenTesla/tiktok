# 启动与部署指南（Startup & Deploy）

本项目是 **Go + Gin + GORM + MySQL + Redis** 的「仿抖音」后端。部署采用 **Docker Compose 容器化编排**：MySQL + Redis + app 三服务，在 **Linux 上一键启动**（Windows + WSL 的 wslc 旧路径降级为附录 A，仍可用）。

> 视频/封面走**本地文件系统存储**（不依赖腾讯云 COS，Phase 7 已替换）；封面由镜像内静态 ffmpeg 抽首帧生成。
> app 镜像为**多阶段构建**（`go/Dockerfile`）：容器内编译静态 Linux 二进制，**宿主机无需安装 Go**，只需 Docker。

---

## 1. 前置条件（Linux 主机）

- 已安装 **Docker**（含 `docker compose` v2 插件；或 `docker-compose` v1）。
- 网络可访问镜像仓库：`mysql:8.4`、`redis:7`、`golang:1.26-alpine`（仅构建期）。
  - 国内网络拉取 golang 镜像困难时：`docker compose build --build-arg GO_IMAGE=docker.m.daocloud.io/library/golang:1.26-alpine`。
- （可选）`openssl`：用于首次引导 `.env` 时生成随机 JWT 密钥（缺失则退路 `/dev/urandom`）。
- **不需要**在宿主机安装 Go / ffmpeg。

---

## 2. 一键启动

```bash
bash scripts/up-linux.sh
```

脚本依次完成：

1. **前置检查**：`docker` / `docker compose` / 守护进程可达，缺失给出明确错误；
2. **引导 `.env`**：首次运行由 `.env.example` 复制生成 `.env`，并把 `JWT_SECRET` 替换为 32 字节随机值（重启复用同一 `.env`，token 不失效）；已存在则跳过；
3. **准备 ffmpeg**：`ffmpeg-bin/ffmpeg-linux-amd64` 已是真实 ELF 则跳过，否则委托 `scripts/fetch-ffmpeg.sh` 拉取；失败仅告警不阻断（封面留空）；
4. **预建 `./data/{videos,covers,tmp}`**（compose bind mount 目标）；
5. **`docker compose up -d --build`**：多阶段构建 app 镜像 + 拉取 mysql/redis + 启动；app 会等 mysql/redis **健康检查通过**后才启动（`depends_on: condition: service_healthy`）；
6. **轮询 `/healthz`** 至多 90s，未就绪则以非零码退出并提示排查命令。

裸命令等价（`.env` 已存在时）：

```bash
docker compose up -d --build
```

### 2.1 配置（.env / 环境变量）

compose 变量写在仓库根 `.env`（首次自动生成，已被 `.gitignore`），也可用 shell 环境变量覆盖（优先级更高）：

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `MYSQL_PWD` | MySQL root 口令（mysql 服务与 app 共用，**务必修改**） | `zhang134679` |
| `JWT_SECRET` | JWT 签名密钥（`.env` 引导时随机生成） | `this_is_a_secret_key`（弱值，启动告警） |
| `PUBLIC_BASE_URL` | 静态资源对外基址（play_url/cover_url 前缀） | `http://localhost:8000/static` |

> 部署到服务器/虚拟机时，把 `PUBLIC_BASE_URL` 改成对外可达地址，例如 `http://192.168.1.7:8000/static`，否则客户端拿到的播放地址指向 localhost。

app 容器内的连接配置由 compose 固定注入（`MYSQL_HOST=mysql`、`REDIS_ADDR=redis:6379` 等），一般无需关心；应用配置体系（viper + 环境变量，本地裸跑用 `go/.env`）见 `docs/ARCHITECTURE.md` 4.1。

数据库表由 app 启动时 **AutoMigrate** 自动建表，无需手动执行 `resources/*.sql`。

---

## 3. 验证

```bash
curl http://localhost:8000/healthz          # {"status":"ok"}
curl http://localhost:8000/douyin/feed/     # 空库返回干净空列表
```

完整链路（注册 → 登录 → 发布 → feed 看到视频）：

```bash
# 1) 注册
curl -X POST http://localhost:8000/douyin/user/register/ \
  -F "username=demo" -F "password=12345678"
# 2) 登录取 token
TOKEN=$(curl -X POST http://localhost:8000/douyin/user/login/ \
  -F "username=demo" -F "password=12345678" | sed -E 's/.*"token":"([^"]+)".*/\1/')
# 3) 发布（token 走表单字段 token，不是 query）
curl -X POST http://localhost:8000/douyin/publish/action/ \
  -F "token=$TOKEN" -F "data=@sample.mp4" -F "title=hello"
# 4) feed 应返回 play_url（/static/videos/...）；配了 ffmpeg 时 cover_url 非空
curl http://localhost:8000/douyin/feed/
```

---

## 4. 关停 / 清理

```bash
bash scripts/down-linux.sh           # 停容器，保留 mysql-data 卷与 ./data
bash scripts/down-linux.sh --purge   # 连 mysql-data 卷一起删（./data 视频/封面仍保留）
```

- 视频/封面：`./data/{videos,covers,tmp}`（bind mount，资源管理器/`ls` 可直接查看、备份）。
- MySQL 数据：命名卷 `mysql-data`（`docker volume inspect tiktok_mysql-data`）。

---

## 5. 编排架构

```
                 ┌──────────────────────── docker network: douyin-net ───────────────────────┐
                 │                                                                          │
  host:8000 ───► │  douyin-app (FROM scratch, /src/app)                                      │
                 │    · 多阶段构建：golang:1.26-alpine 编译静态二进制 → scratch 运行           │
                 │    · 内嵌 /src/ffmpeg（抽帧）、resources/（gin.log）、go/util/key.txt       │
                 │    · 依赖 mysql/redis 健康检查通过后启动；SIGTERM 优雅关闭                  │
                 │        │                    │                                              │
                 │        ▼                    ▼                                              │
                 │  douyin-mysql ◄─────► douyin-redis                                         │
                 │  (mysql:8.4,          (redis:7,                                           │
                 │   healthcheck:         healthcheck:                                        │
                 │   mysqladmin ping)     redis-cli ping)                                     │
                 └──────────────────────────────────────────────────────────────────────────┘
  host:./data ──► /data（视频/封面）        volume mysql-data ──► /var/lib/mysql
```

- 镜像极简：运行阶段 `FROM scratch`，无 shell、无包管理器，攻击面与体积（~20MB）双低。
- 配置全部经环境变量注入（compose `environment:` ← `.env` 变量替换），无密钥入镜像。
- MySQL 8.4 root 口令即 `MYSQL_PWD`；app 与 mysql 服务同源，避免两处口令漂移。

---

## 6. 附录 A：Windows + WSL（wslc 旧路径）

本机 Docker 不可用时的替代方案：`wslc`（微软 WSL 容器 CLI，无 compose 子命令，编排写在脚本里）。

```bash
# 前置：wslc 可用；本机已拉取 docker.m.daocloud.io/library/{mysql:8.4,redis:7}
# 首次建议先准备 ffmpeg（失败不阻断）：bash scripts/fetch-ffmpeg.sh
export TIKTOK_DATA_DIR=/mnt/f/workspace/tiktok/data   # 可选，持久化视频到宿盘
bash scripts/up.sh        # 交叉编译 + FROM scratch 镜像（go/Dockerfile.prebuilt）+ 三容器
bash scripts/down.sh      # 关停清理（保留数据卷）
```

与 compose 路径的差异：二进制在 Windows 侧交叉编译（需 Go 1.26），镜像用 `go/Dockerfile.prebuilt`；环境变量含义一致（`TIKTOK_MYSQL_PWD` ↔ `MYSQL_PWD`）。新部署请优先使用第 2 节 compose 路径。

---

## 7. 附录 B：本地开发模式（不进容器）

需要 MySQL + Redis 可达（可直接用 compose 起的 mysql/redis，端口已映射 3306/6379 到宿主机）：

```bash
cd go
cp .env.example .env          # 按需改 JWT_SECRET、STORAGE 等
export LOCAL_STORAGE_DIR=./data
export PUBLIC_BASE_URL=http://localhost:8000/static
go run .
```

注意：本地裸跑**不抽封面**（ffmpeg 在镜像内 `/src/ffmpeg`，抽帧降级留空）；其余接口与容器部署完全一致。

---

## 8. 已知限制

- 单实例部署，未做读写分离 / 水平扩展（compose `restart: unless-stopped` 仅保进程存活）。
- 关注关系、消息等未做分库分表。
- 封面抽帧依赖镜像内 ffmpeg；未配置时仅封面留空，不影响视频发布与播放。
- arm64 主机：mysql/redis/app 均可原生构建，但 `ffmpeg-bin/ffmpeg-linux-amd64` 为 amd64 二进制，arm 上抽帧降级（可自行替换为 arm64 静态 ffmpeg）。

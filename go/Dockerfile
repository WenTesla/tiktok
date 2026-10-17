# 多阶段构建（Linux 容器编排默认路径，配合仓库根 docker-compose.yml）：
#   阶段一 builder：golang:1.26-alpine 内编译静态 Linux 二进制——宿主机只需 Docker，无需安装 Go；
#   阶段二 runtime：FROM scratch 极简运行镜像（无 shell/CA，体积 ~20MB），
#                   仅含二进制 + 运行资源（resources/、敏感词字典、静态 ffmpeg）。
#
# 构建上下文为仓库根目录：`docker compose up -d --build` 或 `docker build -f go/Dockerfile .`
# 国内网络拉取 golang 镜像困难时，可用镜像覆盖基础镜像：
#   docker compose build --build-arg GO_IMAGE=docker.m.daocloud.io/library/golang:1.26-alpine
#
# 关于 ffmpeg：最终阶段 COPY ffmpeg-bin/ffmpeg-linux-amd64 -> /src/ffmpeg。
#   该文件缺失或仍是占位符（非 ELF）时镜像仍可构建，但发布视频抽帧降级、封面留空
#   （视频发布/播放不受影响）；实机请先 bash scripts/fetch-ffmpeg.sh 拉取真实静态二进制。
#
# 旧路径（Windows 侧交叉编译预置二进制、wslc 部署）见 go/Dockerfile.prebuilt。

ARG GO_IMAGE=golang:1.26-alpine

FROM ${GO_IMAGE} AS builder
WORKDIR /src
# 先拷依赖清单利用层缓存（模块根在仓库根：go.mod / go.sum）
COPY go.mod go.sum ./
RUN go mod download
COPY go ./go
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/douyin-app ./go

FROM scratch
WORKDIR /src
COPY --from=builder /out/douyin-app /src/app
# 运行时资源：gin 日志目录、初始建表 SQL（resources/ 在仓库根目录）
COPY resources ./resources
# 敏感词字典（InitSensitiveFilter 以相对路径 go/util/key.txt 读取，依赖 CWD=/src）
COPY go/util/key.txt ./go/util/key.txt
# 静态 ffmpeg（服务端抽帧用）：构建上下文内的 ffmpeg-bin/ffmpeg-linux-amd64，需可执行权限
COPY ffmpeg-bin/ffmpeg-linux-amd64 /src/ffmpeg
EXPOSE 8000
ENTRYPOINT ["/src/app"]

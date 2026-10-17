#!/usr/bin/env bash
# 用 wslc（微软 WSL 容器 CLI）拉起简易抖音后端全套依赖。
# 用法：bash scripts/up.sh
#
# MySQL 口令不写死在脚本里：优先取环境变量 TIKTOK_MYSQL_PWD，
# 否则从 app 配置(go/config/load.go)中解析默认口令，保持与 app 一致。
set -euo pipefail

# 切换到项目根目录（本脚本位于 <root>/scripts/）
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
echo "==> 项目根目录: $ROOT"

# 解析 MySQL 口令
MYSQL_PWD="${TIKTOK_MYSQL_PWD:-}"
if [ -z "$MYSQL_PWD" ]; then
  MYSQL_PWD="$(grep -m1 'MYSQL_PASSWORD' go/config/load.go 2>/dev/null | sed -E 's/.*", *"([^"]+)".*/\1/')"
fi
if [ -z "$MYSQL_PWD" ]; then
  echo "未解析到口令，使用空口令模式"
  MYSQL_ENV="-e MYSQL_ALLOW_EMPTY_PASSWORD=yes"
  APP_PWD="-e MYSQL_PASSWORD="
else
  MYSQL_ENV="-e MYSQL_ROOT_PASSWORD=$MYSQL_PWD"
  APP_PWD="-e MYSQL_PASSWORD=$MYSQL_PWD"
fi

# 0. 清理已有容器/网络（幂等，便于重复执行）
for c in douyin-app douyin-redis douyin-mysql; do
  wslc stop "$c" >/dev/null 2>&1 || true
  wslc rm -f "$c" >/dev/null 2>&1 || true
done
wslc network rm douyin-net >/dev/null 2>&1 || true

# 0.1 确保静态 ffmpeg 二进制存在（用于服务端抽帧）。
# 仅当「真实 ELF 二进制」缺失时才尝试获取：优先委托 fetch-ffmpeg.sh（多镜像/续传），
# 仍失败（如沙箱网络受限）则不阻断部署——放置占位文件，运行时抽帧降级（封面留空）。
# 注意：占位文件同样 chmod +x，故必须用 is_elf 判定，不能只看 -x，否则占位态会永久绕过下载。
FFMPEG_BIN="$ROOT/ffmpeg-bin/ffmpeg-linux-amd64"
is_elf() {
  local f="$1"
  [ -f "$f" ] || return 1
  case "$(head -c4 "$f" 2>/dev/null | od -An -tx1 | tr -d ' \n')" in
    7f454c46) return 0 ;;   # \x7f E L F
  esac
  return 1
}
if is_elf "$FFMPEG_BIN"; then
  echo "==> ffmpeg 已就绪（真实二进制）：$FFMPEG_BIN"
else
  echo "==> 获取静态 ffmpeg (linux/amd64) 到 $FFMPEG_BIN"
  if [ -f "$ROOT/scripts/fetch-ffmpeg.sh" ]; then
    echo "    委托 scripts/fetch-ffmpeg.sh（多镜像/断点续传）..."
    bash "$ROOT/scripts/fetch-ffmpeg.sh" || true
  fi
  if ! is_elf "$FFMPEG_BIN"; then
    echo "警告: 仍无法获取 ffmpeg；放置占位文件以便镜像构建，运行时抽帧将降级（封面为空）。"
    echo "     实机请运行：bash scripts/fetch-ffmpeg.sh  再 bash scripts/up.sh"
    mkdir -p "$ROOT/ffmpeg-bin"
    printf '# placeholder: replace with a real static ffmpeg (linux/amd64) binary\n' > "$FFMPEG_BIN"
    chmod +x "$FFMPEG_BIN"
  fi
fi

# 0.2 本地存储卷（可选持久化）。
# 设置 TIKTOK_DATA_DIR 为 WSL 绝对路径（如 /mnt/f/workspace/tiktok/data）可启用宿主卷持久化；
# 未设置时数据写在容器内 /data（随容器销毁，验证/演示足够）。wslc 对 Windows 风格路径支持不稳定，
# 故默认不挂载宿主卷，避免部署失败。
DATA_VOL=""
if [ -n "${TIKTOK_DATA_DIR:-}" ]; then
  mkdir -p "$TIKTOK_DATA_DIR"/{videos,covers,tmp}
  DATA_VOL="-v $TIKTOK_DATA_DIR:/data"
  echo "==> 本地存储宿主卷: $TIKTOK_DATA_DIR -> /data"
else
  echo "==> 未设置 TIKTOK_DATA_DIR，数据存储在容器内 /data（非持久化）"
fi

# 1. 创建网络（幂等：已存在则忽略，避免二次运行时 "network already exists" 触发 set -e 中断）
echo "==> 1. 创建网络 douyin-net"
wslc network create douyin-net 2>/dev/null || true

# 2. 启动 MySQL（镜像已在本地：docker.m.daocloud.io/library/mysql:8.4）
echo "==> 2. 启动 MySQL"
wslc run -d --name douyin-mysql \
  --network douyin-net \
  $MYSQL_ENV \
  -e MYSQL_DATABASE=tiktok \
  -e MYSQL_CHARACTER_SET_SERVER=utf8mb4 \
  -e MYSQL_COLLATION_SERVER=utf8mb4_unicode_ci \
  -p 3306:3306 \
  docker.m.daocloud.io/library/mysql:8.4

echo "==> 等待 MySQL 就绪 (ready for connections)..."
for i in $(seq 1 40); do
  if timeout 8 wslc container logs douyin-mysql 2>/dev/null | grep -q "ready for connections"; then
    echo "MySQL 就绪"
    break
  fi
  sleep 3
done

# 3. 启动 Redis
echo "==> 3. 启动 Redis"
wslc run -d --name douyin-redis \
  --network douyin-net \
  -p 6379:6379 \
  docker.m.daocloud.io/library/redis:7

# 4. 交叉编译 Linux 静态二进制并构建最小镜像（FROM scratch，避免拉取 golang/alpine；
#    预置二进制路径用 go/Dockerfile.prebuilt；Linux compose 路径为 go/Dockerfile 多阶段构建）
echo "==> 4. 交叉编译 Linux 静态二进制"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o douyin-app-bin ./go
echo "==> 构建 app 镜像 (FROM scratch)"
wslc build -t douyin-app -f go/Dockerfile.prebuilt .
echo "==> 启动 app"
wslc run -d --name douyin-app \
  --network douyin-net \
  -p 8000:8000 \
  $DATA_VOL \
  -e MYSQL_HOST=douyin-mysql \
  -e MYSQL_PORT=3306 \
  -e MYSQL_USER=root \
  $APP_PWD \
  -e MYSQL_DB=tiktok \
  -e REDIS_ADDR=douyin-redis:6379 \
  -e SERVER_PORT=8000 \
  -e STORAGE_TYPE=local \
  -e LOCAL_STORAGE_DIR=/data \
  -e PUBLIC_BASE_URL=http://localhost:8000/static \
  douyin-app

# 5. 等待健康检查
echo "==> 5. 等待 app 健康检查 /healthz"
for i in $(seq 1 30); do
  if curl -fsS http://localhost:8000/healthz >/dev/null 2>&1; then
    echo "app 就绪 (healthz ok)"
    break
  fi
  sleep 2
done

echo "==> 完成。试试："
echo "    curl http://localhost:8000/douyin/feed/"
echo "    curl http://localhost:8000/healthz"

#!/usr/bin/env bash
# up-linux.sh —— Linux 容器化编排一键启动（docker compose 拉起 MySQL + Redis + app）。
#
# 前置：Linux 主机安装 Docker（compose v2 插件，或 docker-compose v1）。
#       Go 无需安装——app 二进制由多阶段构建在容器内编译（go/Dockerfile）。
#       网络可访问镜像仓库（mysql:8.4 / redis:7 / golang:1.26-alpine）；
#       国内网络可用 docker compose build --build-arg GO_IMAGE=<镜像> 覆盖 golang 基础镜像。
#
# 用法：bash scripts/up-linux.sh
#   首次运行自动由 .env.example 生成 .env（随机 JWT_SECRET），之后重启复用同一 .env。
#   也可预先 export MYSQL_PWD / JWT_SECRET / PUBLIC_BASE_URL 覆盖（优先级高于 .env）。
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

log()  { echo "==> $*"; }
warn() { echo "警告: $*" >&2; }
die()  { echo "错误: $*" >&2; exit 1; }

# ---------- 0) 前置检查：docker / compose / 守护进程 ----------
command -v docker >/dev/null 2>&1 \
  || die "未找到 docker，请先安装 Docker（含 compose 插件）。"
if docker compose version >/dev/null 2>&1; then
  DC="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  DC="docker-compose"
else
  die "未找到 docker compose / docker-compose，请安装 Docker Compose（v2 插件或 v1）。"
fi
docker info >/dev/null 2>&1 \
  || die "docker 守护进程不可达（docker info 失败）：请启动 Docker，或将当前用户加入 docker 组（或使用 sudo 运行）。"

# ---------- 1) 引导 .env（compose 变量文件，首次自动生成） ----------
if [ -f "$ROOT/.env" ]; then
  log ".env 已存在，跳过引导（删除 .env 可重新生成）"
else
  log "首次运行：由 .env.example 生成 .env ..."
  cp "$ROOT/.env.example" "$ROOT/.env"
  if command -v openssl >/dev/null 2>&1; then
    JWT="$(openssl rand -hex 32)"
  else
    JWT="$(head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n')"
  fi
  sed -i "s|^JWT_SECRET=.*|JWT_SECRET=$JWT|" "$ROOT/.env"
  log "已生成随机 JWT_SECRET 并写入 .env（重启不失效；请按需修改 MYSQL_PWD）"
fi

# ---------- 2) 准备静态 ffmpeg（linux/amd64，烤进镜像用于服务端抽帧） ----------
# 已是真实 ELF 则跳过；否则委托 fetch-ffmpeg.sh（多镜像/断点续传）。
# 获取失败不阻断部署：放占位文件，运行期抽帧降级（封面留空，视频发布/播放不受影响）。
is_elf() {
  local f="$1"
  [ -f "$f" ] || return 1
  case "$(head -c4 "$f" 2>/dev/null | od -An -tx1 | tr -d ' \n')" in
    7f454c46) return 0 ;;   # \x7f E L F
  esac
  return 1
}
FFMPEG_BIN="$ROOT/ffmpeg-bin/ffmpeg-linux-amd64"
if is_elf "$FFMPEG_BIN"; then
  log "ffmpeg 已就绪（真实二进制）"
else
  log "获取 ffmpeg（静态 linux/amd64）..."
  if bash "$ROOT/scripts/fetch-ffmpeg.sh"; then
    log "ffmpeg 就绪"
  else
    warn "ffmpeg 获取失败：封面将留空，视频发布/播放不受影响。"
    warn "网络恢复后执行 bash scripts/fetch-ffmpeg.sh 再 bash scripts/up-linux.sh 即可启用抽帧。"
    mkdir -p "$ROOT/ffmpeg-bin"
    printf '# placeholder: replace with a real static ffmpeg (linux/amd64) binary\n' > "$FFMPEG_BIN"
    chmod +x "$FFMPEG_BIN"
  fi
fi

# ---------- 3) 本地存储目录（compose bind mount ./data -> /data） ----------
mkdir -p "$ROOT/data/videos" "$ROOT/data/covers" "$ROOT/data/tmp"

# ---------- 4) 起栈：构建 app 镜像（容器内多阶段编译）+ 拉 mysql/redis + 启动 ----------
# depends_on 已配置 mysql/redis 健康检查，app 会等两者 healthy 后启动。
log "docker compose up -d --build ..."
$DC up -d --build

# ---------- 5) 等待 app 健康检查 ----------
log "等待 app /healthz ..."
ready=0
for _ in $(seq 1 45); do
  if curl -fsS http://localhost:8000/healthz >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 2
done
if [ "$ready" != "1" ]; then
  die "app 未在 90s 内就绪。排查：$DC ps；$DC logs app --tail 100"
fi

echo
echo "完成。试试："
echo "  curl http://localhost:8000/healthz"
echo "  curl http://localhost:8000/douyin/feed/"
echo "视频/封面在 ./data（bind mount，可备份）；MySQL 数据在卷 mysql-data。"
echo "关停：bash scripts/down-linux.sh（保留数据）/ bash scripts/down-linux.sh --purge（连 MySQL 数据一起清）"

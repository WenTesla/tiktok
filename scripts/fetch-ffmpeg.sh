#!/usr/bin/env bash
# fetch-ffmpeg.sh —— 实机拉取「静态 ffmpeg (linux/amd64)」二进制到 ffmpeg-bin/ffmpeg-linux-amd64
#
# 背景：
#   - 沙箱网络限流（~25KB/s + TLS 证书吊销），up.sh 内置的 120s 下载会失败并放置占位文件，
#     导致 /publish/action/ 抽帧降级、封面留空。
#   - 本脚本专供网络正常的实机运行：多镜像、断点续传、校验 ELF 魔数、必要时重建镜像。
#
# 用法：
#   bash scripts/fetch-ffmpeg.sh            # 仅拉取并校验二进制
#   REBUILD=1 bash scripts/fetch-ffmpeg.sh  # 拉取后自动 bash scripts/up.sh 重建并重启容器
#
# 也可手动指定镜像（二选一）：
#   FFMPEG_URL=https://.../ffmpeg.tar.xz  bash scripts/fetch-ffmpeg.sh   # tar.xz，内含 ffmpeg 可执行文件
#   FFMPEG_URL=https://.../ffmpeg-linux    bash scripts/fetch-ffmpeg.sh   # 单文件二进制（直拷）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
TARGET="$ROOT/ffmpeg-bin/ffmpeg-linux-amd64"
# 强制 POSIX 临时目录：避免 Git Bash 下 mktemp 返回 C:\... 路径，导致 tar/find/rm 异常
export TMPDIR=/tmp
TMPD="$(mktemp -d)"
trap 'rm -rf "$TMPD"' EXIT

log() { echo "==> $*"; }
die() { echo "错误: $*" >&2; exit 1; }

# 校验是否为真实 ELF 可执行文件（占位文件是文本，非 ELF）
is_elf() {
  local f="$1"
  [ -f "$f" ] || return 1
  case "$(head -c4 "$f" 2>/dev/null | od -An -tx1 | tr -d ' \n')" in
    7f454c46) return 0 ;;   # \x7f E L F
  esac
  return 1
}

# 若已存在真实二进制，直接收工（可被 FORCE=1 强制覆盖）
if is_elf "$TARGET" && [ "${FORCE:-0}" != "1" ]; then
  log "已存在真实 ffmpeg 二进制：$TARGET（$(du -h "$TARGET" | cut -f1)），跳过。"
  exit 0
fi

mkdir -p "$ROOT/ffmpeg-bin"

# 优先：本机已装 ffmpeg（WSL/Linux 常见），直接复用，零下载
if command -v ffmpeg >/dev/null 2>&1; then
  # 仅接受 linux/amd64 静态产物（file 输出含 "statically linked" 最佳，但动态链接也可先用）
  if cp "$(command -v ffmpeg)" "$TARGET" 2>/dev/null && chmod +x "$TARGET" && is_elf "$TARGET"; then
    log "复用本机 ffmpeg: $(command -v ffmpeg)"
    log "已写入 $TARGET"
    [ "${REBUILD:-0}" = "1" ] && bash "$ROOT/scripts/up.sh"
    exit 0
  fi
fi

# 走到这里说明需要联网下载：确认 curl 可用（缺失则明确报错，非致命，up.sh 会降级占位）
if ! command -v curl >/dev/null 2>&1; then
  die "未找到 curl，无法自动下载 ffmpeg。请手动放置真实静态二进制到 $TARGET 后重跑，或先安装 curl。"
fi

# 下载候选：tar.xz（内含 ffmpeg 可执行文件）与单文件二进制两路
# 顺序：BtbN(最新 master, GitHub) -> johnvansickle(稳定 release) -> 用户覆盖
MIRRORS_TAR=(
  "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-linux64-gpl.tar.xz"
  "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
)
MIRRORS_SINGLE=(
  "https://github.com/eugeneware/ffmpeg-static/releases/latest/download/ffmpeg-linux-x64"
)

# 快速失败策略：受限/防火墙网络下不会长时间卡死（约 1~2 分钟即放弃并降级），
# --ssl-no-revoke 规避 Windows curl 偶发的证书吊销误报。
CURL="curl -fsSL --ssl-no-revoke --connect-timeout 12 --max-time 90 --retry 1 --retry-delay 2 -C -"

download_one() {
  # $1 = url, $2 = 输出路径；成功返回 0（外层 timeout 作为硬上限兜底）
  timeout 120 $CURL "$1" -o "$2" 2>/dev/null
}

# 1) 用户显式指定 FFMPEG_URL
if [ -n "${FFMPEG_URL:-}" ]; then
  log "使用指定地址：$FFMPEG_URL"
  if [[ "$FFMPEG_URL" == *.tar.* ]]; then
    ARCHIVE="$TMPD/ff.txz"; download_one "$FFMPEG_URL" "$ARCHIVE" || die "下载失败：$FFMPEG_URL"
    extract_inner "$ARCHIVE"
  else
    download_one "$FFMPEG_URL" "$TARGET" && chmod +x "$TARGET" || die "下载失败：$FFMPEG_URL"
  fi
else
  # 2) 依次尝试 tar.xz 镜像
  GOT=0
  for u in "${MIRRORS_TAR[@]}"; do
    log "尝试 tar.xz 镜像：$u"
    ARCHIVE="$TMPD/ff.txz"
    if download_one "$u" "$ARCHIVE" && [ -s "$ARCHIVE" ]; then
      if extract_inner "$ARCHIVE"; then GOT=1; break; fi
    fi
  done
  # 3) 退而求其次：单文件二进制镜像
  if [ "$GOT" = "0" ]; then
    for u in "${MIRRORS_SINGLE[@]}"; do
      log "尝试单文件镜像：$u"
      if download_one "$u" "$TARGET" && chmod +x "$TARGET" && is_elf "$TARGET"; then GOT=1; break; fi
    done
  fi
  [ "$GOT" = "1" ] || die "所有镜像下载失败，请手动放置真实静态二进制到 $TARGET 后重跑。"
fi

# 从 tar.xz 内提取名为 ffmpeg 的可执行文件并落盘到 TARGET
extract_inner() {
  local archive="$1"
  log "解包并提取 ffmpeg 二进制 ..."
  tar -xf "$archive" -C "$TMPD" 2>/dev/null || return 1
  # 递归找第一个名为 ffmpeg 的 ELF 可执行文件
  local found=""
  while IFS= read -r cand; do
    if is_elf "$cand" && [ -x "$cand" ]; then found="$cand"; break; fi
  done < <(find "$TMPD" -type f -name ffmpeg 2>/dev/null)
  [ -n "$found" ] || return 1
  cp "$found" "$TARGET" && chmod +x "$TARGET" && is_elf "$TARGET"
}

# 最终校验
is_elf "$TARGET" || die "$TARGET 不是有效的 ELF 可执行文件，请检查下载内容。"
log "ffmpeg 就绪：$TARGET（$(du -h "$TARGET" | cut -f1)）"
log "版本：$("$TARGET" -version 2>/dev/null | head -n1 || echo '（Windows 侧无法执行 linux 二进制，进容器后由抽帧逻辑验证）')"

if [ "${REBUILD:-0}" = "1" ]; then
  log "重建并重启容器 ..."
  bash "$ROOT/scripts/up.sh"
fi
log "完成。发布视频后封面将由服务端 ffmpeg 自动抽取（不再留空）。"

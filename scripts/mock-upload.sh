#!/usr/bin/env bash
# mock-upload.sh —— 模拟一次完整视频上传链路，用于本地联调。
# 自动造测试视频（优先 ffmpeg 生成真视频，否则造最小 ftyp 桩文件）→ 注册/登录取 token
# → 上传到 /douyin/publish/action/ → 拉 /douyin/feed/ 验证。
#
# 用法:
#   bash scripts/mock-upload.sh                                  # 默认 http://localhost:8000，ffmpeg 生成 mp4
#   BASE_URL=http://1.2.3.4:8000 bash scripts/mock-upload.sh     # 指定服务地址
#   VIDEO_FILE=/path/to/a.mov bash scripts/mock-upload.sh         # 直接上传指定文件（测其它格式）
#   FMT=webm bash scripts/mock-upload.sh                          # 用 ffmpeg 生成指定格式（需对应编码器）

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8000}"
FMT="${FMT:-mp4}"
VIDEO_FILE="${VIDEO_FILE:-}"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

log() { echo "==> $*"; }
die() { echo "错误: $*" >&2; exit 1; }

# 1) 准备测试视频
if [ -n "$VIDEO_FILE" ]; then
  [ -f "$VIDEO_FILE" ] || die "VIDEO_FILE 不存在: $VIDEO_FILE"
  VIDEO="$VIDEO_FILE"
  log "使用指定视频: $VIDEO"
else
  VIDEO="$WORK/sample.$FMT"
  if command -v ffmpeg >/dev/null 2>&1; then
    log "用 ffmpeg 生成测试视频 ($FMT)..."
    case "$FMT" in
      mp4)  ffmpeg -y -f lavfi -i testsrc=size=320x240:rate=15 -t 2 -c:v libx264 -pix_fmt yuv420p "$VIDEO" >/dev/null 2>&1 \
              || ffmpeg -y -f lavfi -i testsrc=size=320x240:rate=15 -t 2 "$VIDEO" >/dev/null 2>&1 ;;
      webm) ffmpeg -y -f lavfi -i testsrc=size=320x240:rate=15 -t 2 -c:v libvpx-vp9 "$VIDEO" >/dev/null 2>&1 \
              || ffmpeg -y -f lavfi -i testsrc=size=320x240:rate=15 -t 2 -c:v libvpx "$VIDEO" >/dev/null 2>&1 ;;
      mov)  ffmpeg -y -f lavfi -i testsrc=size=320x240:rate=15 -t 2 -c:v libx264 -pix_fmt yuv420p -f mov "$VIDEO" >/dev/null 2>&1 ;;
      flv)  ffmpeg -y -f lavfi -i testsrc=size=320x240:rate=15 -t 2 -c:v flv1 "$VIDEO" >/dev/null 2>&1 ;;
      avi)  ffmpeg -y -f lavfi -i testsrc=size=320x240:rate=15 -t 2 -c:v mpeg4 "$VIDEO" >/dev/null 2>&1 ;;
      mpg)  ffmpeg -y -f lavfi -i testsrc=size=320x240:rate=15 -t 2 -c:v mpeg2video "$VIDEO" >/dev/null 2>&1 ;;
      *)    die "未内置 $FMT 的 ffmpeg 生成命令，请用 VIDEO_FILE 指定文件" ;;
    esac
    [ -s "$VIDEO" ] || die "ffmpeg 生成 $FMT 失败"
  else
    log "未找到 ffmpeg，生成最小 ftyp 桩文件（仅用于越过格式校验，非可播放视频）"
    # 最小 MP4：ftyp box（size=0x18, 'ftyp','isom'），足以通过 IsVideo 校验
    printf '\x00\x00\x00\x18\x66\x74\x79\x70\x69\x73\x6f\x6d\x00\x00\x00\x00\x00\x00\x00\x00' > "$VIDEO"
    [ -s "$VIDEO" ] || die "无法写入桩文件"
  fi
fi
log "视频大小: $(wc -c < "$VIDEO") 字节"

# 2) 注册并取 token
USER="mock_$$_$RANDOM"
PASS="mock1234"
log "注册用户: $USER"
REG="$(curl -fsS "$BASE_URL/douyin/user/register/?username=$USER&password=$PASS")" \
  || die "注册请求失败（服务可达？BASE_URL=$BASE_URL）"
TOKEN="$(printf '%s' "$REG" | grep -o '"token":"[^"]*"' | head -1 | sed 's/"token":"//;s/"$//')"
[ -n "$TOKEN" ] || die "注册未返回 token，响应: $REG"

# 3) 上传
TITLE="mock-upload-$(date +%s)"
log "上传视频 (title=$TITLE)..."
PUB="$(curl -fsS -X POST "$BASE_URL/douyin/publish/action/" \
  -F "token=$TOKEN" -F "title=$TITLE" -F "data=@$VIDEO")" \
  || die "上传请求失败"
echo "publish 响应: $PUB"

# 4) 拉 feed 验证
log "拉取 feed 验证..."
FEED="$(curl -fsS "$BASE_URL/douyin/feed/")" || die "feed 请求失败"
echo "feed 响应: $FEED"
echo "完成 ✅"

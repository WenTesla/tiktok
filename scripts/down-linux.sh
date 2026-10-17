#!/usr/bin/env bash
# down-linux.sh —— 关停 Linux 容器编排（对应 up-linux.sh）。
#
# 默认仅停止并移除容器；mysql 数据卷（mysql-data）与 ./data（视频/封面）保留。
#   bash scripts/down-linux.sh           # 关停，保留数据
#   bash scripts/down-linux.sh --purge   # 连 mysql-data 卷一起删除（./data 视频/封面仍保留）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if docker compose version >/dev/null 2>&1; then DC="docker compose";
elif command -v docker-compose >/dev/null 2>&1; then DC="docker-compose";
else echo "错误: 未找到 docker compose / docker-compose" >&2; exit 1; fi

if [ "${1:-}" = "--purge" ]; then
  echo "==> 停止并移除容器 + 删除 mysql-data 卷（./data 视频/封面仍保留）"
  $DC down -v
else
  echo "==> 停止并移除容器（保留卷）"
  $DC down
  echo "已关停。视频/封面在 ./data，MySQL 数据在卷 mysql-data，均未删除。"
  echo "如需连 MySQL 数据一起清：bash scripts/down-linux.sh --purge"
fi

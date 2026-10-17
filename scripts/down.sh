#!/usr/bin/env bash
# 停止并移除由 scripts/up.sh 创建的容器（保留数据卷）。
set -euo pipefail

echo "==> 停止容器"
wslc stop douyin-app douyin-redis douyin-mysql 2>/dev/null || true

echo "==> 移除容器"
wslc rm douyin-app douyin-redis douyin-mysql 2>/dev/null || true

# 如需彻底清理 MySQL 数据卷，取消下一行注释：
# wslc volume rm douyin-mysql-data 2>/dev/null || true

echo "已停止并移除容器（数据卷 douyin-mysql-data 保留）。"

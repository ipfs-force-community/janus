#!/bin/bash

# 给执行遇到错误时自动退出
set -e

echo "🚀 开始 Janus 项目全流程自动部署验证"
echo "================================================="

echo -e "\n[1] 环境验证与同步"
if [ ! -f .env ]; then
  echo "👉 正在从 .env.docker 自动生成 .env 环境变量文件..."
  cp .env.docker .env
else
  echo "✔ 检测到 .env 文件已存在"
fi

echo -e "\n[2] 启动底层依赖及外挂件 (MySQL & Mock Node)"
# 确保数据库先启动到后台，同时我们把刚加的 filecoin 测试伪造节点也拉起来
docker compose up -d mysql filecoin-mock

echo "⏳ 等待 MySQL 完全就绪接收连接（约 15 秒）..."
# Docker 有内建的健康检查配合 depends_on，但这仅在组合被拉起时有效。
# 此处简单做个前置等待，确保存储底层健康。
sleep 15

echo -e "\n[3] 开始执行一次性矿工数据同步任务 (janus-miner)"
echo "-------------------------------------------------"
# 由于配置了 profile，我们直接执行它，它前台执行完会自动退出
docker compose --profile init-task up janus-miner
echo "-------------------------------------------------"
echo "✔ 初始化数据任务执行完毕。"


echo -e "\n[4] 一键拉起余下所有常驻全核心服务 (API, Indexer, Frontend, Nginx)"
docker compose up -d

echo -e "\n🎉 部署流程演示搞定！请看当前所有存活的系统服务："
docker compose ps

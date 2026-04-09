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

echo -e "\n[2] 启动 MySQL"
docker compose up -d mysql

echo "⏳ 等待 MySQL 完全就绪接收连接（约 15 秒）..."
sleep 15

echo -e "\n[3] 开始执行一次性矿工数据同步任务 (janus-miner)"
echo "-------------------------------------------------"
docker compose --profile init-task up janus-miner
echo "-------------------------------------------------"
echo "✔ 初始化数据任务执行完毕。"

echo -e "\n[4] 一键拉起余下所有常驻服务 (API, Indexer, Frontend)"
docker compose up -d mysql api indexer frontend

echo -e "\n🎉 部署流程演示搞定！请看当前所有存活的系统服务："
docker compose ps

echo -e "\n📝 服务访问地址："
echo "   - Frontend: http://localhost:50001"
echo "   - API: http://localhost:50002"
echo "   - MySQL: localhost:50003"

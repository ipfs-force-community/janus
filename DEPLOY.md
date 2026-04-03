# Janus 部署指南

## 启动流程

```
┌─────────────────────────────────────────────────────────────┐
│  第一步：一次性数据同步                                       │
│  ─────────────────────                                     │
│  1. 启动 MySQL                                              │
│  2. 启动 janus-miner（等待完成）                              │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  第二步：启动完整服务                                         │
│  ─────────────────────                                     │
│  1. MySQL（已启动）                                         │
│  2. indexer、api（并行）                                     │
│  3. frontend                                                │
│  4. nginx                                                   │
└─────────────────────────────────────────────────────────────┘
```

## 步骤详解

### 第一步：配置环境变量

```bash
# 1. 复制环境变量文件
cp .env.docker .env

# 2. 编辑 .env，配置以下必填项：
#    - MYSQL_ROOT_PASSWORD
#    - FILECOIN_NODE_ENDPOINT
#    - FILECOIN_NODE_TOKEN
#    - MINER_START_EPOCH（可选）
#    - MINER_END_EPOCH（可选）

# 3. ⚠️ 极度重要：同步后端密码映射
# 请务必打开 backend/config/config.docker.yaml 文件
# 并将其中的 "password" 与 "db_name" 改得与你的 .env 内部值一模一样，否则 API 和数据同步程序会报错 Access denied。
```

### 第二步：一次性数据同步

```bash
# 1. 独立启动 MySQL（后台运行），作为数据同步的底层依赖
docker compose up -d mysql

# 2. 指定 init-task 配置组来启动 miner 服务
docker compose --profile init-task up janus-miner

# 终端输出 Exited (0) 则说明完成。
# 若需要放入后台执行，可加上 -d 命令：
# docker compose --profile init-task up -d janus-miner
```

### 第三步：启动完整服务

```bash
# 启动所有常驻服务
# 由于之前的 init-task 隔离，这个默认启动指令会自动忽略一次性同步任务
docker compose up -d
```

## 服务依赖关系

```
mysql (健康检查)
  ├─→ api (depends_on mysql health)
  ├─→ indexer (depends_on mysql health)
  └─→ janus-miner (depends_on mysql health, restart: no)

api ──→ frontend ──→ nginx
```

## 常用命令

```bash
# 查看所有服务状态
docker compose ps

# 查看日志
docker compose logs -f janus-miner  # 查看同步日志
docker compose logs -f api          # 查看 API 日志
docker compose logs -f indexer      # 查看 Indexer 日志

# 停止所有服务
docker compose down

# 重新运行一次性同步（如果需要重新同步）
docker compose --profile init-task up janus-miner
```

## 部署模式

### 生产模式（只暴露 nginx）

```bash
# 不使用 override 文件，只启动完整服务
docker compose -f docker-compose.yml up -d

# 外部访问：http://<服务器IP>:80
```

### 开发模式（暴露所有端口）

```bash
# 默认使用 override 文件，启动全部映射服务
docker compose up -d

# 可通过以下地址访问：
# - Frontend: http://localhost:3001
# - API: http://localhost:10086
# - MySQL: localhost:3306
```

## 环境变量说明

| 变量 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `MYSQL_ROOT_PASSWORD` | 是 | MySQL root 密码 | - |
| `MYSQL_DATABASE` | 否 | 数据库名 | janus |
| `FILECOIN_NODE_ENDPOINT` | 是 | Filecoin 节点 RPC 地址 | - |
| `FILECOIN_NODE_TOKEN` | 是 | Filecoin 节点 Token | - |
| `INDEXER_INTERVAL` | 否 | Indexer 轮询间隔（秒） | 10 |
| `MINER_START_EPOCH` | 否 | 同步起始 epoch | 5260000 |
| `MINER_END_EPOCH` | 否 | 同步结束 epoch | 5261000 |
| `NGINX_PORT` | 否 | 对外暴露端口 | 80 |

## 服务说明

| 服务 | 类型 | 说明 |
|------|------|------|
| nginx | 常驻 | 反向代理，对外提供 HTTP 服务 |
| frontend | 常驻 | Next.js 前端 |
| api | 常驻 | API 服务，提供接口给前端 |
| indexer | 常驻 | 索引服务，持续同步链上数据 |
| janus-miner | 一次性 | 初始化数据同步任务 |
| mysql | 常驻 | MySQL 数据库 |

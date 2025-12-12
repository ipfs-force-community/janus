# Root Makefile

# Directories
BACKEND_DIR := backend
FRONTEND_DIR := frontend
BIN_DIR := $(BACKEND_DIR)/bin

# Backend binaries
API_BIN := $(BIN_DIR)/api
JANUS_BIN := $(BIN_DIR)/janus
INDEXER_BIN := $(BIN_DIR)/indexer

# 检测操作系统
OS := $(shell uname -s)

# 根据系统选择打包命令
ifeq ($(OS),Darwin)
	TAR_CMD = COPYFILE_DISABLE=1 gtar -czf
else
	TAR_CMD = tar -czf
endif

.PHONY: all backend frontend clean

# Default: build both backend and frontend
all: backend frontend

# -------------------
# Backend build rules
# -------------------
backend: $(API_BIN) $(JANUS_BIN) $(INDEXER_BIN)

$(API_BIN): $(BACKEND_DIR)/cmd/api/main.go
	@echo "Building backend API..."
	@mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o bin/api ./cmd/api/main.go

$(JANUS_BIN): $(wildcard $(BACKEND_DIR)/cmd/janus/*.go)
	@echo "Building backend Janus..."
	@mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o bin/janus ./cmd/janus/*.go

$(INDEXER_BIN): $(wildcard $(BACKEND_DIR)/cmd/indexer/*.go)
	@echo "Building backend Indexer..."
	@mkdir -p $(BIN_DIR)
	cd $(BACKEND_DIR) && go build -o bin/indexer ./cmd/indexer/*.go
# -------------------
# Frontend build rules
# -------------------
frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build

# -------------------
# PM2 Deployment Rules
# -------------------
.PHONY: frontend-package pm2-start pm2-stop pm2-restart

# Create a compressed package of the frontend build for deployment
frontend-package:
	@echo "Creating frontend deployment package..."
	cd $(FRONTEND_DIR) && npm install && npm run build && npm ci --omit=dev
	@echo "Compressing build artifacts..."
	cd $(FRONTEND_DIR) && $(TAR_CMD) \
	  janus-frontend-$(shell date +%Y%m%d%H%M%S).tar.gz \
	  --exclude=".next/cache" \
	  --exclude="*.ts" \
	  --exclude="*.DS_Store" \
	  --exclude="*.log" \
	  .next \
	  public \
	  node_modules \
	  package.json \
	  next.config.* \
	  ecosystem.config.js 

# Start PM2 service (for production)
pm2-start:
	@echo "Starting PM2 service..."
	cd $(FRONTEND_DIR) && pm2 start ecosystem.config.js --env production

# Stop PM2 service
pm2-stop:
	@echo "Stopping PM2 service..."
	cd $(FRONTEND_DIR) && pm2 stop ecosystem.config.js

# Restart PM2 service (for hot reload)
pm2-restart:
	@echo "Restarting PM2 service..."
	cd $(FRONTEND_DIR) && pm2 restart ecosystem.config.js

# -------------------
# Cleanup
# -------------------
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	cd $(FRONTEND_DIR) && rm -rf .next
	@echo "Cleaning deployment packages..."
	cd $(FRONTEND_DIR) && rm -f janus-frontend-*.tar.gz

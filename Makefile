# Root Makefile

# Directories
BACKEND_DIR := backend
FRONTEND_DIR := frontend
BIN_DIR := $(BACKEND_DIR)/bin

# Backend binaries
API_BIN := $(BIN_DIR)/api
JANUS_BIN := $(BIN_DIR)/janus
INDEXER_BIN := $(BIN_DIR)/indexer

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
# Cleanup
# -------------------
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	cd $(FRONTEND_DIR) && rm -rf .next

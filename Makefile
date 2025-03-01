# Makefile

GO=go
BINARY=airules
BUILD_DIR=bin
TEMPLATE_DIR=templates

# デフォルトターゲット
.PHONY: all
all: clean build

# ビルド
.PHONY: build
build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY) main.go

# クリーン
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# テスト
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# テストカバレッジ
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# エンドツーエンドテスト
.PHONY: test-e2e
test-e2e:
	@echo "Running end-to-end tests..."
	$(GO) test -v ./e2e

# インストール
.PHONY: install
install: build
	@echo "Installing $(BINARY) to $(GOPATH)/bin"
	@cp $(BUILD_DIR)/$(BINARY) $(GOPATH)/bin/

# 依存関係の更新
.PHONY: deps
deps:
	@echo "Updating dependencies..."
	$(GO) mod tidy

# ヘルプ
.PHONY: help
help:
	@echo "使用可能なターゲット:"
	@echo "  all            - クリーンビルド"
	@echo "  build          - バイナリをビルド"
	@echo "  clean          - ビルドディレクトリを削除"
	@echo "  test           - すべてのテストを実行"
	@echo "  test-coverage  - カバレッジ付きでテストを実行"
	@echo "  test-e2e       - エンドツーエンドテストを実行"
	@echo "  install        - $(GOPATH)/bin にバイナリをインストール"
	@echo "  deps           - 依存関係を更新"
	@echo "  help           - このヘルプを表示"

.PHONY: help build run test clean docker-up docker-down migrate seed

help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## アプリケーションをビルド
	go build -o bin/main ./cmd

run: ## アプリケーションをローカルで実行
	go run ./cmd

dev: ## ホットリロードで開発サーバーを起動
	air

test: ## テストを実行
	go test ./...

test-cover: ## カバレッジ付きテストを実行
	go test -cover ./...

clean: ## ビルドファイルを削除
	rm -rf bin/ tmp/

docker-up: ## Dockerでサービスを起動
	docker-compose up -d

docker-down: ## Dockerサービスを停止
	docker-compose down

docker-logs: ## Dockerのログを表示
	docker-compose logs -f

deps: ## 依存関係をインストール
	go mod tidy

fmt: ## コードをフォーマット
	go fmt ./...

lint: ## コードの静的解析
	golangci-lint run

migrate: ## データベースマイグレーションを実行
	@echo "マイグレーション機能は今後実装予定"

seed: ## シードデータを投入
	@echo "シード機能は今後実装予定"

.DEFAULT_GOAL := help
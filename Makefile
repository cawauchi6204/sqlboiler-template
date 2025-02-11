.PHONY: run build test clean docker-up docker-down docker-build docker-reset db-up db-down db-reset migrate generate-models seed install help

# デフォルトのターゲット
.DEFAULT_GOAL := help

# 変数定義
BINARY_NAME=twitter-clone
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

# データベース接続情報
DB_HOST=localhost
DB_PORT=3307
DB_USER=root
DB_PASSWORD=example
DB_NAME=todoapp

# ヘルプメッセージ
help:
	@echo "利用可能なコマンド:"
	@echo "  make run              - アプリケーションをローカルで実行"
	@echo "  make build            - アプリケーションをビルド"
	@echo "  make test             - テストを実行"
	@echo "  make clean            - ビルドファイルを削除"
	@echo "  make docker-up        - 全てのDockerコンテナを起動"
	@echo "  make docker-down      - 全てのDockerコンテナを停止"
	@echo "  make docker-build     - Dockerイメージを再ビルド"
	@echo "  make docker-reset     - Dockerコンテナをリセット(停止→削除→ビルド→起動)"
	@echo "  make db-up            - データベースコンテナのみを起動"
	@echo "  make db-down          - データベースコンテナのみを停止"
	@echo "  make db-reset         - データベースをリセット(停止→削除→起動)"
	@echo "  make migrate          - データベースマイグレーションを実行"
	@echo "  make generate-models  - SQLBoilerでモデルを生成"
	@echo "  make seed             - テストデータを生成"
	@echo "  make install          - 依存関係をインストール"

# アプリケーションの実行(ローカル)
run:
	go run main.go

# アプリケーションのビルド
build:
	go build -o $(BINARY_NAME) main.go

# テストの実行
test:
	go test -v ./...

# ビルドファイルのクリーンアップ
clean:
	go clean
	rm -f $(BINARY_NAME)

# Dockerコンテナの起動(全て)
docker-up:
	docker-compose up -d
	@echo "全てのコンテナの起動を待機中..."
	@sleep 10

# Dockerコンテナの停止(全て)
docker-down:
	docker-compose down

# Dockerイメージの再ビルド
docker-build:
	docker-compose build

# Dockerコンテナのリセット
docker-reset: docker-down
	docker-compose rm -f
	docker-compose build
	docker-compose up -d
	@echo "コンテナの起動を待機中..."
	@sleep 10
	make migrate

# データベースコンテナの起動
db-up:
	docker-compose up -d db
	@echo "データベースの起動を待機中..."
	@sleep 10

# データベースコンテナの停止
db-down:
	docker-compose stop db

# データベースのリセット
db-reset: db-down
	docker-compose rm -f db
	docker-compose up -d db
	@echo "データベースの起動を待機中..."
	@sleep 10
	make migrate

# マイグレーションの実行
migrate:
	./run_migrations.sh

# SQLBoilerでモデルの生成
generate-models:
	sqlboiler mysql

# シードデータの生成
seed:
	DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) go run cmd/seed/main.go

# 依存関係のインストール
install:
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
	go mod download
	go mod tidy

# 開発用の便利なコマンド
dev: docker-up migrate generate-models

# データベースの再構築とアプリケーションの起動
reset: docker-reset generate-models

# 開発環境の完全セットアップ(データ含む)
setup: docker-reset generate-models seed

# テスト環境のセットアップと実行
test-all: docker-up migrate generate-models test

# コードの静的解析
lint:
	go vet ./...
	golangci-lint run

# モジュールの更新
update:
	go get -u ./...
	go mod tidy

# 本番環境用のビルドと起動
prod: docker-build docker-up
.PHONY: run build test clean docker-up docker-down docker-build docker-reset db-up db-down db-reset migrate generate-models seed install help api-register api-login api-tweet api-profile api-follow api-like

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
	@echo ""
	@echo "開発用コマンド:"
	@echo "  make run              - アプリケーションをローカルで実行"
	@echo "  make dev              - 開発環境を起動"
	@echo "  make build            - アプリケーションをビルド"
	@echo "  make test             - テストを実行"
	@echo ""
	@echo "Docker関連:"
	@echo "  make docker-up        - 全てのDockerコンテナを起動"
	@echo "  make docker-down      - 全てのDockerコンテナを停止"
	@echo "  make docker-build     - Dockerイメージを再ビルド"
	@echo "  make docker-reset     - Dockerコンテナをリセット"
	@echo ""
	@echo "データベース操作:"
	@echo "  make db-up           - データベースコンテナのみを起動"
	@echo "  make db-down         - データベースコンテナのみを停止"
	@echo "  make db-reset        - データベースをリセット"
	@echo "  make migrate         - マイグレーションを実行"
	@echo "  make seed            - テストデータを生成"
	@echo ""
	@echo "APIテスト:"
	@echo "  make api-register    - 新規ユーザー登録"
	@echo "  make api-login       - ログイン"
	@echo "  make api-tweet       - ツイートを投稿"
	@echo "  make api-profile     - プロフィールを取得"
	@echo "  make api-timeline    - タイムラインを取得"
	@echo "  make api-follow      - ユーザーをフォロー"
	@echo "  make api-unfollow    - フォローを解除"
	@echo "  make api-like        - ツイートにいいね"
	@echo "  make api-unlike      - いいねを解除"
	@echo ""
	@echo "その他:"
	@echo "  make clean           - ビルドファイルを削除"
	@echo "  make install         - 依存関係をインストール"
	@echo "  make setup           - 開発環境の完全セットアップ"

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
dev: docker-up migrate generate-models run

# データベースの再構築とアプリケーションの起動
reset: docker-reset generate-models

# 開発環境の完全セットアップ(データ含む)
setup: docker-reset generate-models seed

# APIテストコマンド
api-register:
	./scripts/api_test.sh register

api-login:
	./scripts/api_test.sh login

api-tweet:
	./scripts/api_test.sh tweet

api-profile:
	./scripts/api_test.sh profile 1

api-timeline:
	./scripts/api_test.sh timeline

api-follow:
	./scripts/api_test.sh follow 2

api-unfollow:
	./scripts/api_test.sh unfollow 2

api-like:
	./scripts/api_test.sh like 1

api-unlike:
	./scripts/api_test.sh unlike 1

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
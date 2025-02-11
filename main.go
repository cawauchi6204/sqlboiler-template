package main

import (
	"log"
	"net/http"
	"os"

	"todoapp/internal/infrastructure"
	"todoapp/internal/user/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// データベース接続
	db, err := infrastructure.NewDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		log.Fatal("データベース接続エラー: ", err)
	}
	defer db.Close()

	// ハンドラーの初期化
	userHandler := handler.NewUserHandler(db)

	// Echoの初期化
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 認証不要のエンドポイント
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	// 認証が必要なエンドポイント
	api := e.Group("/api")
	api.Use(userHandler.AuthMiddleware)

	// ユーザー関連
	users := api.Group("/users")
	users.GET("/:id", userHandler.GetProfile)
	users.PUT("/me", userHandler.UpdateProfile)

	// サーバー起動
	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatal("サーバーのシャットダウン中にエラーが発生しました: ", err)
	}
}

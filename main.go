package main

import (
	"log"
	"net/http"
	"os"

	"todoapp/internal/handler"
	"todoapp/internal/infrastructure"
	"todoapp/internal/repository"
	"todoapp/internal/usecase"

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

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	tweetRepo := repository.NewTweetRepository(db)
	followRepo := repository.NewFollowRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// ユースケースの初期化
	userUsecase := usecase.NewUserUsecase(userRepo, followRepo)
	tweetUsecase := usecase.NewTweetUsecase(tweetRepo, likeRepo, userRepo)
	followUsecase := usecase.NewFollowUsecase(followRepo, userRepo)
	likeUsecase := usecase.NewLikeUsecase(likeRepo, tweetRepo)

	// ハンドラーの初期化
	h := handler.NewHandlers(
		userUsecase,
		tweetUsecase,
		followUsecase,
		likeUsecase,
		"your-secret-key", // 本番環境では環境変数から取得すべき
	)

	// Echoの初期化
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 認証不要のエンドポイント
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)

	// 認証が必要なエンドポイント
	api := e.Group("/api")
	api.Use(h.AuthMiddleware)

	// ユーザー関連
	api.GET("/users/:id", h.GetUserProfile)

	// ツイート関連
	api.POST("/tweets", h.CreateTweet)
	api.GET("/tweets/:id", h.GetTweet)
	api.DELETE("/tweets/:id", h.DeleteTweet)
	api.GET("/users/:id/tweets", h.GetUserTweets)
	api.GET("/timeline", h.GetTimeline)

	// フォロー関連
	api.POST("/users/:id/follow", h.Follow)
	api.DELETE("/users/:id/follow", h.Unfollow)
	api.GET("/users/:id/followers", h.GetFollowers)
	api.GET("/users/:id/following", h.GetFollowing)

	// いいね関連
	api.POST("/tweets/:id/like", h.LikeTweet)
	api.DELETE("/tweets/:id/like", h.UnlikeTweet)
	api.GET("/users/:id/likes", h.GetLikedTweets)

	// サーバー起動
	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatal("サーバーのシャットダウン中にエラーが発生しました: ", err)
	}
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"todoapp/internal/infrastructure"
	"todoapp/internal/schema"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"
)

const (
	NumUsers    = 1000  // 生成するユーザー数
	NumTweets   = 10000 // 生成するツイート数
	NumFollows  = 5000  // 生成するフォロー関係数
	NumLikes    = 8000  // 生成するいいね数
	BatchSize   = 100   // 一度に挿入するレコード数
	MaxTweetLen = 280   // ツイートの最大文字数
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

	// シード値を設定
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(time.Now().UnixNano())

	// ユーザーの生成
	userIDs := generateUsers(db)
	fmt.Printf("生成されたユーザー: %d件\n", len(userIDs))

	// ツイートの生成
	tweetIDs := generateTweets(db, userIDs)
	fmt.Printf("生成されたツイート: %d件\n", len(tweetIDs))

	// フォロー関係の生成
	generateFollows(db, userIDs)
	fmt.Printf("生成されたフォロー関係: %d件\n", NumFollows)

	// いいねの生成
	generateLikes(db, userIDs, tweetIDs)
	fmt.Printf("生成されたいいね: %d件\n", NumLikes)
}

func generateUsers(db *sql.DB) []int {
	userIDs := make([]int, 0, NumUsers)
	users := make([]*schema.User, 0, BatchSize)

	for i := 0; i < NumUsers; i++ {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &schema.User{
			Username:        gofakeit.Username(),
			DisplayName:     gofakeit.Name(),
			Email:           gofakeit.Email(),
			PasswordHash:    string(hashedPassword),
			Bio:             null.StringFrom(gofakeit.Sentence(10)),
			ProfileImageURL: null.StringFrom(gofakeit.ImageURL(400, 400)),
		}
		users = append(users, user)

		// バッチサイズに達したらインサート
		if len(users) == BatchSize || i == NumUsers-1 {
			err := user.Insert(context.Background(), db, boil.Infer())
			if err != nil {
				log.Printf("ユーザー作成エラー: %v\n", err)
				continue
			}
			userIDs = append(userIDs, user.ID)
			users = users[:0] // スライスをクリア
		}
	}

	return userIDs
}

func generateTweets(db *sql.DB, userIDs []int) []int {
	tweetIDs := make([]int, 0, NumTweets)
	tweets := make([]*schema.Tweet, 0, BatchSize)

	for i := 0; i < NumTweets; i++ {
		// ツイート内容の生成(最大280文字)
		content := gofakeit.Sentence(rand.Intn(20) + 1) // 1-20文の文章
		if len(content) > MaxTweetLen {
			content = content[:MaxTweetLen]
		}

		// 画像付きツイートの確率(30%)
		var imageURL null.String
		if rand.Float32() < 0.3 {
			imageURL = null.StringFrom(gofakeit.ImageURL(1200, 675)) // 16:9のアスペクト比
		}

		tweet := &schema.Tweet{
			UserID:   userIDs[rand.Intn(len(userIDs))],
			Content:  content,
			ImageURL: imageURL,
		}
		tweets = append(tweets, tweet)

		// バッチサイズに達したらインサート
		if len(tweets) == BatchSize || i == NumTweets-1 {
			err := tweet.Insert(context.Background(), db, boil.Infer())
			if err != nil {
				log.Printf("ツイート作成エラー: %v\n", err)
				continue
			}
			tweetIDs = append(tweetIDs, tweet.ID)
			tweets = tweets[:0] // スライスをクリア
		}
	}

	return tweetIDs
}

func generateFollows(db *sql.DB, userIDs []int) {
	follows := make([]*schema.Follow, 0, BatchSize)
	seen := make(map[string]bool)

	for i := 0; i < NumFollows; i++ {
		followerID := userIDs[rand.Intn(len(userIDs))]
		followingID := userIDs[rand.Intn(len(userIDs))]

		// 自分自身をフォローしない、重複もスキップ
		key := fmt.Sprintf("%d-%d", followerID, followingID)
		if followerID == followingID || seen[key] {
			continue
		}
		seen[key] = true

		follow := &schema.Follow{
			FollowerID:  followerID,
			FollowingID: followingID,
		}
		follows = append(follows, follow)

		// バッチサイズに達したらインサート
		if len(follows) == BatchSize || i == NumFollows-1 {
			err := follow.Insert(context.Background(), db, boil.Infer())
			if err != nil {
				log.Printf("フォロー作成エラー: %v\n", err)
				continue
			}
			follows = follows[:0] // スライスをクリア
		}
	}
}

func generateLikes(db *sql.DB, userIDs []int, tweetIDs []int) {
	likes := make([]*schema.Like, 0, BatchSize)
	seen := make(map[string]bool)

	for i := 0; i < NumLikes; i++ {
		userID := userIDs[rand.Intn(len(userIDs))]
		tweetID := tweetIDs[rand.Intn(len(tweetIDs))]

		// 重複をスキップ
		key := fmt.Sprintf("%d-%d", userID, tweetID)
		if seen[key] {
			continue
		}
		seen[key] = true

		like := &schema.Like{
			UserID:  userID,
			TweetID: tweetID,
		}
		likes = append(likes, like)

		// バッチサイズに達したらインサート
		if len(likes) == BatchSize || i == NumLikes-1 {
			err := like.Insert(context.Background(), db, boil.Infer())
			if err != nil {
				log.Printf("いいね作成エラー: %v\n", err)
				continue
			}
			likes = likes[:0] // スライスをクリア
		}
	}
}

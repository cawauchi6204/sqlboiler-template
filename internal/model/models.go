package model

import "time"

type User struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	DisplayName     string    `json:"display_name"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"` // パスワードハッシュはJSONに含めない
	Bio             *string   `json:"bio,omitempty"`
	ProfileImageURL *string   `json:"profile_image_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Tweet struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	ImageURL  *string   `json:"image_url,omitempty"`
	User      *User     `json:"user,omitempty"` // ツイート取得時にユーザー情報も含める
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Follow struct {
	FollowerID  int       `json:"follower_id"`
	FollowingID int       `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Like struct {
	UserID    int       `json:"user_id"`
	TweetID   int       `json:"tweet_id"`
	CreatedAt time.Time `json:"created_at"`
}

// レスポンス用の構造体
type UserResponse struct {
	User           User `json:"user"`
	FollowersCount int  `json:"followers_count"`
	FollowingCount int  `json:"following_count"`
	IsFollowing    bool `json:"is_following"`
}

type TweetResponse struct {
	Tweet   Tweet `json:"tweet"`
	IsLiked bool  `json:"is_liked"`
}

// リクエスト用の構造体
type CreateUserRequest struct {
	Username    string  `json:"username" validate:"required,min=3,max=50"`
	DisplayName string  `json:"display_name" validate:"required,max=100"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	Bio         *string `json:"bio"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateTweetRequest struct {
	Content  string  `json:"content" validate:"required,max=280"`
	ImageURL *string `json:"image_url"`
}

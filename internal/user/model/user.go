package model

import "time"

type User struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	DisplayName     string    `json:"display_name"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Bio             *string   `json:"bio,omitempty"`
	ProfileImageURL *string   `json:"profile_image_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserProfile struct {
	User           User `json:"user"`
	FollowersCount int  `json:"followers_count"`
	FollowingCount int  `json:"following_count"`
	IsFollowing    bool `json:"is_following"`
}

type RegisterRequest struct {
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

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UpdateProfileRequest struct {
	DisplayName     string  `json:"display_name" validate:"required,max=100"`
	Bio             *string `json:"bio"`
	ProfileImageURL *string `json:"profile_image_url"`
}

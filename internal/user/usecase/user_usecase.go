package usecase

import (
	"database/sql"
	"errors"
	"time"
	"todoapp/internal/user/model"
	"todoapp/internal/user/repository"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(req *model.RegisterRequest) (*model.User, error)
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	GetProfile(userID, currentUserID int) (*model.UserProfile, error)
	UpdateProfile(userID int, req *model.UpdateProfileRequest) (*model.User, error)
}

type userUsecase struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserUsecase(db *sql.DB, jwtSecret string) UserUsecase {
	return &userUsecase{
		repo:      repository.NewUserRepository(db),
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) Register(req *model.RegisterRequest) (*model.User, error) {
	// メールアドレスの重複チェック
	existingUser, err := u.repo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	return u.repo.Create(req)
}

func (u *userUsecase) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// パスワードの検証
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 1週間の有効期限
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: tokenString,
		User:  *user,
	}, nil
}

func (u *userUsecase) GetProfile(userID, currentUserID int) (*model.UserProfile, error) {
	return u.repo.GetProfile(userID, currentUserID)
}

func (u *userUsecase) UpdateProfile(userID int, req *model.UpdateProfileRequest) (*model.User, error) {
	return u.repo.Update(userID, req)
}

package repository

import (
	"context"
	"database/sql"
	"todoapp/internal/schema"
	"todoapp/internal/user/model"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *model.RegisterRequest) (*model.User, error)
	GetByID(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(id int, user *model.UpdateProfileRequest) (*model.User, error)
	GetProfile(userID, currentUserID int) (*model.UserProfile, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) convertToModel(dbUser *schema.User) *model.User {
	return &model.User{
		ID:              dbUser.ID,
		Username:        dbUser.Username,
		DisplayName:     dbUser.DisplayName,
		Email:           dbUser.Email,
		PasswordHash:    dbUser.PasswordHash,
		Bio:             &dbUser.Bio.String,
		ProfileImageURL: &dbUser.ProfileImageURL.String,
		CreatedAt:       dbUser.CreatedAt.Time,
		UpdatedAt:       dbUser.UpdatedAt.Time,
	}
}

func (r *userRepository) Create(req *model.RegisterRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	dbUser := &schema.User{
		Username:        req.Username,
		DisplayName:     req.DisplayName,
		Email:           req.Email,
		PasswordHash:    string(hashedPassword),
		Bio:             null.StringFromPtr(req.Bio),
		ProfileImageURL: null.String{},
	}

	err = dbUser.Insert(context.Background(), r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return r.convertToModel(dbUser), nil
}

func (r *userRepository) GetByID(id int) (*model.User, error) {
	dbUser, err := schema.FindUser(context.Background(), r.db, id)
	if err != nil {
		return nil, err
	}

	return r.convertToModel(dbUser), nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	dbUser, err := schema.Users(qm.Where("email = ?", email)).One(context.Background(), r.db)
	if err != nil {
		return nil, err
	}

	return r.convertToModel(dbUser), nil
}

func (r *userRepository) Update(id int, req *model.UpdateProfileRequest) (*model.User, error) {
	dbUser, err := schema.FindUser(context.Background(), r.db, id)
	if err != nil {
		return nil, err
	}

	dbUser.DisplayName = req.DisplayName
	dbUser.Bio = null.StringFromPtr(req.Bio)
	dbUser.ProfileImageURL = null.StringFromPtr(req.ProfileImageURL)

	_, err = dbUser.Update(context.Background(), r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return r.convertToModel(dbUser), nil
}

func (r *userRepository) GetProfile(userID, currentUserID int) (*model.UserProfile, error) {
	ctx := context.Background()

	dbUser, err := schema.FindUser(ctx, r.db, userID)
	if err != nil {
		return nil, err
	}

	// フォロワー数を取得
	followersCount, err := schema.Follows(
		qm.Where("following_id = ?", userID),
	).Count(ctx, r.db)
	if err != nil {
		return nil, err
	}

	// フォロー数を取得
	followingCount, err := schema.Follows(
		qm.Where("follower_id = ?", userID),
	).Count(ctx, r.db)
	if err != nil {
		return nil, err
	}

	// フォロー状態を確認
	isFollowing := false
	if currentUserID != 0 {
		exists, err := schema.Follows(
			qm.Where("follower_id = ? AND following_id = ?", currentUserID, userID),
		).Exists(ctx, r.db)
		if err != nil {
			return nil, err
		}
		isFollowing = exists
	}

	return &model.UserProfile{
		User:           *r.convertToModel(dbUser),
		FollowersCount: int(followersCount),
		FollowingCount: int(followingCount),
		IsFollowing:    isFollowing,
	}, nil
}

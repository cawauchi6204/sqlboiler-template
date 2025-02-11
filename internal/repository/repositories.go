package repository

import (
	"context"
	"database/sql"
	"errors"
	"todoapp/internal/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type TweetRepository interface {
	Create(tweet *models.Tweet) error
	GetByID(id int) (*models.Tweet, error)
	GetByUserID(userID int, limit, offset int) ([]*models.Tweet, error)
	GetTimeline(userID int, limit, offset int) ([]*models.Tweet, error)
	Update(tweet *models.Tweet) error
	Delete(id int) error
}

type FollowRepository interface {
	Create(follow *models.Follow) error
	Delete(followerID, followingID int) error
	IsFollowing(followerID, followingID int) (bool, error)
	GetFollowers(userID int, limit, offset int) ([]*models.User, error)
	GetFollowing(userID int, limit, offset int) ([]*models.User, error)
	GetFollowersCount(userID int) (int64, error)
	GetFollowingCount(userID int) (int64, error)
}

type LikeRepository interface {
	Create(like *models.Like) error
	Delete(userID, tweetID int) error
	IsLiked(userID, tweetID int) (bool, error)
	GetLikeCount(tweetID int) (int64, error)
	GetLikedTweets(userID int, limit, offset int) ([]*models.Tweet, error)
}

type userRepository struct {
	db *sql.DB
}

type tweetRepository struct {
	db *sql.DB
}

type followRepository struct {
	db *sql.DB
}

type likeRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func NewTweetRepository(db *sql.DB) TweetRepository {
	return &tweetRepository{db: db}
}

func NewFollowRepository(db *sql.DB) FollowRepository {
	return &followRepository{db: db}
}

func NewLikeRepository(db *sql.DB) LikeRepository {
	return &likeRepository{db: db}
}

// UserRepository実装
func (r *userRepository) Create(user *models.User) error {
	return user.Insert(context.Background(), r.db, boil.Infer())
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	user, err := models.Users(qm.Where("id=?", id)).One(context.Background(), r.db)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return user, err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user, err := models.Users(qm.Where("email=?", email)).One(context.Background(), r.db)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return user, err
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	user, err := models.Users(qm.Where("username=?", username)).One(context.Background(), r.db)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	_, err := user.Update(context.Background(), r.db, boil.Infer())
	return err
}

func (r *userRepository) Delete(id int) error {
	user, err := r.GetByID(id)
	if err != nil {
		return err
	}
	_, err = user.Delete(context.Background(), r.db)
	return err
}

// TweetRepository実装
func (r *tweetRepository) Create(tweet *models.Tweet) error {
	return tweet.Insert(context.Background(), r.db, boil.Infer())
}

func (r *tweetRepository) GetByID(id int) (*models.Tweet, error) {
	tweet, err := models.Tweets(
		qm.Where("id=?", id),
		qm.Load("User"),
	).One(context.Background(), r.db)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return tweet, err
}

func (r *tweetRepository) GetByUserID(userID int, limit, offset int) ([]*models.Tweet, error) {
	return models.Tweets(
		qm.Where("user_id=?", userID),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
		qm.Load("User"),
	).All(context.Background(), r.db)
}

func (r *tweetRepository) GetTimeline(userID int, limit, offset int) ([]*models.Tweet, error) {
	return models.Tweets(
		qm.Where("user_id IN (SELECT following_id FROM follows WHERE follower_id = ?) OR user_id = ?", userID, userID),
		qm.OrderBy("created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
		qm.Load("User"),
	).All(context.Background(), r.db)
}

func (r *tweetRepository) Update(tweet *models.Tweet) error {
	_, err := tweet.Update(context.Background(), r.db, boil.Infer())
	return err
}

func (r *tweetRepository) Delete(id int) error {
	tweet, err := r.GetByID(id)
	if err != nil {
		return err
	}
	_, err = tweet.Delete(context.Background(), r.db)
	return err
}

// FollowRepository実装
func (r *followRepository) Create(follow *models.Follow) error {
	return follow.Insert(context.Background(), r.db, boil.Infer())
}

func (r *followRepository) Delete(followerID, followingID int) error {
	_, err := models.Follows(
		qm.Where("follower_id=? AND following_id=?", followerID, followingID),
	).DeleteAll(context.Background(), r.db)
	return err
}

func (r *followRepository) IsFollowing(followerID, followingID int) (bool, error) {
	exists, err := models.Follows(
		qm.Where("follower_id=? AND following_id=?", followerID, followingID),
	).Exists(context.Background(), r.db)
	return exists, err
}

func (r *followRepository) GetFollowers(userID int, limit, offset int) ([]*models.User, error) {
	return models.Users(
		qm.InnerJoin("follows f on f.follower_id = users.id"),
		qm.Where("f.following_id=?", userID),
		qm.OrderBy("f.created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(context.Background(), r.db)
}

func (r *followRepository) GetFollowing(userID int, limit, offset int) ([]*models.User, error) {
	return models.Users(
		qm.InnerJoin("follows f on f.following_id = users.id"),
		qm.Where("f.follower_id=?", userID),
		qm.OrderBy("f.created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(context.Background(), r.db)
}

func (r *followRepository) GetFollowersCount(userID int) (int64, error) {
	return models.Follows(
		qm.Where("following_id=?", userID),
	).Count(context.Background(), r.db)
}

func (r *followRepository) GetFollowingCount(userID int) (int64, error) {
	return models.Follows(
		qm.Where("follower_id=?", userID),
	).Count(context.Background(), r.db)
}

// LikeRepository実装
func (r *likeRepository) Create(like *models.Like) error {
	return like.Insert(context.Background(), r.db, boil.Infer())
}

func (r *likeRepository) Delete(userID, tweetID int) error {
	_, err := models.Likes(
		qm.Where("user_id=? AND tweet_id=?", userID, tweetID),
	).DeleteAll(context.Background(), r.db)
	return err
}

func (r *likeRepository) IsLiked(userID, tweetID int) (bool, error) {
	exists, err := models.Likes(
		qm.Where("user_id=? AND tweet_id=?", userID, tweetID),
	).Exists(context.Background(), r.db)
	return exists, err
}

func (r *likeRepository) GetLikeCount(tweetID int) (int64, error) {
	return models.Likes(
		qm.Where("tweet_id=?", tweetID),
	).Count(context.Background(), r.db)
}

func (r *likeRepository) GetLikedTweets(userID int, limit, offset int) ([]*models.Tweet, error) {
	return models.Tweets(
		qm.InnerJoin("likes l on l.tweet_id = tweets.id"),
		qm.Where("l.user_id=?", userID),
		qm.OrderBy("l.created_at DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
		qm.Load("User"),
	).All(context.Background(), r.db)
}

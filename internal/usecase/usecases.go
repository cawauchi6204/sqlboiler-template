package usecase

import (
	"errors"
	"todoapp/internal/model"
	"todoapp/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrTweetNotFound      = errors.New("tweet not found")
	ErrDuplicateUsername  = errors.New("username already exists")
	ErrDuplicateEmail     = errors.New("email already exists")
)

type UserUsecase interface {
	Register(req *model.CreateUserRequest) (*model.User, error)
	Login(req *model.LoginRequest) (*model.User, error)
	GetUserProfile(userID int, currentUserID int) (*model.UserResponse, error)
	UpdateProfile(userID int, displayName string, bio *string) error
	UpdateProfileImage(userID int, imageURL string) error
}

type TweetUsecase interface {
	CreateTweet(userID int, req *model.CreateTweetRequest) (*model.Tweet, error)
	GetTweet(tweetID int, currentUserID int) (*model.TweetResponse, error)
	GetUserTweets(userID int, currentUserID int, limit, offset int) ([]*model.TweetResponse, error)
	GetTimeline(userID int, limit, offset int) ([]*model.TweetResponse, error)
	DeleteTweet(tweetID, userID int) error
}

type FollowUsecase interface {
	Follow(followerID, followingID int) error
	Unfollow(followerID, followingID int) error
	GetFollowers(userID int, limit, offset int) ([]*model.User, error)
	GetFollowing(userID int, limit, offset int) ([]*model.User, error)
}

type LikeUsecase interface {
	LikeTweet(userID, tweetID int) error
	UnlikeTweet(userID, tweetID int) error
	GetLikedTweets(userID int, limit, offset int) ([]*model.Tweet, error)
}

type userUsecase struct {
	userRepo   repository.UserRepository
	followRepo repository.FollowRepository
}

type tweetUsecase struct {
	tweetRepo repository.TweetRepository
	likeRepo  repository.LikeRepository
	userRepo  repository.UserRepository
}

type followUsecase struct {
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
}

type likeUsecase struct {
	likeRepo  repository.LikeRepository
	tweetRepo repository.TweetRepository
}

func NewUserUsecase(userRepo repository.UserRepository, followRepo repository.FollowRepository) UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func NewTweetUsecase(tweetRepo repository.TweetRepository, likeRepo repository.LikeRepository, userRepo repository.UserRepository) TweetUsecase {
	return &tweetUsecase{
		tweetRepo: tweetRepo,
		likeRepo:  likeRepo,
		userRepo:  userRepo,
	}
}

func NewFollowUsecase(followRepo repository.FollowRepository, userRepo repository.UserRepository) FollowUsecase {
	return &followUsecase{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

func NewLikeUsecase(likeRepo repository.LikeRepository, tweetRepo repository.TweetRepository) LikeUsecase {
	return &likeUsecase{
		likeRepo:  likeRepo,
		tweetRepo: tweetRepo,
	}
}

// UserUsecase実装
func (u *userUsecase) Register(req *model.CreateUserRequest) (*model.User, error) {
	// ユーザー名とメールアドレスの重複チェック
	if _, err := u.userRepo.GetByUsername(req.Username); err == nil {
		return nil, ErrDuplicateUsername
	}
	if _, err := u.userRepo.GetByEmail(req.Email); err == nil {
		return nil, ErrDuplicateEmail
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Bio:          req.Bio,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) Login(req *model.LoginRequest) (*model.User, error) {
	user, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (u *userUsecase) GetUserProfile(userID int, currentUserID int) (*model.UserResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	followersCount, _ := u.followRepo.GetFollowersCount(userID)
	followingCount, _ := u.followRepo.GetFollowingCount(userID)
	isFollowing, _ := u.followRepo.IsFollowing(currentUserID, userID)

	return &model.UserResponse{
		User:           *user,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		IsFollowing:    isFollowing,
	}, nil
}

func (u *userUsecase) UpdateProfile(userID int, displayName string, bio *string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	user.DisplayName = displayName
	user.Bio = bio

	return u.userRepo.Update(user)
}

func (u *userUsecase) UpdateProfileImage(userID int, imageURL string) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	user.ProfileImageURL = &imageURL

	return u.userRepo.Update(user)
}

// TweetUsecase実装
func (u *tweetUsecase) CreateTweet(userID int, req *model.CreateTweetRequest) (*model.Tweet, error) {
	tweet := &model.Tweet{
		UserID:   userID,
		Content:  req.Content,
		ImageURL: req.ImageURL,
	}

	if err := u.tweetRepo.Create(tweet); err != nil {
		return nil, err
	}

	return tweet, nil
}

func (u *tweetUsecase) GetTweet(tweetID int, currentUserID int) (*model.TweetResponse, error) {
	tweet, err := u.tweetRepo.GetByID(tweetID)
	if err != nil {
		return nil, ErrTweetNotFound
	}

	isLiked, _ := u.likeRepo.IsLiked(currentUserID, tweetID)

	return &model.TweetResponse{
		Tweet:   *tweet,
		IsLiked: isLiked,
	}, nil
}

func (u *tweetUsecase) GetUserTweets(userID int, currentUserID int, limit, offset int) ([]*model.TweetResponse, error) {
	tweets, err := u.tweetRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.TweetResponse, len(tweets))
	for i, tweet := range tweets {
		isLiked, _ := u.likeRepo.IsLiked(currentUserID, tweet.ID)
		responses[i] = &model.TweetResponse{
			Tweet:   *tweet,
			IsLiked: isLiked,
		}
	}

	return responses, nil
}

func (u *tweetUsecase) GetTimeline(userID int, limit, offset int) ([]*model.TweetResponse, error) {
	tweets, err := u.tweetRepo.GetTimeline(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.TweetResponse, len(tweets))
	for i, tweet := range tweets {
		isLiked, _ := u.likeRepo.IsLiked(userID, tweet.ID)
		responses[i] = &model.TweetResponse{
			Tweet:   *tweet,
			IsLiked: isLiked,
		}
	}

	return responses, nil
}

func (u *tweetUsecase) DeleteTweet(tweetID, userID int) error {
	tweet, err := u.tweetRepo.GetByID(tweetID)
	if err != nil {
		return ErrTweetNotFound
	}

	if tweet.UserID != userID {
		return errors.New("unauthorized")
	}

	return u.tweetRepo.Delete(tweetID)
}

// FollowUsecase実装
func (u *followUsecase) Follow(followerID, followingID int) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	isFollowing, err := u.followRepo.IsFollowing(followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return errors.New("already following")
	}

	follow := &model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return u.followRepo.Create(follow)
}

func (u *followUsecase) Unfollow(followerID, followingID int) error {
	isFollowing, err := u.followRepo.IsFollowing(followerID, followingID)
	if err != nil {
		return err
	}
	if !isFollowing {
		return errors.New("not following")
	}

	return u.followRepo.Delete(followerID, followingID)
}

func (u *followUsecase) GetFollowers(userID int, limit, offset int) ([]*model.User, error) {
	return u.followRepo.GetFollowers(userID, limit, offset)
}

func (u *followUsecase) GetFollowing(userID int, limit, offset int) ([]*model.User, error) {
	return u.followRepo.GetFollowing(userID, limit, offset)
}

// LikeUsecase実装
func (u *likeUsecase) LikeTweet(userID, tweetID int) error {
	isLiked, err := u.likeRepo.IsLiked(userID, tweetID)
	if err != nil {
		return err
	}
	if isLiked {
		return errors.New("already liked")
	}

	like := &model.Like{
		UserID:  userID,
		TweetID: tweetID,
	}
	return u.likeRepo.Create(like)
}

func (u *likeUsecase) UnlikeTweet(userID, tweetID int) error {
	isLiked, err := u.likeRepo.IsLiked(userID, tweetID)
	if err != nil {
		return err
	}
	if !isLiked {
		return errors.New("not liked")
	}

	return u.likeRepo.Delete(userID, tweetID)
}

func (u *likeUsecase) GetLikedTweets(userID int, limit, offset int) ([]*model.Tweet, error) {
	return u.likeRepo.GetLikedTweets(userID, limit, offset)
}

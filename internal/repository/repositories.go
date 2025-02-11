package repository

import (
	"database/sql"
	"errors"
	"todoapp/internal/model"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
}

type TweetRepository interface {
	Create(tweet *model.Tweet) error
	GetByID(id int) (*model.Tweet, error)
	GetByUserID(userID int, limit, offset int) ([]*model.Tweet, error)
	GetTimeline(userID int, limit, offset int) ([]*model.Tweet, error)
	Update(tweet *model.Tweet) error
	Delete(id int) error
}

type FollowRepository interface {
	Create(follow *model.Follow) error
	Delete(followerID, followingID int) error
	IsFollowing(followerID, followingID int) (bool, error)
	GetFollowers(userID int, limit, offset int) ([]*model.User, error)
	GetFollowing(userID int, limit, offset int) ([]*model.User, error)
	GetFollowersCount(userID int) (int, error)
	GetFollowingCount(userID int) (int, error)
}

type LikeRepository interface {
	Create(like *model.Like) error
	Delete(userID, tweetID int) error
	IsLiked(userID, tweetID int) (bool, error)
	GetLikeCount(tweetID int) (int, error)
	GetLikedTweets(userID int, limit, offset int) ([]*model.Tweet, error)
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
func (r *userRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (username, display_name, email, password_hash, bio, profile_image_url)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		user.Username,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
		user.Bio,
		user.ProfileImageURL,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

func (r *userRepository) GetByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, display_name, email, password_hash, bio, profile_image_url, created_at, updated_at
		FROM users WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ProfileImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, display_name, email, password_hash, bio, profile_image_url, created_at, updated_at
		FROM users WHERE email = ?
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ProfileImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, display_name, email, password_hash, bio, profile_image_url, created_at, updated_at
		FROM users WHERE username = ?
	`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ProfileImageURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(user *model.User) error {
	query := `
		UPDATE users
		SET username = ?, display_name = ?, email = ?, bio = ?, profile_image_url = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query,
		user.Username,
		user.DisplayName,
		user.Email,
		user.Bio,
		user.ProfileImageURL,
		user.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// TweetRepository実装
func (r *tweetRepository) Create(tweet *model.Tweet) error {
	query := `
		INSERT INTO tweets (user_id, content, image_url)
		VALUES (?, ?, ?)
	`
	result, err := r.db.Exec(query,
		tweet.UserID,
		tweet.Content,
		tweet.ImageURL,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	tweet.ID = int(id)
	return nil
}

func (r *tweetRepository) GetByID(id int) (*model.Tweet, error) {
	tweet := &model.Tweet{}
	query := `
		SELECT t.id, t.user_id, t.content, t.image_url, t.created_at, t.updated_at,
			   u.username, u.display_name, u.profile_image_url,
			   (SELECT COUNT(*) FROM likes WHERE tweet_id = t.id) as like_count
		FROM tweets t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?
	`
	var username, displayName string
	var profileImageURL *string
	err := r.db.QueryRow(query, id).Scan(
		&tweet.ID,
		&tweet.UserID,
		&tweet.Content,
		&tweet.ImageURL,
		&tweet.CreatedAt,
		&tweet.UpdatedAt,
		&username,
		&displayName,
		&profileImageURL,
		&tweet.LikeCount,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	tweet.User = &model.User{
		Username:        username,
		DisplayName:     displayName,
		ProfileImageURL: profileImageURL,
	}
	return tweet, nil
}

func (r *tweetRepository) GetByUserID(userID int, limit, offset int) ([]*model.Tweet, error) {
	query := `
		SELECT t.id, t.user_id, t.content, t.image_url, t.created_at, t.updated_at,
			   u.username, u.display_name, u.profile_image_url,
			   (SELECT COUNT(*) FROM likes WHERE tweet_id = t.id) as like_count
		FROM tweets t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id = ?
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`
	return r.queryTweets(query, userID, limit, offset)
}

func (r *tweetRepository) GetTimeline(userID int, limit, offset int) ([]*model.Tweet, error) {
	query := `
		SELECT t.id, t.user_id, t.content, t.image_url, t.created_at, t.updated_at,
			   u.username, u.display_name, u.profile_image_url,
			   (SELECT COUNT(*) FROM likes WHERE tweet_id = t.id) as like_count
		FROM tweets t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id IN (
			SELECT following_id FROM follows WHERE follower_id = ?
		) OR t.user_id = ?
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`
	return r.queryTweets(query, userID, userID, limit, offset)
}

func (r *tweetRepository) queryTweets(query string, args ...interface{}) ([]*model.Tweet, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []*model.Tweet
	for rows.Next() {
		tweet := &model.Tweet{}
		var username, displayName string
		var profileImageURL *string
		err := rows.Scan(
			&tweet.ID,
			&tweet.UserID,
			&tweet.Content,
			&tweet.ImageURL,
			&tweet.CreatedAt,
			&tweet.UpdatedAt,
			&username,
			&displayName,
			&profileImageURL,
			&tweet.LikeCount,
		)
		if err != nil {
			return nil, err
		}

		tweet.User = &model.User{
			Username:        username,
			DisplayName:     displayName,
			ProfileImageURL: profileImageURL,
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}

func (r *tweetRepository) Update(tweet *model.Tweet) error {
	query := `
		UPDATE tweets
		SET content = ?, image_url = ?
		WHERE id = ? AND user_id = ?
	`
	result, err := r.db.Exec(query,
		tweet.Content,
		tweet.ImageURL,
		tweet.ID,
		tweet.UserID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *tweetRepository) Delete(id int) error {
	query := "DELETE FROM tweets WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// FollowRepository実装
func (r *followRepository) Create(follow *model.Follow) error {
	query := "INSERT INTO follows (follower_id, following_id) VALUES (?, ?)"
	_, err := r.db.Exec(query, follow.FollowerID, follow.FollowingID)
	return err
}

func (r *followRepository) Delete(followerID, followingID int) error {
	query := "DELETE FROM follows WHERE follower_id = ? AND following_id = ?"
	result, err := r.db.Exec(query, followerID, followingID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *followRepository) IsFollowing(followerID, followingID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = ? AND following_id = ?)"
	err := r.db.QueryRow(query, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *followRepository) GetFollowers(userID int, limit, offset int) ([]*model.User, error) {
	query := `
		SELECT u.id, u.username, u.display_name, u.email, u.bio, u.profile_image_url, u.created_at, u.updated_at
		FROM users u
		JOIN follows f ON u.id = f.follower_id
		WHERE f.following_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?
	`
	return r.queryUsers(query, userID, limit, offset)
}

func (r *followRepository) GetFollowing(userID int, limit, offset int) ([]*model.User, error) {
	query := `
		SELECT u.id, u.username, u.display_name, u.email, u.bio, u.profile_image_url, u.created_at, u.updated_at
		FROM users u
		JOIN follows f ON u.id = f.following_id
		WHERE f.follower_id = ?
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?
	`
	return r.queryUsers(query, userID, limit, offset)
}

func (r *followRepository) queryUsers(query string, args ...interface{}) ([]*model.User, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.DisplayName,
			&user.Email,
			&user.Bio,
			&user.ProfileImageURL,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *followRepository) GetFollowersCount(userID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM follows WHERE following_id = ?"
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *followRepository) GetFollowingCount(userID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM follows WHERE follower_id = ?"
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// LikeRepository実装
func (r *likeRepository) Create(like *model.Like) error {
	query := "INSERT INTO likes (user_id, tweet_id) VALUES (?, ?)"
	_, err := r.db.Exec(query, like.UserID, like.TweetID)
	return err
}

func (r *likeRepository) Delete(userID, tweetID int) error {
	query := "DELETE FROM likes WHERE user_id = ? AND tweet_id = ?"
	result, err := r.db.Exec(query, userID, tweetID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *likeRepository) IsLiked(userID, tweetID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND tweet_id = ?)"
	err := r.db.QueryRow(query, userID, tweetID).Scan(&exists)
	return exists, err
}

func (r *likeRepository) GetLikeCount(tweetID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM likes WHERE tweet_id = ?"
	err := r.db.QueryRow(query, tweetID).Scan(&count)
	return count, err
}

func (r *likeRepository) GetLikedTweets(userID int, limit, offset int) ([]*model.Tweet, error) {
	query := `
		SELECT t.id, t.user_id, t.content, t.image_url, t.created_at, t.updated_at,
			   u.username, u.display_name, u.profile_image_url,
			   (SELECT COUNT(*) FROM likes WHERE tweet_id = t.id) as like_count
		FROM tweets t
		JOIN users u ON t.user_id = u.id
		JOIN likes l ON t.id = l.tweet_id
		WHERE l.user_id = ?
		ORDER BY l.created_at DESC
		LIMIT ? OFFSET ?
	`
	tr := &tweetRepository{db: r.db}
	return tr.queryTweets(query, userID, limit, offset)
}

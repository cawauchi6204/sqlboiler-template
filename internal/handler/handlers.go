package handler

import (
	"net/http"
	"strconv"
	"time"

	"todoapp/internal/model"
	"todoapp/internal/usecase"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

type Handlers struct {
	userUsecase   usecase.UserUsecase
	tweetUsecase  usecase.TweetUsecase
	followUsecase usecase.FollowUsecase
	likeUsecase   usecase.LikeUsecase
	jwtSecret     []byte
}

func NewHandlers(
	userUsecase usecase.UserUsecase,
	tweetUsecase usecase.TweetUsecase,
	followUsecase usecase.FollowUsecase,
	likeUsecase usecase.LikeUsecase,
	jwtSecret string,
) *Handlers {
	return &Handlers{
		userUsecase:   userUsecase,
		tweetUsecase:  tweetUsecase,
		followUsecase: followUsecase,
		likeUsecase:   likeUsecase,
		jwtSecret:     []byte(jwtSecret),
	}
}

// ミドルウェア
func (h *Handlers) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return h.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		c.Set("user_id", claims.UserID)
		return next(c)
	}
}

// ユーザー関連ハンドラー
func (h *Handlers) Register(c echo.Context) error {
	var req model.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.userUsecase.Register(&req)
	if err != nil {
		switch err {
		case usecase.ErrDuplicateUsername:
			return echo.NewHTTPError(http.StatusConflict, "username already exists")
		case usecase.ErrDuplicateEmail:
			return echo.NewHTTPError(http.StatusConflict, "email already exists")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	// JWTトークンの生成
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (h *Handlers) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.userUsecase.Login(&req)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (h *Handlers) GetUserProfile(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	currentUserID := c.Get("user_id").(int)

	profile, err := h.userUsecase.GetUserProfile(userID, currentUserID)
	if err != nil {
		if err == usecase.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}

// ツイート関連ハンドラー
func (h *Handlers) CreateTweet(c echo.Context) error {
	userID := c.Get("user_id").(int)
	var req model.CreateTweetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	tweet, err := h.tweetUsecase.CreateTweet(userID, &req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, tweet)
}

func (h *Handlers) GetTweet(c echo.Context) error {
	tweetID, _ := strconv.Atoi(c.Param("id"))
	currentUserID := c.Get("user_id").(int)

	tweet, err := h.tweetUsecase.GetTweet(tweetID, currentUserID)
	if err != nil {
		if err == usecase.ErrTweetNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "tweet not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tweet)
}

func (h *Handlers) GetUserTweets(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	currentUserID := c.Get("user_id").(int)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	tweets, err := h.tweetUsecase.GetUserTweets(userID, currentUserID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tweets)
}

func (h *Handlers) GetTimeline(c echo.Context) error {
	userID := c.Get("user_id").(int)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	tweets, err := h.tweetUsecase.GetTimeline(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tweets)
}

func (h *Handlers) DeleteTweet(c echo.Context) error {
	tweetID, _ := strconv.Atoi(c.Param("id"))
	userID := c.Get("user_id").(int)

	if err := h.tweetUsecase.DeleteTweet(tweetID, userID); err != nil {
		if err == usecase.ErrTweetNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "tweet not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// フォロー関連ハンドラー
func (h *Handlers) Follow(c echo.Context) error {
	followerID := c.Get("user_id").(int)
	followingID, _ := strconv.Atoi(c.Param("id"))

	if err := h.followUsecase.Follow(followerID, followingID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handlers) Unfollow(c echo.Context) error {
	followerID := c.Get("user_id").(int)
	followingID, _ := strconv.Atoi(c.Param("id"))

	if err := h.followUsecase.Unfollow(followerID, followingID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handlers) GetFollowers(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	followers, err := h.followUsecase.GetFollowers(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, followers)
}

func (h *Handlers) GetFollowing(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	following, err := h.followUsecase.GetFollowing(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, following)
}

// いいね関連ハンドラー
func (h *Handlers) LikeTweet(c echo.Context) error {
	userID := c.Get("user_id").(int)
	tweetID, _ := strconv.Atoi(c.Param("id"))

	if err := h.likeUsecase.LikeTweet(userID, tweetID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handlers) UnlikeTweet(c echo.Context) error {
	userID := c.Get("user_id").(int)
	tweetID, _ := strconv.Atoi(c.Param("id"))

	if err := h.likeUsecase.UnlikeTweet(userID, tweetID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handlers) GetLikedTweets(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("id"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	tweets, err := h.likeUsecase.GetLikedTweets(userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tweets)
}

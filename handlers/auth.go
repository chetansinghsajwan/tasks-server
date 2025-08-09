package handlers

import (
	"context"
	"fmt"
	"net/http"
	"tasks/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	LoginRequestTimeout       = time.Second * 5
	JwtSigningKey             = ""
	AccessTokenCookiePath     = ""
	AccessTokenCookieDomain   = ""
	AccessTokenCookieLifetime = 3600
)

var ST store.Store

func GenerateToken(userID string) (string, error) {

	var claims = jwt.NewWithClaims(
		jwt.SigningMethodPS256.SigningMethodRSA,
		jwt.MapClaims{
			"id":         userID,
			"expires_at": time.Now().Add(time.Hour * 2),
		},
	)

	var token, err = claims.SignedString(JwtSigningKey)

	if err != nil {
		return "", err
	}

	return token, nil
}

func AuthenticateMiddleware(c *gin.Context) {

	var err error
	var tokenStr string
	if tokenStr, err = c.Cookie("access-token"); err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("no access token provided. Error: %s", err.Error()),
		})
		return
	}

	var token *jwt.Token
	token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		// Make sure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return JwtSigningKey, nil
	})

	if err != nil || !token.Valid {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired token",
		})
		return
	}

	// Extract username if needed
	var userIdStr string
	if userIdStr, err = token.Claims.GetSubject(); err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if userIdStr == "" {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "userid must not be empty",
		})
		return
	}

	c.Set("userid", userIdStr)
}

type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), LoginRequestTimeout)

	defer cancel()

	var body LoginRequest
	var err error
	if err = c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var serr *store.StoreError
	var secret *store.UserSecret
	secret, serr = ST.GetUserSecret(ctx, body.UserID)

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(secret.Pass), []byte(body.Password))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password.",
		})
		return
	}

	var token string
	if token, err = GenerateToken(body.UserID); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access-token", token, AccessTokenCookieLifetime,
		AccessTokenCookiePath, AccessTokenCookieDomain, true, true)

	c.Status(http.StatusOK)
}

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"tasks/db"
	"tasks/sqlc"
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

func GenerateToken(userId string) (string, error) {

	var claims = jwt.NewWithClaims(
		jwt.SigningMethodPS256.SigningMethodRSA,
		jwt.MapClaims{
			"id":         userId,
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

	var tokenStr, err = c.Cookie("access-token")
	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("no access token provided. Error: %s", err.Error()),
		})
		return
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
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
	userId, err := token.Claims.GetSubject()
	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if userId == "" {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "userid must not be empty",
		})
		return
	}

	c.Set("userid", userId)
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
	if err := c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var quries, err = db.Begin(ctx)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer quries.Tx.Rollback(ctx)

	hashedPassword, err := quries.Queries.GetSecret(ctx, sqlc.GetSecretParams{
		Key:   body.UserID,
		Scope: "user-login",
	})

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	quries.Tx.Commit(ctx)

	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(body.Password))

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password.",
		})
		return
	}

	token, err := GenerateToken(body.UserID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("access-token", token, AccessTokenCookieLifetime,
		AccessTokenCookiePath, AccessTokenCookieDomain, true, true)

	c.Status(http.StatusOK)
}

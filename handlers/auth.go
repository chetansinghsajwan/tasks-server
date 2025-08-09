package handlers

import (
	"context"
	"fmt"
	"net/http"
	"tasks/errorcodes"
	"tasks/services"
	"tasks/store"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

const (
	LoginRequestTimeout = time.Second * 5
)

var ST store.Store

func AuthenticateMiddleware(c *gin.Context) {

	var err error
	var tokenStr string
	if tokenStr, err = c.Cookie("access-token"); err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"msg": "access token not found",
			},
		})
		return
	}

	var serr *services.ServiceError
	var authToken *services.AuthToken
	authToken, serr = services.AuthenticateToken(tokenStr)

	if serr != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"msg": "invalid token",
			},
		})
		return
	}

	c.Set("userid", authToken.UserID)
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

	var token string
	var serr *services.ServiceError
	token, serr = services.LoginUser(
		services.ServiceContext{
			Ctx: ctx,
		},
		services.LoginParams{
			UserID:   body.UserID,
			Password: body.Password,
		},
	)

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound:

			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"msg": fmt.Sprintf("user with id '%s' not found", body.UserID),
				},
			})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"msg": "internal server error",
				},
			})
			return
		}
	}

	c.SetCookie(
		"access-token",
		token,
		services.AccessTokenCookieLifetime,
		services.AccessTokenCookiePath,
		services.AccessTokenCookieDomain,
		true,
		true,
	)

	c.Status(http.StatusOK)
}

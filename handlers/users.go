package handlers

import (
	"context"
	"net/http"
	"strings"
	"tasks/errorcodes"
	"tasks/option"
	"tasks/store"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	CreateUserRequestTimeout = time.Second * 5000000000
	GetUserRequestTimeout    = time.Second * 5000000000
	UpdateUserRequestTimeout = time.Second * 5000000000
	DeleteUserRequestTimeout = time.Second * 5000000000
	BcryptUserEncryptionCost = bcrypt.DefaultCost
)

var ST store.Store

type CreateUserBody struct {
	ID          string  `json:"id"`
	Pass        string  `json:"pass"`
	FullName    string  `json:"full_name"`
	DisplayName *string `json:"display_name"`
	Email       string  `json:"email"`
}

type UpdateUserBody struct {
	ID          *string `json:"id"`
	Pass        *string `json:"pass"`
	FullName    *string `json:"full_name"`
	DisplayName *string `json:"display_name"`
	Email       *string `json:"email"`
}

func CreateUser(c *gin.Context) {

	// Setup context
	var ctx, cancel = context.WithTimeout(
		context.Background(), CreateUserRequestTimeout,
	)

	defer cancel()

	// Parse request body
	var body CreateUserBody
	var err error
	if err = c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var serr *store.StoreError = ST.CreateUser(ctx, store.CreateUserParams{
		ID:          body.ID,
		Email:       body.Email,
		FullName:    body.FullName,
		DisplayName: option.FromPtr(body.DisplayName),
	})

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserIDNull,
			errorcodes.UserIDAlreadyExists,
			errorcodes.UserIDFormat,
			errorcodes.UserEmailNull,
			errorcodes.UserEmailAlreadyExists,
			errorcodes.UserEmailFormat,
			errorcodes.UserFullNameFormat,
			errorcodes.UserDisplayNameFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": serr.Msg,
			})
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	// // Encrypt the password
	// var hashedPass []byte
	// hashedPass, err = bcrypt.GenerateFromPassword(
	// 	[]byte(body.Pass), BcryptUserEncryptionCost)

	// if err != nil {

	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// // Create user secrets
	// var createSecretParams = sqlc.CreateSecretParams{
	// 	Key:   body.ID,
	// 	Scope: "user-login",
	// 	Pass:  string(hashedPass),
	// }
	// if err = db.Queries.CreateSecret(ctx, createSecretParams); err != nil {

	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	c.JSON(http.StatusCreated, nil)
}

func GetUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), GetUserRequestTimeout,
	)

	defer cancel()

	var userID string = c.Param("id")
	if len(strings.TrimSpace(userID)) == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id must not be empty",
		})
		return
	}

	var user *store.User
	var serr *store.StoreError
	if user, serr = ST.GetUser(ctx, userID); serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": serr.Msg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func UpdateUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), UpdateUserRequestTimeout,
	)

	defer cancel()

	var userID string = c.Param("id")
	if len(strings.TrimSpace(userID)) == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id must not be empty",
		})
		return
	}

	var err error
	var body UpdateUserBody
	if err = c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var serr *store.StoreError
	var args = store.UpdateUserParams{
		ID:          option.Some(body.ID),
		Email:       option.Some(body.Email),
		FullName:    option.Some(body.FullName),
		DisplayName: option.Some(body.DisplayName),
	}
	if serr = ST.UpdateUser(ctx, userID, args); serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound,
			errorcodes.UserIDNull,
			errorcodes.UserIDAlreadyExists,
			errorcodes.UserIDFormat,
			errorcodes.UserEmailNull,
			errorcodes.UserEmailAlreadyExists,
			errorcodes.UserEmailFormat,
			errorcodes.UserFullNameFormat,
			errorcodes.UserDisplayNameFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": serr.Msg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func DeleteUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), DeleteUserRequestTimeout,
	)

	defer cancel()

	var userID string = c.Param("id")
	if len(strings.TrimSpace(userID)) == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id must not be empty",
		})
		return
	}

	var serr *store.StoreError
	if serr = ST.DeleteUser(ctx, userID); serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": serr.Msg,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.Status(http.StatusOK)
}

package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"tasks/errorcodes"
	"tasks/services"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUserBody struct {
	ID          string  `json:"id"`
	Pass        string  `json:"pass"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	DisplayName *string `json:"display_name"`
}

type UpdateUserBody struct {
	ID          *string `json:"id"`
	Email       *string `json:"email"`
	FullName    *string `json:"full_name"`
	DisplayName *string `json:"display_name"`
}

const (
	createUserRequestTimeout                = time.Second * 5000000000
	getUserRequestTimeout                   = time.Second * 5000000000
	updateUserRequestTimeout                = time.Second * 5000000000
	deleteUserRequestTimeout                = time.Second * 5000000000
	invalidUserIDFormatHint          string = ""
	invalidUserEmailFormatHint       string = ""
	invalidUserFullNameFormatHint    string = ""
	invalidUserDisplayNameFormatHint string = ""
)

func CreateUser(c *gin.Context) {

	// Setup context
	var ctx, cancel = context.WithTimeout(
		context.Background(), createUserRequestTimeout,
	)

	defer cancel()

	// Parse request body
	var body CreateUserBody
	var err error
	if err = c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code": errorcodes.Internal,
				"msg":  err.Error(),
			},
		})
		return
	}

	var serr *services.ServiceError = services.CreateUser(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		services.CreateUserParams{
			ID:          body.ID,
			Email:       body.Email,
			FullName:    body.FullName,
			DisplayName: body.DisplayName,
		},
	)

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserIDNull:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": errorcodes.UserIDNull,
					"msg":  "user id is null",
				},
			})
			return

		case errorcodes.UserIDAlreadyExists:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": errorcodes.UserIDAlreadyExists,
					"msg":  fmt.Sprintf("user with id '%s' already exists", body.ID),
				},
			})
			return

		case errorcodes.UserEmailAlreadyExists:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": errorcodes.UserEmailAlreadyExists,
					"msg":  fmt.Sprintf("user with email '%s' already exists", body.Email),
				},
			})
			return

		case errorcodes.InvalidUserIDFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": errorcodes.InvalidUserIDFormat,
					"msg":  fmt.Sprintf("user id '%s' format is invalid", body.ID),
					"hint": invalidUserIDFormatHint,
				},
			})
			return

		case errorcodes.InvalidUserEmailFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": errorcodes.InvalidUserEmailFormat,
					"msg":  fmt.Sprintf("user email '%s' format is invalid", body.Email),
					"hint": invalidUserEmailFormatHint,
				},
			})
			return

		case errorcodes.InvalidUserFullNameFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": errorcodes.InvalidUserFullNameFormat,
					"msg":  fmt.Sprintf("user full name '%s' format is invalid", body.FullName),
					"hint": invalidUserFullNameFormatHint,
				},
			})
			return

		case errorcodes.InvalidUserDisplayNameFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code": errorcodes.InvalidUserDisplayNameFormat,
					"msg":  fmt.Sprintf("user display name '%s' format is invalid", *body.DisplayName),
					"hint": invalidUserDisplayNameFormatHint,
				},
			})
			return

		default:

			log.Printf("HANDLER: CreateUser: unexpected error: %v", serr)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code": errorcodes.Internal,
					"msg":  "internal server error",
				},
			})
			return
		}
	}

	c.Status(http.StatusCreated)
}

func GetUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), getUserRequestTimeout,
	)

	defer cancel()

	var userID = c.Param("id")

	var user *services.User
	var serr *services.ServiceError
	user, serr = services.GetUser(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		userID,
	)

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound:

			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"msg": fmt.Sprintf("user with id '%s' not found", userID),
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

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func UpdateUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), updateUserRequestTimeout,
	)

	defer cancel()

	var userID string = c.Param("id")
	if len(strings.TrimSpace(userID)) == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"msg": "user id must not be empty",
			},
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

	var serr *services.ServiceError = services.UpdateUser(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		userID,
		services.UpdateUserParams{
			ID:          body.ID,
			Email:       body.Email,
			FullName:    body.FullName,
			DisplayName: &body.DisplayName,
		},
	)

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserIDNull:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"msg": "user id is null",
				},
			})
			return

		case errorcodes.UserIDAlreadyExists:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"msg": fmt.Sprintf("user with id '%s' already exists", *body.ID),
				},
			})
			return

		case errorcodes.UserEmailAlreadyExists:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"msg": fmt.Sprintf("user with email '%s' already exists", *body.Email),
				},
			})
			return

		case errorcodes.InvalidUserIDFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"msg":  fmt.Sprintf("user id '%s' format is invalid", *body.ID),
					"hint": invalidUserIDFormatHint,
				},
			})
			return

		case errorcodes.InvalidUserEmailFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"msg":  fmt.Sprintf("user email '%s' format is invalid", *body.Email),
					"hint": invalidUserEmailFormatHint,
				},
			})
			return

		case errorcodes.InvalidUserFullNameFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"msg":  fmt.Sprintf("user full name '%s' format is invalid", *body.FullName),
					"hint": invalidUserFullNameFormatHint,
				},
			})
			return

		case errorcodes.InvalidUserDisplayNameFormat:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"msg":  fmt.Sprintf("user display name '%s' format is invalid", *body.DisplayName),
					"hint": invalidUserDisplayNameFormatHint,
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

	c.Status(http.StatusOK)
}

func DeleteUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), deleteUserRequestTimeout,
	)

	defer cancel()

	var userID string = c.Param("id")
	if len(strings.TrimSpace(userID)) == 0 {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"msg": "user id must not be empty",
			},
		})
		return
	}

	var serr *services.ServiceError = services.DeleteUser(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		userID,
	)

	if serr != nil {

		switch serr.Code {
		case errorcodes.UserNotFound:

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"msg": fmt.Sprintf("user with id '%s' not found", userID),
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

	c.Status(http.StatusOK)
}

package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"tasks/db"
	"tasks/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CreateUserRequestTimeout = time.Second * 5
	GetUserRequestTimeout    = time.Second * 5
	UpdateUserRequestTimeout = time.Second * 5
	DeleteUserRequestTimeout = time.Second * 5
)

func CreateUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), CreateUserRequestTimeout,
	)

	defer cancel()

	var createUserParams sqlc.CreateUserParams
	if err := c.BindJSON(&createUserParams); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := db.Queries.CreateUser(ctx, createUserParams); err != nil {

		// var pgerr *pgconn.PgError
		// if errors.As(err, &pgerr) {
		// }

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, nil)
}

func GetUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), GetUserRequestTimeout,
	)

	defer cancel()

	var userId = c.Param("id")
	if userId == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id must not be empty.",
		})

		return
	}

	var user sqlc.User
	var err error
	if user, err = db.Queries.GetUser(ctx, userId); err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Either the user doesn't exist or you don't have access.",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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

	var updateParams sqlc.UpdateUserParams
	if err := c.BindJSON(&updateParams); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	if err := db.Queries.UpdateUser(ctx, updateParams); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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

	var userId = c.Param("id")
	if userId == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id must not be empty.",
		})

		return
	}

	if err := db.Queries.DeleteUser(ctx, userId); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.Status(http.StatusOK)
}

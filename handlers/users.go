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
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

const (
	CreateUserRequestTimeout = time.Second * 5
	GetUserRequestTimeout    = time.Second * 5
	UpdateUserRequestTimeout = time.Second * 5
	DeleteUserRequestTimeout = time.Second * 5
	BcryptUserEncryptionCost = bcrypt.DefaultCost
)

type CreateUserBody struct {
	ID          string `json:"id"`
	Pass        string `json:"pass"`
	FullName    string `json:"full_name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func CreateUser(c *gin.Context) {

	// Setup context
	var ctx, cancel = context.WithTimeout(
		context.Background(), CreateUserRequestTimeout,
	)

	defer cancel()

	// Parse request body
	var body CreateUserBody
	if err := c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create the transaction
	var queries, err = db.Begin(ctx)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer queries.Tx.Rollback(ctx)

	// Create user
	var createUserParams = sqlc.CreateUserParams{
		ID:          body.ID,
		FullName:    body.FullName,
		DisplayName: pgtype.Text{String: body.DisplayName, Valid: true},
		Email:       body.Email,
	}
	if err := queries.Queries.CreateUser(ctx, createUserParams); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	// Encrypt the password
	var hashedPass []byte
	hashedPass, err = bcrypt.GenerateFromPassword(
		[]byte(body.Pass), BcryptUserEncryptionCost)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create user secrets
	var createSecretParams = sqlc.CreateSecretParams{
		Key:   body.ID,
		Scope: "user-login",
		Pass:  string(hashedPass),
	}
	if err := queries.Queries.CreateSecret(ctx, createSecretParams); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Commit the transaction
	if err = queries.Tx.Commit(ctx); err != nil {

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

package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"tasks/db"
	"tasks/sqlc"
	"tasks/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CreateListRequestTimeout = time.Second * 5
	GetListRequestTimeout    = time.Second * 5
	UpdateListRequestTimeout = time.Second * 5
	DeleteListRequestTimeout = time.Second * 5
)

func CreateList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), CreateListRequestTimeout)

	defer cancel()

	var list sqlc.CreateListParams
	var err = c.BindJSON(&list)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var listId int64
	listId, err = db.Queries.CreateList(ctx, list)

	if err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			if pgerr.Code == "25303" && pgerr.ColumnName == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": pgerr.Detail,
				})

				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"listId": listId,
	})
}

func GetList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), GetListRequestTimeout)

	defer cancel()

	var listId, err = utils.ParseInt64(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var list sqlc.List
	list, err = db.Queries.GetList(ctx, listId)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("List with id '%d' doesn't exist.", listId),
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list": gin.H{
			"id":       list.ID,
			"owner_id": list.OwnerID,
			"name":     list.Name,
		},
	})
}

func UpdateList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), UpdateListRequestTimeout)

	defer cancel()

	var listUpdate *sqlc.UpdateListParams
	var err error
	err = c.BindJSON(&listUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = db.Queries.UpdateList(ctx, *listUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List updated successfully",
	})
}

func DeleteList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), DeleteListRequestTimeout)

	defer cancel()

	var listId, err = utils.ParseInt64(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.Queries.DeleteList(ctx, listId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

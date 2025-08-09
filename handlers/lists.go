package handlers

import (
	"context"
	"net/http"
	"tasks/services"
	"tasks/utils"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CreateListRequestTimeout    = time.Second * 5
	GetListRequestTimeout       = time.Second * 5
	UpdateListRequestTimeout    = time.Second * 5
	DeleteListRequestTimeout    = time.Second * 5
	GetListAccessRequestTimeout = time.Second * 5
	SetListAccessRequestTimeout = time.Second * 5
)

type CreateListRequest struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type UpdateListRequest struct {
	Name *string `json:"name"`
}

func CreateList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), CreateListRequestTimeout)

	defer cancel()

	// Parse request body
	var body CreateListRequest
	var err error
	if err = c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Perform creation
	var listID uint64
	var serr *services.ServiceError
	listID, serr = services.CreateList(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		services.CreateListParams{
			Name: body.Name,
		},
	)

	if serr != nil {

		switch serr.Code {
		default:

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"msg": "internal server error",
				},
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"list": gin.H{
			"id": listID,
		},
	})
}

func GetList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), GetListRequestTimeout)

	defer cancel()

	// Parse list id
	var listID, err = utils.ParseUint64(c.Param("id"))
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get list
	var list *services.List
	var serr *services.ServiceError
	list, serr = services.GetList(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		listID,
	)

	if serr != nil {

		switch serr.Code {
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
		"list": list,
	})
}

func UpdateList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), UpdateListRequestTimeout)

	defer cancel()

	// Parse list id
	var listID uint64
	var err error
	if listID, err = utils.ParseUint64(c.Param("id")); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Parse request body
	var body *UpdateListRequest
	err = c.BindJSON(&body)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update the list
	var serr *services.ServiceError = services.UpdateList(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		listID,
		services.UpdateListParams{
			Name: body.Name,
		},
	)

	if serr != nil {

		switch serr.Code {
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

func DeleteList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), DeleteListRequestTimeout)

	defer cancel()

	// Parse list id
	var listID uint64
	var err error
	if listID, err = utils.ParseUint64(c.Param("id")); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Perform deletion
	var serr *services.ServiceError = services.DeleteList(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		listID,
	)

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"msg": "internal server error",
			},
		})
		return
	}

	c.Status(http.StatusOK)
}

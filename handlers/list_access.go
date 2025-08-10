package handlers

import (
	"context"
	"net/http"
	"tasks/services"
	"tasks/utils"

	"github.com/gin-gonic/gin"
)

type GetListAccessRequest struct {
	UserID string
}

type AddListAccessRequest struct {
	UserID string
	Access []services.ListAccessType
}

type RemoveListAccessRequest struct {
	UserID string
	Access []services.ListAccessType
}

func GetListAccess(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), GetListAccessRequestTimeout)

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
	var body GetListAccessRequest
	if err = c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get list access
	var listAccess *services.ListAccess
	var serr *services.ServiceError
	listAccess, serr = services.GetListAccess(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		services.GetListAccessParams{
			UserID: body.UserID,
			ListID: listID,
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

	c.JSON(http.StatusOK, gin.H{
		"access": gin.H{
			"user_id":  listAccess.UserID,
			"list_id":  listAccess.ListID,
			"accesses": listAccess.Access,
		},
	})
}

func AddListAccess(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), GetListAccessRequestTimeout)

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
	var body AddListAccessRequest
	if err = c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Add list access
	var serr *services.ServiceError = services.AddListAccess(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		services.AddListAccessParams{
			UserID: body.UserID,
			ListID: listID,
			Access: body.Access,
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

func RemoveListAccess(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), GetListAccessRequestTimeout)

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
	var body RemoveListAccessRequest
	if err = c.BindJSON(&body); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Add list access
	var serr *services.ServiceError = services.RemoveListAccess(
		services.ServiceContext{
			Ctx:    ctx,
			UserID: c.GetString("userid"),
		},
		services.RemoveListAccessParams{
			UserID: &body.UserID,
			ListID: &listID,
			Access: &body.Access,
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

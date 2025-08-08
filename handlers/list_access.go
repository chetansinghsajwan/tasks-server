package handlers

import (
	"context"
	"net/http"
	"tasks/option"
	"tasks/store"
	"tasks/utils"

	"github.com/gin-gonic/gin"
)

type GetListAccessRequest struct {
	UserID string
}

type AddListAccessRequest struct {
	UserID string
	Access []store.ListAccessType
}

type RemoveListAccessRequest struct {
	UserID string
	Access []store.ListAccessType
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
	var listAccess *store.ListAccess
	var serr *store.StoreError
	listAccess, serr = ST.GetListAccess(ctx, store.GetListAccessParams{
		UserID: body.UserID,
		ListID: listID,
	})

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access": *listAccess,
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
	var serr *store.StoreError = ST.AddListAccess(ctx, store.AddListAccessParams{
		UserID: body.UserID,
		ListID: listID,
		Access: body.Access,
	})

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
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
	var serr *store.StoreError = ST.RemoveListAccess(ctx, store.RemoveListAccessParams{
		UserID: option.Some(body.UserID),
		ListID: option.Some(listID),
		Access: option.Some(body.Access),
	})

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.Status(http.StatusOK)
}

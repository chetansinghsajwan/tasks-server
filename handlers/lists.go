package handlers

import (
	"context"
	"net/http"
	"tasks/option"
	"tasks/store"
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
	ID   store.ListID `json:"id"`
	Name string       `json:"name"`
}

type UpdateListRequest struct {
	Name option.Option[string] `json:"name"`
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
	var listID store.ListID
	var serr *store.StoreError
	listID, serr = ST.CreateList(ctx, store.CreateListParams{
		Name: body.Name,
	})

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
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
	var listID, err = store.ParseListID(c.Param("id"))
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get list
	var list *store.List
	var serr *store.StoreError
	list, serr = ST.GetList(ctx, listID)

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
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
	var listID store.ListID
	var err error
	if listID, err = store.ParseListID(c.Param("id")); err != nil {

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
	var serr *store.StoreError = ST.UpdateList(ctx, listID, store.UpdateListParams{
		Name: body.Name,
	})

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.Status(http.StatusOK)
}

func DeleteList(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), DeleteListRequestTimeout)

	defer cancel()

	// Parse list id
	var listID store.ListID
	var err error
	if listID, err = store.ParseListID(c.Param("id")); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Perform deletion
	var serr *store.StoreError = ST.DeleteList(ctx, listID)
	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.Status(http.StatusOK)
}

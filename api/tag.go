package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/db/Result"
)
type TagResponseFormat struct{
	TagId int32 `json:"tag_id"`
	Name string `json:"name"`
}

func TagResponse(tag Database.Tag)TagResponseFormat{
	return TagResponseFormat{
		TagId: tag.TagID,
		Name: tag.Name,
	}
}

func formatManytags(tags []Database.Tag) []TagResponseFormat {

	var formattedtags []TagResponseFormat

	for _, tag := range tags {
		formattedtags = append(formattedtags, TagResponse(tag))
	}

	return formattedtags
}

type CreateTagsRequest struct {
	Owner string `json:"owner" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

func (server *Server) CreateTags(ctx *gin.Context) {
	var req CreateTagsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	arg := Database.CreateTagsParams{
		Owner: sql.NullString{String: req.Owner, Valid: true},
		Name: req.Name,
	}

	tag, err := server.store.CreateTags(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, TagResponse(tag))
}

type GetTagRequest struct {
	TagId int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetTag(ctx *gin.Context) {
	var req GetTagRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	tag, err := server.store.GetTag(ctx, req.TagId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

type ListTagsRequest struct {
	TagId    int32 `form:"tag_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) ListTags(ctx *gin.Context) {
	var req ListTagsRequest

	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := Database.ListTagsParams{
		TagID: req.TagId,
		Limit: req.PageSize,
	}
	tags, err := server.store.ListTags(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, formatManytags(tags))
}
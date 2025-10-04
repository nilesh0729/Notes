package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/db/Result"
)

type CreateTagsRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) CreateTags(ctx *gin.Context) {
	var req CreateTagsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	tag, err := server.store.CreateTags(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

type GetTagRequest struct {
	TagId int32 `uri:"tag_id" binding:"required"`
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
	PageSize int32 `form:"PageSize" binding:"required,min=5,max=20"`
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
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tags)
}
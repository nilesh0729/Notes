package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/db/Result"
)

type AddTagToNoteRequest struct {
	NoteId int32 `json:"note_id" binding:"required"`
	TagId  int32 `json:"tag_id" binding:"required"`
}

func (server *Server) AddTagToNote(ctx *gin.Context) {
	var req AddTagToNoteRequest

	err := ctx.ShouldBindJSON(&req)
	if err!=nil{
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}


	arg := Database.AddTagToNoteParams{
		NoteID: req.NoteId,
		TagID: req.TagId,
	}
	notetag, err := server.store.AddTagToNote(ctx, arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, notetag)

}



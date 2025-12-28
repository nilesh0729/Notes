package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/internal/db/Result"
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


type ListNotesForTagRequest struct {
	TagID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) ListNotesForTag(ctx *gin.Context) {
	var req ListNotesForTagRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	notes, err := server.store.GetNotesForTag(ctx, req.TagID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, notes)
}

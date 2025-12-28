package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/internal/db/Result"
	"github.com/nilesh0729/Notes/internal/tokens"
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


	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*tokens.Payload)
	
	// Check Note Ownership
	note, err := server.store.GetNoteById(ctx, req.NoteId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if note.Owner.String != authPayload.Username {
		err := errors.New("note doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	
	// Check Tag Ownership
	tag, err := server.store.GetTag(ctx, req.TagId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if tag.Owner.String != authPayload.Username {
		err := errors.New("tag doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
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

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*tokens.Payload)
	
	arg := Database.GetNotesForTagParams{
		TagID: req.TagID,
		Owner: sql.NullString{String: authPayload.Username, Valid: true},
	}
	notes, err := server.store.GetNotesForTag(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	var dbNotes []Database.Note
	for _, rawNote := range notes {
		dbNotes = append(dbNotes, Database.Note{
			NoteID:    rawNote.NoteID,
			Title:     rawNote.Title,
			Owner:     rawNote.Owner,
			Content:   rawNote.Content,
			Pinned:    rawNote.Pinned,
			Archived:  rawNote.Archived,
			CreatedAt: rawNote.CreatedAt,
			UpdatedAt: rawNote.UpdatedAt,
		})
	}

	ctx.JSON(http.StatusOK, server.formatManyNotes(ctx, dbNotes))
}

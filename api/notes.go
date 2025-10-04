package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/db/Result"
)

type ResponseFormat struct {
	NoteId    int32     `json:"note_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ResponseFormating(note Database.Note) ResponseFormat {
	return ResponseFormat{
		NoteId:    note.NoteID,
		Title:     note.Title.String,
		Content:   note.Content.String,
		CreatedAt: note.CreatedAt.Time,
		UpdatedAt: note.CreatedAt.Time,
	}
}

func formatManyNotes(notes []Database.Note) []ResponseFormat {

	var formattedNotes []ResponseFormat

	for _, note := range notes {
		formattedNotes = append(formattedNotes, ResponseFormating(note))
	}

	return formattedNotes
}

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (server *Server) CreateNote(ctx *gin.Context) {
	var req CreateNoteRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	arg := Database.CreateNoteParams{
		Title:   sql.NullString{String: req.Title, Valid: true},
		Content: sql.NullString{String: req.Content, Valid: true},
	}
	note, err := server.store.CreateNote(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ResponseFormating(note))
}

type GetNoteByIdRequest struct {
	NoteID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetNoteById(ctx *gin.Context) {
	var req GetNoteByIdRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	note, err := server.store.GetNoteById(ctx, (req.NoteID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ResponseFormating(note))

}

type ListNotesRequest struct {
	cursor   int32 `form:"cursor" binding:"required,min=0"`
	PageSize int32 `form:"page_size" binding:"required,max=20,min=5"`
}

func (server *Server) ListNotes(ctx *gin.Context) {
	var req ListNotesRequest

	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := Database.ListNotesParams{
		NoteID: req.cursor,
		Limit:  req.PageSize,
	}
	notes, err := server.store.ListNotes(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, formatManyNotes(notes))
}

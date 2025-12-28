package api
// Force rebuild


import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/internal/db/Result"
	"github.com/nilesh0729/Notes/internal/tokens"
)

type ResponseFormat struct {
	NoteId    int32     `json:"note_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Pinned    bool      `json:"pinned"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []TagResponseFormat `json:"tags"`
}

func ResponseFormating(note Database.Note, tags []TagResponseFormat) ResponseFormat {
	return ResponseFormat{
		NoteId:    note.NoteID,
		Title:     note.Title.String,
		Content:   note.Content.String,
		Pinned:    note.Pinned.Bool,
		Archived:  note.Archived.Bool,
		CreatedAt: note.CreatedAt.Time,
		UpdatedAt: note.UpdatedAt.Time,
		Tags:      tags,
	}
}

func transformTagRows(rows []Database.GetTagsForNoteRow) []TagResponseFormat {
	var tags []TagResponseFormat
	for _, row := range rows {
		tags = append(tags, TagResponseFormat{
			TagId: row.TagID,
			Name:  row.Name,
		})
	}
	return tags
}

func (server *Server) formatManyNotes(ctx *gin.Context, notes []Database.Note) []ResponseFormat {
	var formattedNotes []ResponseFormat

	for _, note := range notes {
		tags, _ := server.store.GetTagsForNote(ctx, note.NoteID)
		formattedNotes = append(formattedNotes, ResponseFormating(note, transformTagRows(tags)))
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
	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*tokens.Payload)

	arg := Database.CreateNoteParams{
		Owner:   sql.NullString{String: authPayload.Username, Valid: true},
		Title:   sql.NullString{String: req.Title, Valid: true},
		Content: sql.NullString{String: req.Content, Valid: true},
	}
	note, err := server.store.CreateNote(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ResponseFormating(note, []TagResponseFormat{}))
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
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	tags, _ := server.store.GetTagsForNote(ctx, req.NoteID)
	ctx.JSON(http.StatusOK, ResponseFormating(note, transformTagRows(tags)))

}

type ListNotesRequest struct {
	Cursor   int32  `form:"cursor"`
	PageSize int32  `form:"page_size" binding:"required,max=100,min=5"`
	Search   string `form:"search"`
}

func (server *Server) ListNotes(ctx *gin.Context) {
	var req ListNotesRequest

	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var notes []Database.Note
	
	if req.Search != "" {
		// Use SearchNotes query
		arg := Database.SearchNotesParams{
			Column1: sql.NullString{String: req.Search, Valid: true},
			Limit:   req.PageSize,
			Offset:  req.Cursor, // Cursor acts as Offset for search
		}
		notes, err = server.store.SearchNotes(ctx, arg)
	} else {
		// Use standard ListNotes query
		arg := Database.ListNotesParams{
			NoteID: req.Cursor,
			Limit:  req.PageSize,
		}
		notes, err = server.store.ListNotes(ctx, arg)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, server.formatManyNotes(ctx, notes))
}

type UpdateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (server *Server) UpdateNote(ctx *gin.Context) {
	var req UpdateNoteRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	noteIdStr := ctx.Param("id")
	var noteId int32
	// simple atoi
	fmt.Sscanf(noteIdStr, "%d", &noteId)
	if noteId == 0 {
		ctx.JSON(http.StatusBadRequest, errResponse(fmt.Errorf("invalid note id")))
		return
	}

	arg := Database.UpdateNoteParams{
		NoteID:  noteId,
		Title:   sql.NullString{String: req.Title, Valid: true},
		Content: sql.NullString{String: req.Content, Valid: true},
	}

	note, err := server.store.UpdateNote(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	tags, _ := server.store.GetTagsForNote(ctx, note.NoteID)
	ctx.JSON(http.StatusOK, ResponseFormating(note, transformTagRows(tags)))
}

func (server *Server) DeleteNote(ctx *gin.Context) {
	noteIdStr := ctx.Param("id")
	var noteId int32
	fmt.Sscanf(noteIdStr, "%d", &noteId)
	
	// Manually cascade delete
	err := server.store.DeleteNoteTagsByNoteId(ctx, noteId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = server.store.DeleteNote(ctx, noteId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}

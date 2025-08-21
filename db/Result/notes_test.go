package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomNotes(t *testing.T) Note {
	User := CreateRandomUser(t)
	arg := CreateNotesParams{
		UserID:   sql.NullInt32{Int32: User.ID, Valid: true},
		Title:    sql.NullString{String: util.RandomString(4), Valid: true},
		Content:  sql.NullString{String: util.RandomString(10), Valid: true},
		Archived: sql.NullBool{Bool: false, Valid: true},
		Pinned:   sql.NullBool{Bool: false, Valid: true},
	}
	note, err := testQueries.CreateNotes(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, note)

	require.Equal(t, note.UserID.Int32, User.ID)
	require.Equal(t, note.Title, arg.Title)
	require.Equal(t, note.Content, arg.Content)
	require.Equal(t, note.Archived, arg.Archived)
	require.Equal(t, note.Pinned, arg.Pinned)

	require.NotZero(t, note.ID)

	require.WithinDuration(t, note.CreatedAt.Time, User.CreatedAt.Time, time.Second)

	return note
}

func TestCreateNotes(t *testing.T) {
	CreateRandomNotes(t)
}

func TestGetNotesById(t *testing.T) {
	note := CreateRandomNotes(t)
	note2, err := testQueries.GetNoteById(context.Background(), note.ID)

	require.NoError(t, err)
	require.NotEmpty(t, note2)

	require.Equal(t, note2.ID, note.ID)
	require.Equal(t, note2.Title, note.Title)
	require.Equal(t, note2.Content, note.Content)
	require.Equal(t, note2.Archived, note.Archived)
	require.Equal(t, note2.Pinned, note.Pinned)
	require.Equal(t, note2.UserID, note.UserID)

	require.Equal(t, note2.CreatedAt, note.CreatedAt)
	require.Equal(t, note2.UpdatedAt, note.UpdatedAt)
}

func TestListUserNotes(t *testing.T) {
	User := CreateRandomUser(t)

	for i := 0; i < 10; i++ {
		testQueries.CreateNotes(context.Background(), CreateNotesParams{
			UserID:   sql.NullInt32{Int32: User.ID, Valid: true},
			Title:    sql.NullString{String: util.RandomString(4), Valid: true},
			Content:  sql.NullString{String: util.RandomString(10), Valid: true},
			Archived: sql.NullBool{Bool: false, Valid: true},
			Pinned:   sql.NullBool{Bool: false, Valid: true},
		})
	}
	notes, err := testQueries.ListUserNotes(context.Background(), util.ToNullInt32(User.ID))
	require.NoError(t, err)

	require.Len(t, notes, 10)

	for _, note := range notes {
		require.NotEmpty(t, note)
	}

}

func TestUpdateNotes(t *testing.T){
	note1 := CreateRandomNotes(t)

	arg := UpdateNotesParams{
		ID: note1.ID,
		Title: note1.Title,
		Content: sql.NullString{String: util.RandomString(8), Valid: true},
		Archived: sql.NullBool{Bool: false, Valid: true},
		Pinned: sql.NullBool{Bool: false, Valid: true},
	}

	note2, err := testQueries.UpdateNotes(context.Background(),arg )

	require.NoError(t, err)
	require.NotEmpty(t, note2)

	require.Equal(t, note1.ID, note2.ID)
	require.Equal(t, note1.Title, note2.Title)
	require.Equal(t, note2.Content, arg.Content)
	require.Equal(t, note2.Archived, arg.Archived)
	require.Equal(t, note2.Pinned, arg.Pinned)

	require.Equal(t,note2.CreatedAt, note1.CreatedAt)

	require.WithinDuration(t, note2.UpdatedAt.Time, time.Now(), time.Second)
}

func TestDeleteNotes(t *testing.T){
	note1 := CreateRandomNotes(t)

	err := testQueries.DeleteNotes(context.Background(), note1.ID)
	require.NoError(t, err)

	note, err := testQueries.GetNoteById(context.Background(), note1.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, note)
}

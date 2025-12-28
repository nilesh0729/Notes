package Database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nilesh0729/Notes/internal/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomNote(t *testing.T) Note {
	user := RandomUser(t)
	arg := CreateNoteParams{
		Owner:   sql.NullString{String:user.Username , Valid: true},
		Title:   sql.NullString{String: util.RandomString(6), Valid: true},
		Content: sql.NullString{String: util.RandomString(8), Valid: true},
	}
	note, err := testQueries.CreateNote(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, note)

	require.Equal(t, note.Owner, arg.Owner)
	require.Equal(t, note.Title, arg.Title)
	require.Equal(t, note.Content, arg.Content)

	require.NotZero(t, note.NoteID)

	require.False(t, note.Pinned.Bool)
	require.False(t, note.Archived.Bool)

	return note

}

func TestCreateNote(t *testing.T) {
	CreateRandomNote(t)
}

func TestGetNoteById(t *testing.T) {
	note1 := CreateRandomNote(t)
	note, err := testQueries.GetNoteById(context.Background(), note1.NoteID)
	require.NoError(t, err)
	require.NotEmpty(t, note)

	require.Equal(t, note.Owner, note1.Owner)
	require.Equal(t, note.NoteID, note1.NoteID)
	require.Equal(t, note.Title, note1.Title)
	require.Equal(t, note.Content, note1.Content)
	require.Equal(t, note.Archived, note1.Archived)
	require.Equal(t, note.Pinned, note1.Pinned)

	require.WithinDuration(t, note.CreatedAt.Time, note1.CreatedAt.Time, time.Second)
	require.WithinDuration(t, note.UpdatedAt.Time, note1.UpdatedAt.Time, time.Second)
}

func TestListNotes(t *testing.T) {
	var Cnotes []Note

	for i := 0; i < 10; i++ {
		note := CreateRandomNote(t)
		Cnotes = append(Cnotes, note)
	}
	StartsAfterId := Cnotes[0].NoteID - 1

	arg := ListNotesParams{
		NoteID: StartsAfterId,
		Limit:  5,
	}

	Notes, err := testQueries.ListNotes(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, Notes)

	require.Len(t, Notes, 5)

	for i, note := range Notes {
		require.NotEmpty(t, note)

		expected := Cnotes[i]

		require.Equal(t, expected.Owner, note.Owner)
		require.Equal(t, expected.Title, note.Title)
		require.Equal(t, expected.Content, expected.Content)
		require.Equal(t, expected.Archived, expected.Archived)
		require.Equal(t, expected.Pinned, expected.Pinned)
	}
}

func TestUpdateNote(t *testing.T) {
	note1 := CreateRandomNote(t)
	arg := UpdateNoteParams{
		NoteID:  note1.NoteID,
		Title:   sql.NullString{String: util.RandomString(6), Valid: true},
		Content: sql.NullString{String: util.RandomString(10), Valid: true},
	}
	note2, err := testQueries.UpdateNote(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, note2)

	require.Equal(t, note2.Owner, note1.Owner)
	require.Equal(t, note2.NoteID, note1.NoteID)
	require.Equal(t, note2.Title, arg.Title)
	require.Equal(t, note2.Content, arg.Content)

	require.Equal(t, note2.Archived, note1.Archived)
	require.Equal(t, note2.Pinned, note1.Pinned)

	require.WithinDuration(t, note1.CreatedAt.Time, note2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), note2.UpdatedAt.Time, time.Second)

	require.True(t, note2.UpdatedAt.Time.After(note1.CreatedAt.Time))
}

func TestDeleteNote(t *testing.T) {
	note1 := CreateRandomNote(t)
	err := testQueries.DeleteNote(context.Background(), note1.NoteID)
	require.NoError(t, err)

	note2, err := testQueries.GetNoteById(context.Background(), note1.NoteID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, note2)
}

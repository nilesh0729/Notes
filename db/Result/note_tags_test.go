package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

func AddRandomTagToNote(t *testing.T) NoteTag{
	note := CreateRandomNotes(t)
	tag := CreateRandomTags(t)
	arg := AddTagToNoteParams{
		NoteID: note.ID,
		TagID: tag.ID,
	}
	noteTag, err := testQueries.AddTagToNote(context.Background(),arg)
	require.NoError(t, err)
	require.NotEmpty(t, noteTag)

	require.Equal(t, noteTag.NoteID, note.ID)
	require.Equal(t, noteTag.TagID, tag.ID)

	return noteTag
}

func TestGetNotesForTags(t *testing.T){
	noteTag := AddRandomTagToNote(t)
	
	notes, err := testQueries.GetNotesForTag(context.Background(), noteTag.TagID)
	require.NoError(t, err)
	require.NotEmpty(t, notes)

	for _, note := range notes{
		require.NotEmpty(t, note)
	}
}

func TestGetTagsForNote(t *testing.T){
	noteTag := AddRandomTagToNote(t)
	tags, err:= testQueries.GetTagsForNote(context.Background(), noteTag.NoteID)
	require.NoError(t,err)

	for _, tag := range tags {
		require.NotEmpty(t, tag)
	}

}

func TestRemoveTagFromNote(t *testing.T){
	user := CreateRandomUser(t)

	Note , err := testQueries.CreateNotes(context.Background(), CreateNotesParams{
		UserID: sql.NullInt32{Int32: user.ID,Valid: true},
		Title: sql.NullString{String: util.RandomString(4), Valid: true},
		Content: sql.NullString{String: util.RandomString(10), Valid: true},
		Archived: sql.NullBool{Bool: false, Valid: true},
		Pinned: sql.NullBool{Bool: false, Valid: true},
	})
	require.NoError(t,err)


	tag, err := testQueries.CreateTags(context.Background(), CreateTagsParams{
		Name: util.RandomString(6),
		UserID: user.ID,
	})
	require.NoError(t, err)

	noteTag ,err:= testQueries.AddTagToNote(context.Background(),AddTagToNoteParams{
		NoteID: Note.ID,
		TagID: tag.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, noteTag)

	require.Equal(t, noteTag.NoteID, Note.ID)
	require.Equal(t, noteTag.TagID, tag.ID)


	arg := RemoveTagFromNoteParams{
		NoteID: Note.ID,
		TagID: tag.ID,
	}
	hell := testQueries.RemoveTagFromNote(context.Background(),arg)
	require.NoError(t, hell)

	tags , err := testQueries.GetTagsForNote(context.Background(), Note.ID)
	require.NoError(t,err)
	require.Empty(t, tags)
}

func TestNoteTag(t *testing.T){
	noteTag := AddRandomTagToNote(t)
	err := testQueries.DeleteNoteTags(context.Background(),noteTag.NoteID)
	require.NoError(t, err)

	notes, err := testQueries.GetNotesForTag(context.Background(),noteTag.TagID)
	require.NoError(t,err)
	require.Empty(t, notes)
}
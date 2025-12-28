package Database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomNoteTag(t *testing.T, note Note, tag Tag) NoteTag {
	arg := AddTagToNoteParams{
		NoteID: note.NoteID,
		TagID:  tag.TagID,
	}

	noteTag, err := testQueries.AddTagToNote(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, noteTag.NoteID, note.NoteID)
	require.Equal(t, noteTag.TagID, tag.TagID)

	return noteTag
}

func TestAddTagToNote(t *testing.T) {
	note := CreateRandomNote(t)
	tag := CreateRandomTags(t)
	CreateRandomNoteTag(t, note, tag)
}

func TestGetNotesForTag(t *testing.T) {
	note1 := CreateRandomNote(t)
	note2 := CreateRandomNote(t)
	tag1 := CreateRandomTags(t)

	_ = CreateRandomNote(t)

	noteTag1 := CreateRandomNoteTag(t, note1, tag1)
	require.NotEmpty(t, noteTag1)

	noteTag2 := CreateRandomNoteTag(t, note2, tag1)
	require.NotEmpty(t, noteTag2)

	notes, err := testQueries.GetNotesForTag(context.Background(), tag1.TagID)
	require.NoError(t, err)
	require.NotEmpty(t, notes)

	require.Len(t, notes, 2)

	expected := map[int64]bool{
		int64(note1.NoteID): true,
		int64(note2.NoteID): true,
	}

	retrievedNoteIDs := map[int64]bool{}
	for _, note := range notes {
		retrievedNoteIDs[int64(note.NoteID)] = true
	}

	require.Equal(t, expected, retrievedNoteIDs)
}

func TestGetTagsForNote(t *testing.T) {
	note1 := CreateRandomNote(t)
	tag1 := CreateRandomTags(t)
	tag2 := CreateRandomTags(t)

	_ = CreateRandomTags(t)

	notetag1 := CreateRandomNoteTag(t, note1, tag1)
	require.NotEmpty(t, notetag1)

	notetag2 := CreateRandomNoteTag(t, note1, tag2)
	require.NotEmpty(t, notetag2)

	tags, err := testQueries.GetTagsForNote(context.Background(), note1.NoteID)
	require.NoError(t, err)
	require.NotEmpty(t, tags)

	require.Len(t, tags, 2)

	expected := map[int64]bool{
		int64(tag1.TagID): true,
		int64(tag2.TagID): true,
	}

	retrievedTagIDs := map[int64]bool{}
	for _, tag := range tags {
		retrievedTagIDs[int64(tag.TagID)] = true
	}

	require.Equal(t, expected, retrievedTagIDs)

}

func TestRemoveTagFromNote(t *testing.T) {
	note1 := CreateRandomNote(t)
	tag1 := CreateRandomTags(t)

	notetag := CreateRandomNoteTag(t, note1, tag1)
	require.NotEmpty(t, notetag)

	arg := RemoveTagFromNoteParams{
		NoteID: note1.NoteID,
		TagID:  tag1.TagID,
	}

	err := testQueries.RemoveTagFromNote(context.Background(), arg)
	require.NoError(t, err)

	note, err := testQueries.GetNotesForTag(context.Background(), tag1.TagID)
	require.NoError(t, err)
	require.Empty(t, note)

	tag, err := testQueries.GetTagsForNote(context.Background(), note1.NoteID)
	require.NoError(t, err)
	require.Empty(t, tag)
}

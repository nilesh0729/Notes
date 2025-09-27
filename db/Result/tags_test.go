package Database

import (
	"context"
	"database/sql"
	"testing"

	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTags(t *testing.T) Tag {
	name := util.RandomString(6)
	tag, err := testQueries.CreateTags(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, tag)

	require.Equal(t, tag.Name, name)
	require.NotZero(t, tag.TagID)

	return  tag

}

func TestCreateTags(t *testing.T){
	CreateRandomTags(t)
}

func TestGetTag(t *testing.T){
	tag1 := CreateRandomTags(t)
	tag,err := testQueries.GetTag(context.Background(),tag1.TagID)
	require.NoError(t, err)
	require.NotEmpty(t, tag)
	
	require.Equal(t, tag.TagID, tag.TagID)
	require.Equal(t, tag.Name, tag.Name)
}

func TestListTags(t *testing.T){
	var CreatedTags []Tag

	for i:=0; i<10; i++{
		tag := CreateRandomTags(t)
		CreatedTags = append(CreatedTags, tag)
	}

	TagStartingPoint := CreatedTags[0].TagID - 1
	arg := ListTagsParams{
		TagID: TagStartingPoint,
		Limit: 5,
	}

	tags, err := testQueries.ListTags(context.Background(),arg)
	require.NoError(t, err)
	require.NotEmpty(t, tags)
	require.Len(t,tags, 5)

	for i, tag := range tags{
		require.NotEmpty(t, tag)

		expected := CreatedTags[i]

		require.Equal(t, expected.Name, tag.Name)
		require.Equal(t, expected.TagID, tag.TagID)
	}
}

func TestUpdateTag(t *testing.T){
	tag1 := CreateRandomTags(t)

	arg := UpdateTagParams{
		TagID: tag1.TagID,
		Name: util.RandomString(6),
	}
	tag2, err := testQueries.UpdateTag(context.Background(),arg)
	require.NoError(t, err)
	require.NotEmpty(t, tag2)

	require.Equal(t, tag1.TagID, tag2.TagID)
	require.Equal(t, tag2.Name, arg.Name)
}

func TestDeleteTag(t *testing.T){
	tag1 := CreateRandomTags(t)

	err := testQueries.DeleteTag(context.Background(),tag1.TagID)
	require.NoError(t, err)

	tag2, err := testQueries.GetTag(context.Background(), tag1.TagID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tag2)

}
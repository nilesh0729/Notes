package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTags(t *testing.T) Tag{
	User := CreateRandomUser(t)
	arg := CreateTagsParams{
		Name: util.RandomString(4),
		UserID: User.ID,
	}
	tag, err := testQueries.CreateTags(context.Background(),arg)
	require.NoError(t, err)
	require.NotEmpty(t, tag)

	require.Equal(t, tag.UserID, User.ID)
	require.Equal(t, tag.Name, arg.Name)

	require.NotZero(t, tag.ID)

	return tag
}

func TestCreateTags(t *testing.T){
	CreateRandomTags(t)
}

func TestGetUserTag(t *testing.T){
	tag1 := CreateRandomTags(t)

	arg := GetUserTagParams{
		UserID: tag1.UserID,
		Name: tag1.Name,
	}
	tag2, err := testQueries.GetUserTag(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tag2)

	require.Equal(t,tag2.ID, tag1.ID)
	require.Equal(t, tag2.Name, tag1.Name)
	require.Equal(t, tag2.UserID, tag1.UserID)
	
}

func TestListUserTags(t *testing.T){
	User := CreateRandomUser(t)
	for i:= 0; i<10; i++ {
		testQueries.CreateTags(context.Background(), CreateTagsParams{
			Name: util.RandomString(6),
			UserID: User.ID,
		})
	}
	
	tags, err := testQueries.ListUserTags(context.Background(), User.ID )
	require.NoError(t, err)
	require.Len(t, tags, 10)

	for _, tag := range tags{
		require.NotEmpty(t, tag)
	}

}

func TestRenameTag(t *testing.T){
	tag1 := CreateRandomTags(t)

	arg := RenameTagParams{
		ID: tag1.ID,
		UserID: tag1.UserID,
		Name: util.RandomString(4),
	}
	tag2, err := testQueries.RenameTag(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tag2)

	require.Equal(t, tag2.ID, tag1.ID)
	require.Equal(t, tag2.UserID, tag1.UserID)
	require.Equal(t, tag2.Name, arg.Name)

}

func TestDeletetag(t *testing.T){
	tag1 := CreateRandomTags(t)

	arg := DeleteTagsParams{
		ID: tag1.ID,
		UserID: tag1.UserID,
	}
	err := testQueries.DeleteTags(context.Background(), arg)
	require.NoError(t, err)

	hell := GetUserTagParams{
		UserID: tag1.UserID,
		Name: tag1.Name,
	}

	tag2, err := testQueries.GetUserTag(context.Background(), hell)
	require.EqualError(t, err ,sql.ErrNoRows.Error())
	require.Empty(t, tag2)
}


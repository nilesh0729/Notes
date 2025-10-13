package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	Database "github.com/nilesh0729/Notes/db/Result"
	MockDB "github.com/nilesh0729/Notes/db/mock"
	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTag(t *testing.T) {
	tag := RandomTag()

	testcases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *MockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"name": tag.Name,
			},
			buildStubs: func(store *MockDB.MockStore) {
				store.EXPECT().
					CreateTags(gomock.Any(), gomock.Eq(tag.Name)).
					Times(1).
					Return(tag, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				TagResponseMatching(t, recorder.Body, tag)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"name": tag.TagID,
			},
			buildStubs: func(store *MockDB.MockStore) {
				store.EXPECT().
					CreateTags(gomock.Any(), gomock.Eq(tag.Name)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"name": tag.Name,
			},
			buildStubs: func(store *MockDB.MockStore) {
				store.EXPECT().
					CreateTags(gomock.Any(), gomock.Eq(tag.Name)).
					Times(1).
					Return(Database.Tag{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := MockDB.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/tags"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			newBody := bytes.NewBuffer(data)

			request, err := http.NewRequest(http.MethodPost, url, newBody)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetTag(t *testing.T) {
	tag := RandomTag()

	testcases := []struct {
		name          string
		tagId         int32
		buildstubs    func(store *MockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			tagId: tag.TagID,
			buildstubs: func(store *MockDB.MockStore) {
				store.EXPECT().
					GetTag(gomock.Any(), gomock.Eq(tag.TagID)).
					Times(1).
					Return(tag, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				TagResponseMatching(t, recorder.Body, tag)
			},
		},
		{
			name:  "BadRequest",
			tagId: 0,
			buildstubs: func(store *MockDB.MockStore) {
				store.EXPECT().
					GetTag(gomock.Any(), gomock.Eq(tag.TagID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "InternalServerError",
			tagId: tag.TagID,
			buildstubs: func(store *MockDB.MockStore) {
				store.EXPECT().
					GetTag(gomock.Any(), gomock.Any()).
					Times(1).
					Return(Database.Tag{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := MockDB.NewMockStore(ctrl)
			tc.buildstubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tags/%d", tc.tagId)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListTags(t *testing.T) {
	n := 5

	tags := make([]Database.Tag, n)

	for i := 0; i < n; i++ {
		tags[i] = RandomTag()
	}
	type Query struct {
		TagId    int32
		PageSize int32
	}
	testcases := []struct {
		name          string
		query         Query
		buildstubs    func(store *MockDB.MockStore, query Query)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				TagId:    1,
				PageSize: int32(n),
			},
			buildstubs: func(store *MockDB.MockStore, query Query) {
				arg := Database.ListTagsParams{
					TagID: query.TagId,
					Limit: query.PageSize,
				}
				store.EXPECT().
					ListTags(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(tags, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				TagsResponseMatching(t, recorder.Body, tags)
			},
		},
		{
			name: "BadRequest",
			query: Query{
				TagId: 0,
				PageSize: int32(0),
			},
			buildstubs: func(store *MockDB.MockStore, query Query) {
				arg := Database.ListTagsParams{
					TagID: 0,
					Limit: 2,
				}
				store.EXPECT().
				    ListTags(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			query: Query{
				TagId: 1,
				PageSize: int32(n),
			},
			buildstubs: func(store *MockDB.MockStore, query Query) {
				arg := Database.ListTagsParams{
					TagID: query.TagId,
					Limit: int32(n),
				}
				store.EXPECT().
				    ListTags(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]Database.Tag{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := MockDB.NewMockStore(ctrl)
			tc.buildstubs(store, tc.query)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tags?tag_id=%d&page_size=%d", tc.query.TagId, tc.query.PageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}
}

func TagResponseMatching(t *testing.T, body *bytes.Buffer, tag Database.Tag) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var GotTag Database.Tag

	err = json.Unmarshal(data, &GotTag)

	require.Equal(t, tag, GotTag)
}
func TagsResponseMatching(t *testing.T, body *bytes.Buffer, tags []Database.Tag) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var GotTags []Database.Tag

	err = json.Unmarshal(data, &GotTags)
	require.NoError(t, err)

	require.Equal(t, GotTags, tags)
}

func RandomTag() Database.Tag {
	return Database.Tag{
		TagID: int32(util.RandomInt(1, 10)),
		Name:  util.RandomString(4),
	}
}

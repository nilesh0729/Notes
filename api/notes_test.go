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
	mockDB "github.com/nilesh0729/Notes/db/Mock"
	Database "github.com/nilesh0729/Notes/db/Result"

	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

func TestCreateNoteApi(t *testing.T) {
	note := RandomNotes()
	arg := Database.CreateNoteParams{
		Owner:   note.Owner,
		Title:   note.Title,
		Content: note.Content,
	}

	testcases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"owner":   note.Owner.String,
				"title":   note.Title.String,
				"content": note.Content.String,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateNote(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(note, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				NoteBodyMatching(t, recorder.Body, note)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"title":   note.Title,
				"content": note.Content.String,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateNote(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"owner":   note.Owner.String,
				"title":   note.Title.String,
				"content": note.Content.String,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateNote(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(Database.Note{}, sql.ErrNoRows)
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

			store := mockDB.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/notes"

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

func TestGetNotesApi(t *testing.T) {
	note := RandomNotes()

	testcases := []struct {
		name          string
		noteId        int32
		buildStubs    func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			noteId: note.NoteID,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetNoteById(gomock.Any(), gomock.Eq(note.NoteID)).
					Times(1).
					Return(note, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				NoteBodyMatching(t, recorder.Body, note)
			},
		},
		{
			name:   "BadRequest",
			noteId: 0,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetNoteById(gomock.Any(), gomock.Eq(note.NoteID)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			noteId: note.NoteID,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetNoteById(gomock.Any(), gomock.Eq(note.NoteID)).
					Times(1).
					Return(Database.Note{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			noteId: note.NoteID,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetNoteById(gomock.Any(), gomock.Eq(note.NoteID)).
					Times(1).
					Return(Database.Note{}, sql.ErrConnDone)
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

			store := mockDB.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/notes/%d", tc.noteId)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListNotes(t *testing.T) {
	n := 10
	notes := make([]Database.Note, n)
	for i := 0; i < n; i++ {
		notes[i] = RandomNotes()
	}

	type Query struct {
		cursor    int32
		page_size int32
	}
	testcases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockDB.MockStore, query Query)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				cursor:    1,
				page_size: int32(n),
			},
			buildStubs: func(store *mockDB.MockStore, query Query) {
				arg := Database.ListNotesParams{
					NoteID: query.cursor,
					Limit:  query.page_size,
				}
				store.EXPECT().
					ListNotes(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(notes, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				NotesBodyMatching(t, recorder.Body, notes)
			},
		},
		{
			name: "BadRequest",
			query: Query{
				cursor:    0,
				page_size: 3,
			},
			buildStubs: func(store *mockDB.MockStore, query Query) {
				arg := Database.ListNotesParams{
					NoteID: query.cursor,
					Limit:  query.page_size,
				}
				store.EXPECT().
					ListNotes(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
		{
			name: "InternalServerError",
			query: Query{
				cursor:    1,
				page_size: int32(n),
			},
			buildStubs: func(store *mockDB.MockStore, query Query) {
				arg := Database.ListNotesParams{
					NoteID: query.cursor,
					Limit:  query.page_size,
				}
				store.EXPECT().
					ListNotes(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]Database.Note{}, sql.ErrConnDone)
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

			store := mockDB.NewMockStore(ctrl)
			tc.buildStubs(store, tc.query)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/notes?cursor=%d&page_size=%d", tc.query.cursor, tc.query.page_size)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func RandomNotes() Database.Note {
	return Database.Note{
		NoteID:   int32(util.RandomInt(1, 100)),
		Owner:    sql.NullString{String: util.RandomString(5), Valid: true},
		Title:    sql.NullString{String: util.RandomString(5), Valid: true},
		Content:  sql.NullString{String: util.RandomString(8), Valid: true},
		Pinned:   sql.NullBool{Bool: false, Valid: true},
		Archived: sql.NullBool{Bool: false, Valid: true},
	}
}

func NoteBodyMatching(t *testing.T, body *bytes.Buffer, note Database.Note) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var GotNote ResponseFormat

	err = json.Unmarshal(data, &GotNote)
	require.NoError(t, err)

	expected := ResponseFormating(note)
	require.Equal(t, expected, GotNote)
}

func NotesBodyMatching(t *testing.T, body *bytes.Buffer, expected []Database.Note) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	expectedFormatted := formatManyNotes(expected)

	var gotFormatted []ResponseFormat
	err = json.Unmarshal(data, &gotFormatted)
	require.NoError(t, err)

	require.Equal(t, expectedFormatted, gotFormatted)
}

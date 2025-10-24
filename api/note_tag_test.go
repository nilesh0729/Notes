package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	Database "github.com/nilesh0729/Notes/db/Result"
	mockDB "github.com/nilesh0729/Notes/db/Mock"

	"github.com/stretchr/testify/require"
)

func TestAddTagToNote(t *testing.T) {
	notetag := RandomNoteTag()
	testcase := []struct {
		name          string
		body          gin.H
		buildstubs    func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"note_id": notetag.NoteID,
				"tag_id":  notetag.TagID,
			},
			buildstubs: func(store *mockDB.MockStore) {
				arg := Database.AddTagToNoteParams(notetag)
				store.EXPECT().
					AddTagToNote(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(notetag, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				BodyMatching(t, recorder.Body, notetag)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"tag_id": notetag.TagID,
			},
			buildstubs: func(store *mockDB.MockStore) {
				arg := Database.AddTagToNoteParams(notetag)
				store.EXPECT().
					AddTagToNote(gomock.Any(), gomock.Eq(arg)).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"note_id": notetag.NoteID,
				"tag_id":  notetag.TagID,
			},
			buildstubs: func(store *mockDB.MockStore) {
				arg := Database.AddTagToNoteParams(notetag)
				store.EXPECT().
					AddTagToNote(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(Database.NoteTag{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testcase {
		tc := testcase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildstubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/note_tags"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func RandomNoteTag() Database.NoteTag {
	note := RandomNotes()
	tag := RandomTag()

	return Database.NoteTag{
		NoteID: note.NoteID,
		TagID:  tag.TagID,
	}

}

func BodyMatching(t *testing.T, body *bytes.Buffer, notetag Database.NoteTag) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var GotNoteTag Database.NoteTag
	err = json.Unmarshal(data, &GotNoteTag)
	require.NoError(t, err)

	require.Equal(t, GotNoteTag, notetag)
}

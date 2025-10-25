package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockDB "github.com/nilesh0729/Notes/db/Mock"
	Database "github.com/nilesh0729/Notes/db/Result"
	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
)

type eqCreateParamsMatcher struct {
	arg      Database.CreateUserParams
	password string
}

func (e eqCreateParamsMatcher) Matches(x interface{}) bool {
	// In case, some value is nil
	arg, ok := x.(Database.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg Database.CreateUserParams, password string) gomock.Matcher {
	return eqCreateParamsMatcher{arg, password}
}


func TestCreateUser(t *testing.T) {
	password, user1 := RandomUser(t)
	arg := Database.CreateUserParams(user1)
	testcases := []struct {
		name          string
		body          gin.H
		buildstubbs   func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user1.Username,
				"password": password,
				"email":    user1.Email,
			},
			buildstubbs: func(store *mockDB.MockStore) {

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user1, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				UserBodyMatching(t, recorder.Body, user1)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"username": user1.Username,
				"password": "jello",
				"email":    "Hiiiii",
			},
			buildstubbs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"username": user1.Username,
				"password": password,
				"email":    user1.Email,
			},
			buildstubbs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(Database.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UniqueUserVoilation",
			body: gin.H{
				"username": user1.Username,
				"password": user1.HashedPassword,
				"email":    user1.Email,
			},
			buildstubbs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(Database.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "PasswordTooShort",
			body: gin.H{
				"username": user1.Username,
				"password": "hello",
				"email":    user1.Email,
			},
			buildstubbs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)

			tc.buildstubbs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/user"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			Body := bytes.NewBuffer(data)

			request, err := http.NewRequest(http.MethodPost, url, Body)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func RandomUser(t *testing.T) (password string, user Database.User) {
	password = util.RandomString(8)

	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	user = Database.User{
		Username: util.RandomString(6),
		HashedPassword: hashedPassword,
		Email: util.RandomEmail(),
	}
	return password, user
}

func TestGetUser(t *testing.T) {
	_, user2 := RandomUser(t)

	testcases := []struct {
		name          string
		username      string
		buildstubs    func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user2.Username,
			buildstubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user2.Username)).
					Times(1).
					Return(user2, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				UserBodyMatching(t, recorder.Body, user2)
			},
		},
		{
			name:     "BadRequest",
			username: "hello",
			buildstubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq("hello")).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "InternalServerError",
			username: user2.Username,
			buildstubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user2.Username)).
					Times(1).
					Return(Database.User{}, sql.ErrConnDone)
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
			tc.buildstubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/user/%s", tc.username)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func UserBodyMatching(t *testing.T, body *bytes.Buffer, user Database.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	// The response structure is masked
	var gotResponse UserResponseFormat
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)

	// Convert the DB user into masked form (the same way your API does)
	expectedResponse := UserResponse(user)

	// Compare the two masked responses
	require.Equal(t, expectedResponse, gotResponse)
}

package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	Database "github.com/nilesh0729/Notes/db/Result"
	"github.com/nilesh0729/Notes/tokens"
	"github.com/nilesh0729/Notes/util"
	"github.com/stretchr/testify/require"
	"fmt"
	"net/http"
)

func newTestServer(t *testing.T, store Database.Store) (*Server, tokens.Maker) {
	config := util.Config{
		Secret:              util.RandomString(32),
		AccessTokenDuration: time.Minute,
		ServerAddress:       "0.0.0.0:8080",
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server, server.tokenMaker
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker tokens.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(AuthorizationHeaderKey, authorizationHeader)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/auth"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

func newTestServer(t *testing.T, store db.Store, authClient auth.AuthImpl) *Server {
	config := util.Config{
		TokenSymmetricKey:      util.RandomString(32),
		AccessTokenDuration:    time.Minute,
		AuthenticationFilePath: "../simple-bank-firebase-authentication.json",
	}

	server, err := NewServer(config, store, authClient)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

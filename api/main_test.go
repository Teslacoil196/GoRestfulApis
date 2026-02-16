package api

import (
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		SymmetricKey:  util.RandomString(32),
		TokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)
	return server
}

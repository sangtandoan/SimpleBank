package api

import (
	"os"
	"testing"
	"time"

	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store query.Store) *Server {
	config := utils.Config{
		TokenSymKey:         utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

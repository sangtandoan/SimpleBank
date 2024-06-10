package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FrostJ143/simplebank/internal/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuth(
	t *testing.T,
	tokenMaker token.Maker,
	req *http.Request,
	authType string,
	username string,
	duration time.Duration,
) {
	accessToken, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authType, accessToken)
	req.Header.Set(autHeaderKey, authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(tokenMaker token.Maker, req *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, authTypeBearer, "user", time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthHeader",
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidFormat",
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, "", "user", time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedHeader",
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, "supported", "user", time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, authTypeBearer, "user", -time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range len(testCases) {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			server.router.GET("/auth", authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
				return
			})

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)

			tc.setupAuth(server.tokenMaker, req)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

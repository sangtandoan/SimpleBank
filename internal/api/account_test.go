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
	"time"

	mockdb "github.com/FrostJ143/simplebank/internal/mock"
	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/token"
	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	owner, _ := randomUser(t)
	account := randomAccount(owner.Username)

	// table driven tests set
	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(token.Maker, *http.Request)
		buildStubs    func(*mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "Status_OK",
			accountID: account.ID,
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, authTypeBearer, owner.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), (*account).ID).
					Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)

			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, authTypeBearer, owner.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, authTypeBearer, owner.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, authTypeBearer, owner.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
				addAuth(t, tokenMaker, req, authTypeBearer, "user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).
					Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(tokenMaker token.Maker, req *http.Request) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store)
			// save response data to recorder
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			// create http request
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(server.tokenMaker, request)
			// call http request, and response is saved in recorderd
			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})
	}

}

func randomAccount(owner string) *query.Account {
	return &query.Account{
		ID:        utils.RandomInt(0, 1000),
		Owner:     owner,
		Balance:   utils.RandomMoney(),
		Currenncy: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account *query.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	gotAccount := &query.Account{}

	err = json.Unmarshal(data, gotAccount)
	require.NoError(t, err)

	require.Equal(t, gotAccount, account)
}

package api

import (
	mock_db "TeslaCoil196/db/mock"
	db "TeslaCoil196/db/sqlc"
	"TeslaCoil196/token"
	"TeslaCoil196/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RamdonBalnce(),
		Currency: util.RandomCurrency(),
	}
}

func TestGetAccount(t *testing.T) {

	user, _ := randomUser(t)

	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(stomockStore *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, user.Username, time.Minute*2)
			},
			buildStub: func(mockStore *mock_db.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkGotAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, user.Username, time.Minute*2)
			},
			buildStub: func(mockStore *mock_db.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, user.Username, time.Minute*2)
			},
			buildStub: func(mockStore *mock_db.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidId",
			accountID: 0,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationType, user.Username, time.Minute*2)
			},
			buildStub: func(mockStore *mock_db.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			controller.Finish()

			mockStore := mock_db.NewMockStore(controller)
			tc.buildStub(mockStore)

			mockServer := newTestServer(t, mockStore)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/account/%d", tc.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setUpAuth(t, request, mockServer.tokenMaker)
			mockServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})

	}
}

func checkGotAccount(t *testing.T, body *bytes.Buffer, acccount db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, acccount, gotAccount)
}

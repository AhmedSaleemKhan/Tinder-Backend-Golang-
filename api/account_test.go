package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockauth "github.com/techschool/simplebank/auth/mock"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
)

func TestGetAccountAPI(t *testing.T) {
	var account db.Account
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwidWlkIjoiWHhpNElFQkJPWlZnN0ViYjNFWlNvRWpOOVdBbyJ9.3tUTZc4k5pNCsI2UGIXUEtEOA_E0SiDZtvDpDo-sA8I"
	firebaseTokenResponse := auth.Token{
		UID: "Xxi4IEBBOZVg7Ebb3EZSoEjN9WAo",
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					VerifyIDToken(token).
					Times(1).
					Return(&firebaseTokenResponse, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(firebaseTokenResponse.UID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NotFound",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					VerifyIDToken(token).
					Times(1).
					Return(&firebaseTokenResponse, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(firebaseTokenResponse.UID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					VerifyIDToken(token).
					Times(1).
					Return(&firebaseTokenResponse, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(firebaseTokenResponse.UID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			authClient := mockauth.NewMockAuthImpl(ctrl)
			tc.buildStubs(store, authClient)

			server := newTestServer(t, store, authClient)
			recorder := httptest.NewRecorder()

			url := "/accounts/get"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	var account db.Account
	firebaseTokenResponse := auth.Token{
		UID: "Xxi4IEBBOZVg7Ebb3EZSoEjN9WAo",
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					VerifyIDToken(gomock.Any()).
					Times(1).
					Return(&firebaseTokenResponse, nil)
				arg := db.CreateAccountParams{
					UserID:  firebaseTokenResponse.UID,
					Balance: 0,
				}

				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					VerifyIDToken(gomock.Any()).
					Times(1).
					Return(&firebaseTokenResponse, nil)
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			authClient := mockauth.NewMockAuthImpl(ctrl)
			tc.buildStubs(store, authClient)

			server := newTestServer(t, store, authClient)
			recorder := httptest.NewRecorder()

			url := "/accounts/create"
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// func TestListAccountsAPI(t *testing.T) {
// 	user := randomUser(t)

// 	n := 5
// 	accounts := make([]db.Account, n)
// 	for i := 0; i < n; i++ {
// 		accounts[i] = randomAccount(user.ID)
// 	}

// 	type Query struct {
// 		pageID   int
// 		pageSize int
// 	}

// 	testCases := []struct {
// 		name          string
// 		query         Query
// 		setupAuth     func(t *testing.T, request *http.Request, auth auth.AuthImpl)
// 		buildStubs    func(store *mockdb.MockStore, auth *mockauth.MockAuthImpl)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			buildStubs: func(store *mockdb.MockStore, auth *mockauth.MockAuthImpl) {
// 				arg := db.ListAccountsParams{
// 					UserID: user.ID,
// 					Limit:  int32(n),
// 					Offset: 0,
// 				}

// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Eq(arg)).
// 					Times(1).
// 					Return(accounts, nil)

// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchAccounts(t, recorder.Body, accounts)
// 			},
// 		},
// 		{
// 			name: "NoAuthorization",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			setupAuth: func(t *testing.T, request *http.Request, auth auth.AuthImpl) {
// 			},
// 			buildStubs: func(store *mockdb.MockStore, auth *mockauth.MockAuthImpl) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: n,
// 			},
// 			buildStubs: func(store *mockdb.MockStore, auth *mockauth.MockAuthImpl) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return([]db.Account{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPageID",
// 			query: Query{
// 				pageID:   -1,
// 				pageSize: n,
// 			},
// 			buildStubs: func(store *mockdb.MockStore, auth *mockauth.MockAuthImpl) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPageSize",
// 			query: Query{
// 				pageID:   1,
// 				pageSize: 100000,
// 			},
// 			buildStubs: func(store *mockdb.MockStore, auth *mockauth.MockAuthImpl) {
// 				store.EXPECT().
// 					ListAccounts(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			auth := mockauth.NewMockAuthImpl(ctrl)
// 			tc.buildStubs(store, auth)
// 			authClient := mockauth.NewMockAuthImpl(ctrl)

// 			server := newTestServer(t, store, authClient)
// 			recorder := httptest.NewRecorder()

// 			url := "/accounts"
// 			request, err := http.NewRequest(http.MethodGet, url, nil)
// 			require.NoError(t, err)

// 			// Add query parameters to request URL
// 			q := request.URL.Query()
// 			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
// 			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
// 			request.URL.RawQuery = q.Encode()

// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

func randomAccount(userID string) db.Account {
	return db.Account{
		UserID:  userID,
		Balance: 0,
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account.UserID, gotAccount.UserID)
	require.Equal(t, account.Balance, gotAccount.Balance)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}

package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	auth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockauth "github.com/techschool/simplebank/auth/mock"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

type eqCreateUserParamsMatcher struct {
	arg db.CreateUserParams
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	return reflect.DeepEqual(e.arg, arg)
}

func EqCreateUserParams(arg db.CreateUserParams) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg}
}

func TestCreateUserAPI(t *testing.T) {
	user := randomUser(t)

	firebaseAuthResponse := auth.UserRecord{
		UserInfo: &auth.UserInfo{
			UID: "Xxi4IEBBOZVg7Ebb3EZSoEjN9WAo",
		},
	}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
				"password":   util.RandomString(10),
				"phone":      user.Phone,
				"age":        user.Age,
				"gender":     user.Gender,
				"ethnicity":  user.Ethnicity,
				"nsfw":       user.Nsfw,
				"metadata":   string(user.Metadata),
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {

				fbAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&firebaseAuthResponse, nil)

				arg := db.CreateUserParams{
					ID:        firebaseAuthResponse.UID,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Email:     user.Email,
					Phone:     user.Phone,
					Age:       user.Age,
					Gender:    user.Gender,
					Ethnicity: user.Ethnicity,
					Nsfw:      user.Nsfw,
					Metadata:  user.Metadata,
				}

				store.EXPECT().
					CheckUserExists(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg)).
					Times(1).
					Return(user, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
				"password":   util.RandomString(10),
				"phone":      user.Phone,
				"age":        user.Age,
				"gender":     user.Gender,
				"ethnicity":  user.Ethnicity,
				"nsfw":       user.Nsfw,
				"metadata":   string(user.Metadata),
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&firebaseAuthResponse, nil)

				store.EXPECT().
					CheckUserExists(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
				"password":   util.RandomString(10),
				"phone":      user.Phone,
				"age":        user.Age,
				"gender":     user.Gender,
				"ethnicity":  user.Ethnicity,
				"nsfw":       user.Nsfw,
				"metadata":   string(user.Metadata),
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&firebaseAuthResponse, nil)

				store.EXPECT().
					CheckUserExists(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"id":         user.ID,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"password":   util.RandomString(10),
				"phone":      user.Phone,
				"age":        user.Age,
				"gender":     user.Gender,
				"ethnicity":  user.Ethnicity,
				"nsfw":       user.Nsfw,
				"metadata":   string(user.Metadata),
				"email":      "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CheckUserExists(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FireBaseError",
			body: gin.H{
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
				"password":   util.RandomString(10),
				"phone":      user.Phone,
				"age":        user.Age,
				"gender":     user.Gender,
				"ethnicity":  user.Ethnicity,
				"nsfw":       user.Nsfw,
				"metadata":   string(user.Metadata),
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, fmt.Errorf("error"))

				store.EXPECT().
					CheckUserExists(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UserAlreadyExists",
			body: gin.H{
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"email":      user.Email,
				"password":   util.RandomString(10),
				"phone":      user.Phone,
				"age":        user.Age,
				"gender":     user.Gender,
				"ethnicity":  user.Ethnicity,
				"nsfw":       user.Nsfw,
				"metadata":   string(user.Metadata),
			},
			buildStubs: func(store *mockdb.MockStore, fbAuth *mockauth.MockAuthImpl) {

				fbAuth.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					CheckUserExists(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{Email: user.Email, Phone: user.Phone}, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(user, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User) {
	user = db.User{
		ID:        util.RandomUID(28),
		FirstName: util.RandomString(6),
		LastName:  util.RandomString(6),
		Email:     util.RandomEmail(),
		Phone:     util.RandomPhoneNumber(),
		Age:       util.RandomInt(0, 150),
		Gender:    util.PickRandomGender(),
		Ethnicity: util.PickRandomEthnicity(),
		Nsfw:      util.RandomBool(),
		Metadata:  json.RawMessage("{}"),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.FirstName, gotUser.FirstName)
	require.Equal(t, user.LastName, gotUser.LastName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Age, gotUser.Age)
	require.Equal(t, user.Phone, gotUser.Phone)
}

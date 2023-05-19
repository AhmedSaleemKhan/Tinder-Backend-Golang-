package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockauth "github.com/techschool/simplebank/auth/mock"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	authorizationType string,
) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwidWlkIjoiWHhpNElFQkJPWlZnN0ViYjNFWlNvRWpOOVdBbyJ9.3tUTZc4k5pNCsI2UGIXUEtEOA_E0SiDZtvDpDo-sA8I"

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	firebaseTokenPayload := auth.Token{
		UID: "Xxi4IEBBOZVg7Ebb3EZSoEjN9WAo",
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		buildStubs    func(fbAuth *mockauth.MockAuthImpl)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, authorizationTypeBearer)
			},
			buildStubs: func(fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					VerifyIDToken(gomock.Any()).
					Times(1).
					Return(&firebaseTokenPayload, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
			buildStubs: func(fbAuth *mockauth.MockAuthImpl) {},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, "unsupported")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
			buildStubs: func(fbAuth *mockauth.MockAuthImpl) {},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, "")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
			buildStubs: func(fbAuth *mockauth.MockAuthImpl) {},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, authorizationTypeBearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
			buildStubs: func(fbAuth *mockauth.MockAuthImpl) {
				fbAuth.EXPECT().
					VerifyIDToken(gomock.Any()).
					Times(1).
					Return(nil, fmt.Errorf("expired token"))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			authClient := mockauth.NewMockAuthImpl(ctrl)
			tc.buildStubs(authClient)
			server := newTestServer(t, nil, authClient)
			authPath := "/auth"
			server.router.GET(
				authPath,
				server.authMiddleware(),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

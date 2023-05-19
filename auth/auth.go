package auth

import (
	"context"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/techschool/simplebank/util"
	"google.golang.org/api/option"
)

type AuthImpl interface {
	CreateUser(ctx *gin.Context, user *auth.UserToCreate) (*auth.UserRecord, error)
	VerifyIDToken(ctx *gin.Context, accessToken string) (*auth.Token, error)
	// DeleteUser(ctx *gin.Context, uid string) error
}

type Auth struct {
	AuthClient *auth.Client
}

func NewAuth(config util.Config) *Auth {
	authClient := SetupFirebase(config.AuthenticationFilePath)
	auth := &Auth{
		AuthClient: authClient,
	}
	return auth
}

func SetupFirebase(authenticationFilePath string) *auth.Client {
	serviceAccountKeyFilePath, err := filepath.Abs(authenticationFilePath)
	if err != nil {
		log.Fatalln("Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln("Firebase load error")
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalln("Firebase load error")
	}
	return auth
}

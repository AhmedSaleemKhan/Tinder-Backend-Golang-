package auth

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func (a *Auth) VerifyIDToken(ctx *gin.Context, accessToken string) (*auth.Token, error) {
	token, err := a.AuthClient.VerifyIDToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

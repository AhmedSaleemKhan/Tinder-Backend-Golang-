package auth

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func (a *Auth) CreateUser(ctx *gin.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
	resp, err := a.AuthClient.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// func (a *Auth) DeleteUser(ctx *gin.Context, uid string) error {
// 	err := a.AuthClient.DeleteUser(ctx, uid)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

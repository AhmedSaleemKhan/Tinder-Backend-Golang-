package api

import (
	"database/sql"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

func (server *Server) uploadPicture(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(form.File["file"]) == 0 {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	file := form.File["file"][0]
	err = fileTypeValidator(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	url, err := server.linodeSession.UploadToBucket(ctx, file, server.config)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newUploadPictureResponse(*url)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) addPictures(ctx *gin.Context) {
	var req addPictureRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	getPictures, err := server.store.GetAccountPictures(ctx, authPayload.UID)
	getPictures = append(getPictures, req.PictureURL...)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.AddAccountPicturesParams{
		ID:      authPayload.UID,
		Picture: getPictures,
	}

	updatedAccount, err := server.store.AddAccountPictures(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newModifyPictureResponse(updatedAccount, "successfully updated")
	ctx.JSON(http.StatusOK, rsp)
}

package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	checkAccountArg := db.CheckAccountExistsParams{
		Email: req.Email,
		Phone: req.Phone,
	}
	userExists, err := server.store.CheckAccountExists(ctx, checkAccountArg)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = CheckRecordExists(req.Email, req.Phone, userExists.Email, userExists.Phone)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	university := CheckNullString(req.University)
	arg := db.CreateAccountParams{
		ID:         authPayload.UID,
		FirstName:  req.FirstName,
		Email:      req.Email,
		Phone:      req.Phone,
		BirthDate:  req.BirthDate,
		Gender:     req.Gender,
		ShowMe:     req.ShowMe,
		University: university,
		Nsfw:       req.Nsfw,
		Ethnicity:  req.Ethnicity,
		Interests:  req.Interests,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newCreateAccountResponse(account)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getAccount(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	account, err := server.store.GetAccount(ctx, authPayload.UID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newGetAccountResponse(account)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) modifyAccount(ctx *gin.Context) {
	var req modifyAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	verifyYourself := CheckNullBool(req.VerifyYourself)
	aboutMe := CheckNullString(req.AboutMe)
	gender := CheckNullString(req.Gender)
	timezone := CheckNullString(req.TimeZone)
	nsfw := CheckNullBool(req.Nsfw)
	ethnicity := CheckNullString(req.Ethnicity)
	arg := db.UpdateAccountParams{
		ID:             authPayload.UID,
		VerifyYourself: verifyYourself,
		AboutMe:        aboutMe,
		Interests:      req.Interests,
		Gender:         gender,
		TimeZone:       timezone,
		Ethnicity:      ethnicity,
		Nsfw:           nsfw,
		Picture:        req.Picture,
	}

	updatedAccount, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rsp := newModifyAccountResponse(updatedAccount)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) discoverAccount(ctx *gin.Context) {
	filtratedJson := ctx.Query("json")
	filters := discoverFilters{
		PageID:   1,
		PageSize: 10,
	}
	if len(filtratedJson) > 0 {
		if err := json.Unmarshal([]byte(filtratedJson), &filters); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid filter request passed")))
			return
		}
	}
	filters.Filters.MinAge = util.ConvertAgeIntoEpochTime(filters.Filters.MinAge)
	filters.Filters.MaxAge = util.ConvertAgeIntoEpochTime(filters.Filters.MaxAge)
	arg := db.DiscoverAccountsWithFilterParams{
		Limit:       filters.PageSize,
		Offset:      (filters.PageID - 1) * filters.PageSize,
		Gender:      filters.Filters.ShowMe,
		Ethnicity:   filters.Filters.Ethnicity,
		BirthDate:   filters.Filters.MinAge,
		BirthDate_2: filters.Filters.MaxAge,
	}
	accounts, err := server.store.DiscoverAccountsWithFilter(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	metadataArg := db.DiscoverdAccountsMetadataParams{
		Gender:      filters.Filters.ShowMe,
		Ethnicity:   filters.Filters.Ethnicity,
		BirthDate:   int64(filters.Filters.MinAge),
		BirthDate_2: int64(filters.Filters.MaxAge),
	}
	totalDiscoverAccounts, err := server.store.DiscoverdAccountsMetadata(ctx, metadataArg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	metadata := computeMetadata(totalDiscoverAccounts, filters.PageID, filters.PageSize)
	rsp := newGetDiscoveredAccountsResponse(accounts, metadata)
	ctx.JSON(http.StatusOK, rsp)
}

// func (server *Server) deleteAccount(ctx *gin.Context) {
// 	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
// 	err := server.authClient.DeleteUser(ctx, authPayload.UID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	id, err := server.store.DeleteAccount(ctx, authPayload.UID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	rsp := newDeletedAccountResponse(id)
// 	ctx.JSON(http.StatusOK, rsp)
// }

// type loginUserRequest struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=6"`
// }

// type loginUserResponse struct {
// 	SessionID             uuid.UUID    `json:"session_id"`
// 	AccessToken           string       `json:"access_token"`
// 	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
// 	RefreshToken          string       `json:"refresh_token"`
// 	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
// 	User                  userResponse `json:"user"`
// }

// func (server *Server) loginUser(ctx *gin.Context) {
// 	var req loginUserRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	user, err := server.store.GetUser(ctx, req.Email)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
// 		user.Email,
// 		server.config.AccessTokenDuration,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
// 		user.Email,
// 		server.config.RefreshTokenDuration,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
// 		ID:           refreshPayload.ID,
// 		UserID:       user.ID,
// 		RefreshToken: refreshToken,
// 		UserAgent:    ctx.Request.UserAgent(),
// 		ClientIp:     ctx.ClientIP(),
// 		IsBlocked:    false,
// 		ExpiresAt:    refreshPayload.ExpiredAt,
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	rsp := loginUserResponse{
// 		SessionID:             session.ID,
// 		AccessToken:           accessToken,
// 		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
// 		RefreshToken:          refreshToken,
// 		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
// 		User:                  newUserResponse(user),
// 	}
// ctx.JSON(http.StatusOK, rsp)
// }

func CheckRecordExists(newEmail, newPhoneNumber, existedEmail, existedPhoneNumber string) error {
	if (len(existedEmail) > 0 && existedEmail == newEmail) && (len(existedPhoneNumber) > 0 && existedPhoneNumber == newPhoneNumber) {
		return fmt.Errorf("email and phone number already exist")
	} else if len(existedEmail) > 0 && existedEmail == newEmail {
		return fmt.Errorf("email already exist")
	} else if len(existedPhoneNumber) > 0 && existedPhoneNumber == newPhoneNumber {
		return fmt.Errorf("phone number already exist")
	}
	return nil
}

func CheckNullString(input string) sql.NullString {
	var value sql.NullString
	if len(input) <= 0 {
		value = sql.NullString{Valid: false}
	} else {
		value = sql.NullString{String: input, Valid: true}
	}
	return value
}

func CheckNullBool(boolean bool) sql.NullBool {
	result := sql.NullBool{Bool: boolean, Valid: true}
	return result
}

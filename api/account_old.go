package api

// import (
// 	"database/sql"
// 	"net/http"

// 	"firebase.google.com/go/v4/auth"
// 	"github.com/gin-gonic/gin"
// 	"github.com/lib/pq"
// 	db "github.com/techschool/simplebank/db/sqlc"
// )

// func (server *Server) createAccount(ctx *gin.Context) {
// 	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
// 	arg := db.CreateAccountParams{
// 		UserID:  authPayload.UID,
// 		Balance: 0,
// 	}

// 	account, err := server.store.CreateAccount(ctx, arg)
// 	if err != nil {
// 		if pqErr, ok := err.(*pq.Error); ok {
// 			switch pqErr.Code.Name() {
// 			case "foreign_key_violation", "unique_violation":
// 				ctx.JSON(http.StatusForbidden, errorResponse(err))
// 				return
// 			}
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, account)
// }

// func (server *Server) getAccount(ctx *gin.Context) {
// 	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
// 	account, err := server.store.GetAccount(ctx, authPayload.UID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, account)
// }

// type listAccountRequest struct {
// 	PageID   int32 `form:"page_id" binding:"required,min=1"`
// 	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
// }

// func (server *Server) listAccounts(ctx *gin.Context) {
// 	var req listAccountRequest
// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
// 	arg := db.ListAccountsParams{
// 		UserID: authPayload.UID,
// 		Limit:  req.PageSize,
// 		Offset: (req.PageID - 1) * req.PageSize,
// 	}

// 	accounts, err := server.store.ListAccounts(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, accounts)
// }

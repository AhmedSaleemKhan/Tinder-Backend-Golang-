package api

// import (
// 	"database/sql"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	db "github.com/techschool/simplebank/db/sqlc"
// )

// type transferRequest struct {
// 	ID            int64  `json:"id" binding:"required"`
// 	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
// 	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
// 	Amount        int64  `json:"amount" binding:"required,gt=0"`
// 	Type          string `json:"type" binding:"required"`
// }

// func (server *Server) createTransfer(ctx *gin.Context) {
// 	var req transferRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	// fromAccount, valid := server.validAccount(ctx, req.FromAccountID)
// 	// if !valid {
// 	// 	return
// 	// }
// 	// authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
// 	// if fromAccount.UserID != authPayload.UID {
// 	// 	err := errors.New("from account doesn't belong to the authenticated user")
// 	// 	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
// 	// 	return
// 	// }

// 	// _, valid = server.validAccount(ctx, req.ToAccountID)
// 	// if !valid {
// 	// 	return
// 	// }

// 	arg := db.TransferTxParams{
// 		FromAccountID: req.FromAccountID,
// 		ToAccountID:   req.ToAccountID,
// 		Amount:        req.Amount,
// 		Type:          req.Type,
// 	}

// 	result, err := server.store.TransferTx(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, result)
// }

// func (server *Server) validAccount(ctx *gin.Context, accountID string) (db.Account, bool) {
// 	account, err := server.store.GetAccount(ctx, accountID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return account, false
// 		}

// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return account, false
// 	}
// 	return account, true
// }

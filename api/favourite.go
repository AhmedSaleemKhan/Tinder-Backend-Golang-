package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

func (server *Server) addFavourite(ctx *gin.Context) {
	targetID := ctx.Param("target_id")
	favID := util.GenerateUUID()
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	checkFav := db.GetFavoriteParams{
		UserID:   authPayload.UID,
		TargetID: targetID,
	}
	isExists, err := server.store.GetFavorite(ctx, checkFav)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if isExists {
		ctx.JSON(http.StatusForbidden, errorResponse(fmt.Errorf("user already favourited")))
		return
	}
	arg := db.AddFavouriteParams{
		FavID:    favID,
		UserID:   authPayload.UID,
		TargetID: targetID,
	}
	favourite, err := server.store.AddFavourite(ctx, arg)
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
	rsp := newcreateFavouriteResponse(favourite)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getFavourites(ctx *gin.Context) {
	filtratedJson := ctx.Query("json")
	filter := getFavouritesFilter{
		PageID:   1,
		PageSize: 10,
	}
	if len(filtratedJson) > 0 {
		if err := json.Unmarshal([]byte(filtratedJson), &filter); err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid filter request passed")))
			return
		}
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	arg := db.GetAllFavouritesParams{
		ID:     authPayload.UID,
		Limit:  filter.PageSize,
		Offset: (filter.PageID - 1) * filter.PageSize,
	}
	favourites, err := server.store.GetAllFavourites(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	totalFavourites, err := server.store.FavoritesMetadata(ctx, authPayload.UID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	metadata := computeMetadata(totalFavourites, filter.PageID, filter.PageSize)
	rsp := newGetFavouritesResponse(favourites, metadata)
	ctx.JSON(http.StatusOK, rsp)
}

func computeMetadata(totalRecords int64, pageID, pageSize int32) *metadata {
	if totalRecords == 0 {
		return &metadata{}
	}
	return &metadata{
		CurrentPage:  pageID,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

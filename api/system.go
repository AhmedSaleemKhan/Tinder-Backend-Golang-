package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok\n")
}

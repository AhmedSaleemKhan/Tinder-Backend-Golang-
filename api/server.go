package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/techschool/simplebank/auth"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/linode"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config        util.Config
	store         db.Store
	tokenMaker    token.Maker
	router        *gin.Engine
	authClient    auth.AuthImpl
	linodeSession linode.LinodeImpl
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store, authClient auth.AuthImpl, r *gin.Engine) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	linodeSession := linode.NewLinodeSession(config)
	server := &Server{
		config:        config,
		store:         store,
		tokenMaker:    tokenMaker,
		authClient:    authClient,		
		linodeSession: linodeSession,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.router = server.setupRouter(r)
	return server, nil
}

func (server *Server) setupRouter(router *gin.Engine) *gin.Engine {

	router.GET("/health", server.health)
	router.Use(server.CORSMiddleware())

	accountRoutes := router.Group("/account")
	{
		accountRoutes.POST("/create", server.authMiddleware(), server.createAccount)
		accountRoutes.GET("/", server.authMiddleware(), server.getAccount)
		accountRoutes.PATCH("/modify", server.authMiddleware(), server.modifyAccount)
		accountRoutes.GET("/discover", server.authMiddleware(), server.discoverAccount)
		// accountRoutes.DELETE("/delete", server.authMiddleware(), server.deleteAccount)
		pictureRoutes := accountRoutes.Group("/picture")
		{
			pictureRoutes.POST("/upload", server.authMiddleware(), server.uploadPicture)
			pictureRoutes.PATCH("/add", server.authMiddleware(), server.addPictures)
		}

		favouriteRoutes := accountRoutes.Group("/favourite")
		{
			favouriteRoutes.POST("/add/:target_id", server.authMiddleware(), server.addFavourite)
			favouriteRoutes.GET("/", server.authMiddleware(), server.getFavourites)
		}
	}

	// authRoutes := router.Group("/accounts").Use(server.authMiddleware())

	// authRoutes.POST("/create", server.createAccount)
	// authRoutes.GET("/get", server.getAccount)
	// authRoutes.GET("/list", server.listAccounts)

	// authRoutes.POST("/transfers", server.createTransfer)

	// server.router = router
	return router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

package api

import (
	"fmt"
	"log"

	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/token"
	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      query.Store
	router     *gin.Engine
	config     utils.Config
	tokenMaker token.Maker
}

func NewServer(store query.Store, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := Server{store: store, config: config, tokenMaker: tokenMaker}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Fatal("Can not register validations")
		}
	}

	server.setupRouter()

	return &server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authGroup := router.Group("/", authMiddleware(server.tokenMaker))

	authGroup.POST("/accounts", server.createAccount)
	authGroup.GET("/accounts/:id", server.getAccount)
	authGroup.GET("/accounts", server.listAccounts)
	// authGroup.PUT("/accounts/:id", server.updateAccount)
	authGroup.DELETE("/accounts/:id", server.deleteAccount)
	authGroup.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}

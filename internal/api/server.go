package api

import (
	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *query.Store
	router *gin.Engine
}

func NewServer(store *query.Store) *Server {
	server := Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router
	return &server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"err": err.Error()}
}

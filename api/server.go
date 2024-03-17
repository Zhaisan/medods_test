package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"zhaisan-medods/utils"
)

type Server struct {
	db     *mongo.Client
	router *gin.Engine
	httpServer *http.Server
}

func NewServer(db *mongo.Client, config *utils.Config) *Server {
	server := &Server{
		db: db,
		router: gin.Default(),
	}
	server.router.POST("/token", NewAuthAPI(server.db, config).GenerateToken)
	server.router.POST("/token/refresh", NewAuthAPI(server.db, config).RefreshToken)

	server.httpServer = &http.Server{
		Addr:    config.ServerAddress,
		Handler: server.router,
	}

	return server
}

func (server *Server) Start() error {
	return server.httpServer.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.httpServer.Shutdown(ctx)
}
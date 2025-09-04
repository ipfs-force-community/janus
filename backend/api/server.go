package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server api server struct
type Server struct {
	db     *gorm.DB
	engine *gin.Engine
}

// NewServer creates a new API server instance
func NewServer(db *gorm.DB) *Server {
	s := &Server{
		db:     db,
		engine: gin.Default(),
	}
	s.registerRouter()
	return s
}

// registerRouter registers the API routes
func (s *Server) registerRouter() {
	s.engine.GET("/miners", s.GetDailyMinerStats)
}

// Run starts the server on the specified port
func (s *Server) Run(port uint16) error {
	return s.engine.Run(fmt.Sprintf(":%d", port))
}

package server

import (
	"fmt"
	"log"
	"url-collector/internal/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg config.Config
}

func New(cfg config.Config) *Server {
	return &Server{cfg}
}

func (s *Server) Start() error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	setupCollectorRoutes(router, s.cfg)

	log.Println("Starting the server with the configuration:")
	log.Println(s.cfg.String())
	
	return router.Run(fmt.Sprintf(":%d", s.cfg.Port))
}

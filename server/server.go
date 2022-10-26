package server

import (
	"github.com/gin-gonic/gin"

	"github.com/Drozd0f/csv-app/conf"
	"github.com/Drozd0f/csv-app/service"
)

type Server struct {
	*gin.Engine
	service *service.Service
	c       *conf.Config
}

func New(s *service.Service, c *conf.Config) *Server {
	if !c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	serv := &Server{
		Engine:  gin.Default(),
		service: s,
		c:       c,
	}
	serv.RegisterHandlers()

	return serv
}

func (s *Server) RegisterHandlers() {
	v1 := s.Group("/api/v1")
	{
		v1.GET("/ping", s.ping)
	}
}

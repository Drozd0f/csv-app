package server

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Drozd0f/csv-app/conf"
	_ "github.com/Drozd0f/csv-app/docs"
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
	//serv.setupSwagger()

	return serv
}

func (s *Server) RegisterHandlers() {
	v1 := s.Group("/api/v1")
	{
		v1.GET("/ping", s.ping)
	}

	s.registerCsvHandlers(v1)
	s.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

//func (s *Server) setupSwagger() {
//	ginSwagger.WrapHandler(swaggerfiles.Handler,
//		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
//		ginSwagger.DefaultModelsExpandDepth(-1))
//}

package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Drozd0f/csv-app/conf"
	_ "github.com/Drozd0f/csv-app/docs"
	"github.com/Drozd0f/csv-app/server/middleware"
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

func toGin(f http.HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		f(c.Writer, c.Request)
	}
}

func (s *Server) RegisterHandlers() {
	v1 := s.Group("/api/v1", middleware.ErrorHandler)
	{
		v1.GET("/ping", s.ping)
	}

	if s.c.Debug {
		s.GET("/debug/pprof/", toGin(pprof.Index))
		s.GET("/debug/pprof/cmdline", toGin(pprof.Cmdline))
		s.GET("/debug/pprof/profile", toGin(pprof.Profile))
		s.GET("/debug/pprof/symbol", toGin(pprof.Symbol))
		s.GET("/debug/pprof/trace", toGin(pprof.Trace))
		s.GET("/debug/pprof/allocs", gin.WrapH(pprof.Handler("allocs")))
		s.GET("/debug/pprof/block", gin.WrapH(pprof.Handler("block")))
		s.GET("/debug/pprof/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		s.GET("/debug/pprof/heap", gin.WrapH(pprof.Handler("heap")))
		s.GET("/debug/pprof/mutex", gin.WrapH(pprof.Handler("mutex")))
		s.GET("/debug/pprof/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.registerCsvHandlers(v1)
	s.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

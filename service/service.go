package service

import (
	"github.com/Drozd0f/csv-app/conf"
	"github.com/Drozd0f/csv-app/repository"
)

type Service struct {
	r *repository.Repository
	c *conf.Config
}

func New(r *repository.Repository, c *conf.Config) *Service {
	return &Service{
		r: r,
		c: c,
	}
}

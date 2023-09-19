package controller

import (
	"example.com/gitlab-ci-exporter/config"
)

type Controller struct {
	Config *config.Config
}

func New(config *config.Config) *Controller {
	return &Controller{
		Config: config,
	}
}

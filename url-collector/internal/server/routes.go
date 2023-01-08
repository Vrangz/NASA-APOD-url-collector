package server

import (
	"url-collector/internal/config"
	"url-collector/internal/server/nasa"

	"github.com/gin-gonic/gin"
)

const (
	apiV1Prefix = "/api/v1"
	nasaPrefix  = apiV1Prefix + "/nasa"
)

func setupCollectorRoutes(r *gin.Engine, cfg config.Config) {
	nasaController := nasa.New(cfg.Nasa, cfg.Timeout)
	nasaCollectorAPI := r.Group(nasaPrefix)
	{
		nasaCollectorAPI.GET("/pictures", nasaController.GetPictures)
	}
}

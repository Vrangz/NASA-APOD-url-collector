package nasa

import (
	"net/http"
	"time"
	"url-collector/internal/collector"
	"url-collector/internal/collector/nasa"
	"url-collector/internal/config"
)

type Controller struct {
	nasa collector.Collector
}

func New(cfg config.Nasa, timeout uint) *Controller {
	return &Controller{
		nasa: nasa.New(
			http.Client{Timeout: time.Duration(timeout) * time.Second},
			cfg.URL,
			cfg.ApiKey,
			cfg.ConcurrentRequests,
		),
	}
}

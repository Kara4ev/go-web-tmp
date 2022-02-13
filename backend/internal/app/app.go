package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kara4ev/go-web-tmp/internal/config"
	"github.com/Kara4ev/go-web-tmp/pkg/httpserver"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {

	logger.Debug("start app %v version %v", cfg.AppName, cfg.AppVersion)

	gin.SetMode(gin.ReleaseMode)
	// if cfg.AppDebug == 0 {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	ds, err := initDS(cfg)

	if err != nil {
		logger.Fatal("unable initialize data sources : %v", err)
	}

	router, err := inject(ds, *cfg)

	if err != nil {
		logger.Fatal("failure to inject data sources: %v\n", err)
	}

	httpServer := httpserver.New(httpserver.SConfig{
		Hendler: router,
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTPHost, cfg.HTTPPort),
	})

	logger.Info(fmt.Sprintf("Listening port %s:%s", cfg.HTTPHost, cfg.HTTPPort))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("in os signal: %s", s.String())
	case err := <-httpServer.Notify():
		logger.Fatal("in http server notify: %w", err)
	}

	// Shutdown
	if err = httpServer.Shutdown(); err != nil {
		logger.Error("error http server shutdown: %w", err)
	}

	if err := ds.close(); err != nil {
		logger.Error("error data sourse close: %w", err)
	}

}

package main

import (
	"fmt"

	"github.com/Kara4ev/go-web-tmp/internal/app"
	"github.com/Kara4ev/go-web-tmp/internal/config"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		panic(fmt.Sprintf("init config error: %s", err))
	}

	logger.InitLogger(cfg.LoggerLevel, cfg.AppLogFile)
	app.Run(cfg)
}

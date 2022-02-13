package app

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Kara4ev/go-web-tmp/internal/config"
	"github.com/Kara4ev/go-web-tmp/internal/handler"
	"github.com/Kara4ev/go-web-tmp/internal/repository"
	"github.com/Kara4ev/go-web-tmp/internal/service"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func inject(d *dataSource, cfg config.Config) (*gin.Engine, error) {
	logger.Debug("injecting data source")

	/*
	*	init
	 */

	// load rsa key
	logger.Debug("read private key")
	priv, err := ioutil.ReadFile(cfg.AppPrivateKeyFile)
	if err != nil {
		logger.Debug("could not read private key pem file: %w", err)
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	logger.Debug("parse private key")
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		logger.Debug("could not parse private key: %w", err)
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	logger.Debug("read public key")
	pub, err := ioutil.ReadFile(cfg.AppPublicKeyFile)

	if err != nil {
		logger.Debug("could not read public key pem file: %w", err)
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	logger.Debug("parse public key")
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		logger.Debug("could not parse public key: %w", err)
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// gin init
	logger.Debug("create router")
	router := gin.Default()

	/*
	* reposytory layer
	 */
	logger.Debug("create user repository")
	userReposytory := repository.NewUserReposytory(d.DB)
	toketRepository := repository.NewTokenRepository(d.Radis)

	/*
	* service layer
	 */
	logger.Debug("create user services")
	userService := service.NewUserServices(&service.USConfig{
		UserRepository: userReposytory,
	})

	logger.Debug("create token services")
	tokenService := service.NewTokenService(&service.TSConfig{
		TokenRepository:       toketRepository,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         cfg.AppSecret,
		RefrashExpirationSecs: cfg.AppRefreshTokenExpiration,
		IDExpirationSecs:      cfg.AppIDTokenExpiration,
	})

	/*
	* hendler layer
	 */

	logger.Debug("create handler")
	handler.NewHandler(&handler.Config{
		Router:          router,
		UserService:     userService,
		TokenService:    tokenService,
		BaseUrl:         cfg.HTTPBaseURL,
		TimeoutDuration: time.Duration(time.Duration(cfg.HTTPHendlerTimeOut) * time.Second),
	})

	logger.Debug("data source injecting")
	return router, nil

}

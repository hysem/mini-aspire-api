package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hysem/mini-aspire-api/app/config"
	"github.com/hysem/mini-aspire-api/app/core/bcrypt"
	"github.com/hysem/mini-aspire-api/app/core/db"
	"github.com/hysem/mini-aspire-api/app/core/jwt"
	"github.com/hysem/mini-aspire-api/app/handler"
	"github.com/hysem/mini-aspire-api/app/repository"
	"github.com/hysem/mini-aspire-api/app/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// App holds application details
type App struct {
	e *echo.Echo

	provider struct {
		masterDB *sqlx.DB
		bcrypt   bcrypt.Bcrypt
		jwt      jwt.JWT
	}

	repository struct {
		user repository.User
	}

	usecase struct {
		user usecase.User
	}

	handler struct {
		misc *handler.Misc
		user *handler.User
	}
}

// New instatiates the app
func New() *App {
	return &App{
		e: echo.New(),
	}
}

// Run starts the application and wait for requests until stopped
func (a *App) Run() {
	for _, initFn := range []func() error{
		a.initProviders,
		a.initRepositories,
		a.initUsecases,
		a.initHandlers,
		a.initMiddlewares,
		a.initRoutes,
	} {
		if err := initFn(); err != nil {
			zap.L().Fatal("failed to intialize", zap.Error(err))
		}
	}

	// Start server
	go func() {
		address := fmt.Sprintf(":%d", config.Current().App.Port)
		if err := a.e.Start(address); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("shutting down the server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		zap.L().Fatal("failed shutdown", zap.Error(err))
	}
}

// initRoutes intialize the routes
func (a *App) initRoutes() error {
	publicV1Group := a.e.Group("/v1")
	publicV1UserGroup := publicV1Group.Group("/user")

	// GET /v1/ping
	publicV1Group.GET("/ping", a.handler.misc.Ping)

	// POST /v1/user/token
	publicV1UserGroup.POST("/token", a.handler.user.GenerateToken)

	// GET /docs
	a.e.Group("/docs", middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  config.Current().App.Docs,
		Index: "index.html",
		HTML5: true,
	}))
	return nil
}

// initProviders initializes the middlewares
func (a *App) initProviders() error {
	var err error

	a.provider.masterDB, err = db.Connect(config.Current().Database.Master)
	if err != nil {
		return errors.Wrap(err, "failed to connect to master db")
	}

	a.provider.bcrypt, err = bcrypt.New(config.Current().Bcrypt)
	if err != nil {
		return errors.Wrap(err, "failed to create bcrypt provider")
	}

	a.provider.jwt = jwt.New(config.Current().JWT)

	return nil
}

// initRepositories initializes the repositories
func (a *App) initRepositories() error {
	a.repository.user = repository.NewUser(a.provider.masterDB)
	return nil
}

// initUsecases initializes the usecases
func (a *App) initUsecases() error {
	a.usecase.user = usecase.NewUser(a.repository.user, a.provider.bcrypt, a.provider.jwt)
	return nil
}

// initHandlers initializes the handlers
func (a *App) initHandlers() error {
	a.handler.misc = handler.NewMisc()
	a.handler.user = handler.NewUser(a.usecase.user)
	return nil
}

// initMiddlewares initializes the middlewares
func (a *App) initMiddlewares() error {
	return nil
}

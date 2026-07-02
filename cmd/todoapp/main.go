package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/mlkad/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/mlkad/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/mlkad/golang-todoapp/internal/core/transport/http/server"
	users_postgres_repository "github.com/mlkad/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/mlkad/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/mlkad/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	//Если придет SIGINT или SIGTERM, этот ctx станет отмененным
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
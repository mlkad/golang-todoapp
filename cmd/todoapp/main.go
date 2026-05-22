package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_http_server "github.com/mlkad/golang-todoapp/internal/core/transport/http/server"
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

	logger.Debug("starting ToDo application!")

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
	)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
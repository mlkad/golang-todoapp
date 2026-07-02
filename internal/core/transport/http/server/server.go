package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/mlkad/golang-todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux *http.ServeMux
	config Config
	log *core_logger.Logger
	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
	) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		config: config,
		log: log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

// метод, который запускает HTTP-сервер и умеет его корректно останавливать (это называется graceful shutdown — "мягкое завершение").
func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware... )

	server := &http.Server{
		Addr: s.config.Address,
		Handler: mux,
	}

	ch := make(chan error, 1) // канал ошибок
	go func() {
		defer close(ch) //когда горутина закончится, закроем канал

		s.log.Warn("Start HTTP server", zap.String("addr", s.config.Address))
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select { // ждём, что произойдёт первым
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout,
		)
		defer cancel()

		//просим завершиться аккуратно и не принимать новые http запросы
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close() // возвр ошибка

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		s.log.Warn("HTTP server stopped")
	}
	return nil
}
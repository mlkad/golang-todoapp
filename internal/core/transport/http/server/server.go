package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux *http.ServeMux
	config Config
	log *core_logger.Logger
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		config: config,
		log: log,
	}
}

func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

// метод, который запускает HTTP-сервер и умеет его корректно останавливать (это называется graceful shutdown — "мягкое завершение").
func (h *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr: h.config.Address,
		Handler: h.mux,
	}

	ch := make(chan error, 1) // канал ошибок
	go func() {
		defer close(ch) //когда горутина закончится, закроем канал

		h.log.Warn("Start HTTP server", zap.String("addr", h.config.Address))
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
		h.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout,
		)
		defer cancel()

		//просим завершиться аккуратно и не принимать новые http запросы
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close() // возвр ошибка

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		h.log.Warn("HTTP server stopped")
	}
	return nil
}
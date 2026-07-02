package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_http_response "github.com/mlkad/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestID := r.Header.Get(requestIDHeader)
				// клиент не передал айдишник запроса
				if requestID == "" {
					requestID = uuid.NewString()
				}
				r.Header.Set(requestIDHeader, requestID)
				w.Header().Set(requestIDHeader, requestID)

				next.ServeHTTP(w, r)
		})
	}
}

// зашивает логгер в каждый хэндлер
func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()

			log.Debug(
				">> incoming HTTP request",
				zap.String("http_method", r.Method),
			 	zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r) //след это уже хэндлер

			log.Debug("<<< done http request", zap.Int("status_code", rw.GetStatusCode()), zap.Duration("latency", time.Now().Sub(before)))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handler HTTP request got unexpected panic",)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}


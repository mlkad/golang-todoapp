package users_transport_http

import (
	"net/http"

	"github.com/mlkad/golang-todoapp/internal/core/domain"
	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_http_request "github.com/mlkad/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/mlkad/golang-todoapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

// http ответ это ответ на user dto
type CreateUserResponse UserDTOResponse

// CreateUser обрабатывает HTTP POST запрос на создание нового пользователя.
// Процесс:
// 1. Извлекает логгер и HTTP ResponseWriter из контекста
// 2. Создает обработчик ошибок для красивого ответа клиенту
// 3. Декодирует JSON из тела запроса в структуру CreateUserRequest
// 4. Валидирует данные по правилам (required, min, max и т.д.)
// 5. Если ошибка → отправляет HTTP ответ (400/500) и останавливает функцию
// 6. Если успех → продолжает обработку (создание в БД, отправка ответа)

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke CreateUser handle")

	var req CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	userDomain := domainFromDTO(req)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}

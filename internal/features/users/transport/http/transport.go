package users_transport_http

import (
	"net/http"

	core_http_server "github.com/mlkad/golang-todoapp/internal/core/transport/http/server"
)

//Сюда будут приходить запросы типа POST /users
type UsersHTTPHandler struct {
	usersService UsersService
}

//Интерфейс описывает — какие методы у сервиса должны быть
type UsersService interface {
}

func NewUsersHTTPHandler(usersService UsersService,) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method: http.MethodPost,
			Path: "/users",
			Handler: h.CreateUser,
		},
	}
}
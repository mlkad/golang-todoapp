package users_transport_http

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


package users_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserRequest struct {
	FullName 		string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

type CreateUserResponse struct {
	ID 					int			`json:"id"`
	Version 		int			`json:"version"`
	FullName 		string	`json:"full_name"`
	PhoneNumber *string	`json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		 fmt.Println("ошибка")
	}
}


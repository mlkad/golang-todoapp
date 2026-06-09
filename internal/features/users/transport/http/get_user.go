package users_transport_http

import (
	"net/http"

	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_http_response "github.com/mlkad/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/mlkad/golang-todoapp/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	// GET /users/{id}
	userID, err := core_http_utils.GetIntPassValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}
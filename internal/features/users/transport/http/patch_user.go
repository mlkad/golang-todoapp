package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_http_request "github.com/mlkad/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/mlkad/golang-todoapp/internal/core/transport/http/response"
)
/*
нужно отличать сценарий, когда поле просто не было передано, от сценария, когда поле было явно передано со значением null

если *string то
            ┌─→ "88776"  → позиция "значение"  ✅ (видно отдельно)
клиент ─────┤
            ├─→ null      ─┐
            └─→ молчание  ─┴→ позиция "nil"  ❌ (слиплись, не разобрать кто есть кто)
*/


type PatchUserRequest struct {
	FullName string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var req PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)

		return
	}

	log.Debug(
		fmt.Sprintf(
			"PatchUserRequest fields:\nFullName: '%s'\nPhoneNumber: '%s'",
			req.FullName,
			req.PhoneNumber,
		),
	)
	rw.WriteHeader(http.StatusOK)
}
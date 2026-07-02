package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mlkad/golang-todoapp/internal/core/domain"
	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_http_request "github.com/mlkad/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/mlkad/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/mlkad/golang-todoapp/internal/core/transport/http/types"
)

/*
нужно отличать сценарий, когда поле просто не было передано, от сценария, когда поле было явно передано со значением null

если *string то
            ┌─→ "88776"  → позиция "значение"  ✅ (видно отдельно) но не отличает явный null от отсутствия поля
клиент ─────┤
            ├─→ null      ─┐
            └─→ молчание  ─┴→ позиция "nil"  ❌ (слиплись, не разобрать кто есть кто)
*/


type PatchUserRequest struct {
	FullName core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName can't be NULL`")
		}
		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName must be between 3 and 100 symbols`")
		}
	}
	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must be startswith '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPassValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}

	var req PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)

		return
	}

	userPatch := userPatchFromRequest(req)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	) 
}
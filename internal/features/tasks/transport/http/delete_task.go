package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/mlkad/golang-todoapp/internal/core/logger"
	core_http_request "github.com/mlkad/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/mlkad/golang-todoapp/internal/core/transport/http/response"
)

func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.NoContentResponse() //204
}

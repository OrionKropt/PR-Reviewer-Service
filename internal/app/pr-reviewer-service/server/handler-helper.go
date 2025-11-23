package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/OrionKropt/PRReviewerService/api/types"
)

func writeError(w http.ResponseWriter, log *slog.Logger,
	status int, code types.ErrorCode, msg, handlerName string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := types.ErrorResponse{
		Error: types.ErrorDetail{
			Code:    code,
			Message: msg,
		},
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to write error response", "handler", handlerName, "error", err)
	}
}

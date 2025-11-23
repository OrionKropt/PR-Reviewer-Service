package server

import (
	"encoding/json"
	"net/http"

	"github.com/OrionKropt/PRReviewerService/api/types"
)

func (s *Server) handleUsersSetIsActive() http.HandlerFunc {
	request := struct {
		IsActive bool   `json:"is_active"`
		UserId   string `json:"user_id"`
	}{}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			s.log.Error("failed to decode JSON request", "handler", "handleUsersSetIsActive", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		user, err := s.prRevService.SetIsActive(request.UserId, request.IsActive)
		if err != nil {
			s.log.Error("failed to set isActive", "handler", "handleUsersSetIsActive", "error", err.Error())
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode(types.ErrorResponse{Error: types.ErrorDetail{Code: types.NotFound}})
			if err != nil {
				s.log.Error("failed to encode JSON response", "handler", "handleUsersSetIsActive", "error", err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(user); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handleUsersSetIsActive", "error", err.Error())
		}
	}
}

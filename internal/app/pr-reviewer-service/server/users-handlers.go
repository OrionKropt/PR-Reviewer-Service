package server

import (
	"encoding/json"
	"net/http"

	"github.com/OrionKropt/PRReviewerService/api/types"
)

func (s *Server) handleUsersSetIsActive() http.HandlerFunc {
	type Request struct {
		IsActive bool   `json:"is_active"`
		UserId   string `json:"user_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		request := Request{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			s.log.Error("failed to decode JSON request", "handler", "handleUsersSetIsActive", "error", err.Error())
			writeError(w, s.log, http.StatusBadRequest, types.BadRequest, "", "handleUsersSetIsActive")
			return
		}

		user, err := s.prRevService.SetIsActive(request.UserId, request.IsActive)
		if err != nil {
			s.log.Error("failed to set isActive", "handler", "handleUsersSetIsActive", "error", err.Error())
			writeError(w, s.log, http.StatusNotFound, types.NotFound, "", "handleUsersSetIsActive")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(user); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handleUsersSetIsActive", "error", err.Error())
		}
	}
}

func (s *Server) handleUsersGetReview() http.HandlerFunc {
	type Response struct {
		UserID string                   `json:"user_id"`
		PRs    []types.PullRequestShort `json:"pull_requests"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		prs, err := s.prRevService.GetPullRequestsAsReviewer(userID)
		response := Response{UserID: userID}
		if err != nil {
			s.log.Error("failed to get pull requests as reviewer", "handler", "handleUsersGetReview")
			response.PRs = make([]types.PullRequestShort, 0)
		} else {
			response.PRs = prs
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handleUsersGetReview", "error", err.Error())
		}
	}
}

package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OrionKropt/PRReviewerService/api/types"
	"github.com/OrionKropt/PRReviewerService/internal/app/pr-reviewer-service/domain"
)

func (s *Server) handlePRCreate() http.HandlerFunc {
	request := struct {
		PRID     string `json:"pull_request_id"`
		PRName   string `json:"pull_request_name"`
		AuthorId string `json:"author_id"`
	}{}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			s.log.Error("failed to decode JSON request", "handler", "handleUsersSetIsActive", "error", err.Error())
			writeError(w, s.log, http.StatusBadRequest, types.BadRequest, "", "handleUsersSetIsActive")
		}

		pr, err := s.prRevService.CreatePullRequest(request.PRID, request.PRName, request.AuthorId)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrPRExists):
				s.log.Error("failed to create the pull request", "handler", "handleUsersSetIsActive", "error", err.Error())
				writeError(w, s.log, http.StatusConflict, types.PRExists, "PR id already exists", "handleUsersSetIsActive")

			case errors.Is(err, domain.ErrAuthorNotFound):
				s.log.Error("failed to create the pull request", "handler", "handleUsersSetIsActive", "error", err.Error())
				writeError(w, s.log, http.StatusNotFound, types.NotFound, "author not found", "handleUsersSetIsActive")

			case errors.Is(err, domain.ErrTeamNotFound):
				s.log.Error("failed to create the pull request", "handler", "handleUsersSetIsActive", "error", err.Error())
				writeError(w, s.log, http.StatusNotFound, types.NotFound, "team not found", "handleUsersSetIsActive")

			default:
				s.log.Error("failed to create the pull request", "handler", "handleUsersSetIsActive", "error", err.Error())
				writeError(w, s.log, http.StatusInternalServerError, types.InternalError, "", "handleUsersSetIsActive")
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(pr); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handleUsersSetIsActive", "error", err.Error())
		}
	}
}

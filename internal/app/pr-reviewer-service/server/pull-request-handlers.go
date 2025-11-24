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
			s.log.Error("failed to decode JSON request", "handler", "handlePRCreate", "error", err.Error())
			writeError(w, s.log, http.StatusBadRequest, types.BadRequest, "", "handlePRCreate")
			return
		}

		pr, err := s.prRevService.CreatePullRequest(request.PRID, request.PRName, request.AuthorId)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrPRExists):
				s.log.Error("failed to create the pull request", "handler", "handlePRCreate", "error", err.Error())
				writeError(w, s.log, http.StatusConflict, types.PRExists, "PR id already exists", "handlePRCreate")

			case errors.Is(err, domain.ErrAuthorNotFound):
				s.log.Error("failed to create the pull request", "handler", "handlePRCreate", "error", err.Error())
				writeError(w, s.log, http.StatusNotFound, types.NotFound, "author not found", "handlePRCreate")

			case errors.Is(err, domain.ErrTeamNotFound):
				s.log.Error("failed to create the pull request", "handler", "handlePRCreate", "error", err.Error())
				writeError(w, s.log, http.StatusNotFound, types.NotFound, "team not found", "handlePRCreate")

			default:
				s.log.Error("failed to create the pull request", "handler", "handlePRCreate", "error", err.Error())
				writeError(w, s.log, http.StatusInternalServerError, types.InternalError, "", "handlePRCreate")
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(pr); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handlePRCreate", "error", err.Error())
		}
	}
}

func (s *Server) handlePRMerge() http.HandlerFunc {
	request := struct {
		PRID string `json:"pull_request_id"`
	}{}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			s.log.Error("failed to decode JSON request", "handler", "handlePRMerge", "error", err.Error())
			writeError(w, s.log, http.StatusBadRequest, types.BadRequest, "", "handlePRMerge")
			return
		}

		pr, err := s.prRevService.MergePullRequest(request.PRID)
		if err != nil {
			s.log.Error("failed to merge pull request", "handler", "handlePRMerge", "error", err.Error())
			writeError(w, s.log, http.StatusNotFound, types.NotFound, "", "handlePRMerge")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(pr); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handlePRMerge", "error", err.Error())
		}
	}
}

func (s *Server) handlePRReassign() http.HandlerFunc {
	request := struct {
		PRID      string `json:"pull_request_id"`
		OldUserID string `json:"old_user_id"`
	}{}
	response := struct {
		PR         types.PullRequest `json:"pr"`
		ReplacedBy string            `json:"replaced_by"`
	}{}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			s.log.Error("failed to decode JSON request", "handler", "handlePRReassign", "error", err.Error())
			writeError(w, s.log, http.StatusBadRequest, types.BadRequest, "", "handlePRReassign")
			return
		}
		pr, replacedByID, err := s.prRevService.ReassignReviewerPullRequest(request.PRID, request.OldUserID)
		if err != nil {
			s.log.Error("failed to reassign pull requests reviewer", "handler", "handlePRReassign", "error", err.Error())
			switch {
			case errors.Is(err, domain.ErrTeamNotFound):
				fallthrough
			case errors.Is(err, domain.ErrPRNotFound):
				fallthrough
			case errors.Is(err, domain.ErrUserNotFound):
				writeError(w, s.log, http.StatusNotFound, types.NotFound, "", "handlePRReassign")
				return
			case errors.Is(err, domain.ErrUserNotActive):
				writeError(w, s.log, http.StatusConflict, types.NoCandidate, "no active replacement candidate in team", "handlePRReassign")
				return
			case errors.Is(err, domain.ErrPRMerged):
				writeError(w, s.log, http.StatusConflict, types.PRMerged, "cannot reassign on merged PR", "handlePRReassign")
				return
			case errors.Is(err, domain.ErrUserNotAssigned):
				writeError(w, s.log, http.StatusNotFound, types.NotAssigned, "reviewer is not assigned to this PR", "handlePRReassign")
				return
			default:
				writeError(w, s.log, http.StatusBadRequest, types.BadRequest, "", "handlePRReassign")
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response.PR = *pr
		response.ReplacedBy = replacedByID
		if err := json.NewEncoder(w).Encode(response); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handlePRReassign", "error", err.Error())
		}
	}
}

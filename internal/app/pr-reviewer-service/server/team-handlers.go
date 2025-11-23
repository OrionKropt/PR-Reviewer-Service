package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OrionKropt/PRReviewerService/api/types"
)

func (s *Server) handleTeamAdd() http.HandlerFunc {
	team := types.Team{}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := s.prRevService.AddTeam(team.TeamName, team.Members); err != nil {
			s.log.Error("failed to add team", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(types.ErrorResponse{Error: types.ErrorDetail{
				Message: fmt.Sprintf("%s already exists", team.TeamName),
				Code:    types.TeamExists,
			}})
			if err != nil {
				s.log.Error("failed to encode JSON response", "handler", "handleTeamAdd", "error", err)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(team); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handleTeamAdd", "error", err.Error())
		}
	}
}

func (s *Server) handleTeamGet() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		teamName := r.URL.Query().Get("team_name")

		team, err := s.prRevService.GetTeam(teamName)
		if err != nil {
			s.log.Error("failed to get team", "error", err.Error())
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(types.ErrorResponse{Error: types.ErrorDetail{
				Code: types.NotFound,
			}})
			if err != nil {
				s.log.Error("failed to encode JSON response", "handler", "handleTeamGet", "error", err.Error())
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(team); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handleTeamGet", "error", err.Error())
		}
	}
}

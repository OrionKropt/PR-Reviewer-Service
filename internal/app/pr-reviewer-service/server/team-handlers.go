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
			s.log.Error("failed to decode JSON request", "handler", "handleTeamAdd", "error", err.Error())
			writeError(w, s.log, http.StatusBadRequest, types.BadRequest, "", "handleTeamAdd")
			return
		}

		if err := s.prRevService.CreateTeam(team.TeamName, team.Members); err != nil {
			s.log.Error("failed to add team", "error", err.Error())
			writeError(w, s.log, http.StatusBadRequest, types.TeamExists, fmt.Sprintf("%s already exists", team.TeamName), "handleTeamAdd")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(team); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handleTeamAdd", "error", err.Error())
		}
	}
}

func (s *Server) handleTeamGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")

		team, err := s.prRevService.GetTeam(teamName)
		if err != nil {
			s.log.Error("failed to get team", "error", err.Error())
			writeError(w, s.log, http.StatusNotFound, types.NotFound, "", "handleTeamGet")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(team); err != nil {
			s.log.Error("failed to encode JSON response", "handler", "handleTeamGet", "error", err.Error())
		}
	}
}

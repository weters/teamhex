/*
Copyright 2020 Tom Peters

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package controller Team Hex API
//
// The Team Hex API provides information on professional and collegiate sports teams' colors. This API powers the Team Hex website: https://teamhex.dev/.
//
// Terms Of Service: Use at your own risk.
// Host: api.teamhex.dev
// Version: 1.0
// License: Apache License, Version 2.0
// Schemes: https
// Produces:
// - application/json
//
// swagger:meta
package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/weters/teamhex/internal/model"
)

//Controller provides capabilities for handling HTTP requests
type Controller struct {
	*mux.Router
	model   *model.Model
	version string
}

//New returns a new instance of the controller
//This instance implements the methods required of an HTTP handler
func New(m *model.Model, version string) *Controller {
	c := Controller{
		model:   m,
		version: version,
	}

	router := mux.NewRouter()
	c.Router = router

	router.Methods(http.MethodGet).Path("/").Handler(c.getRoot())
	router.Methods(http.MethodGet).Path("/swagger.json").Handler(c.getSwaggerJSON())
	router.Methods(http.MethodGet).Path("/teams").Handler(c.getTeams())
	router.Methods(http.MethodGet).Path("/leagues").Handler(c.getLeagues())
	router.Methods(http.MethodGet).Path("/leagues/{league:[^/]+}").Handler(c.getLeaguesLeague())
	router.Methods(http.MethodGet).Path("/leagues/{league:[^/]+}/{team:[^/]+}").Handler(c.getLeaguesLeagueTeam())

	return &c
}

// Successful response
// swagger:response rootResponse
type rootResponse struct {
	Version        string    `json:"version"`
	GenerationDate time.Time `json:"generationDate"`
	Links          []string  `json:"_links"`
}

// swagger:route GET / version
//
// Get health and version
//
// Provides version information and links to other resources
//
// Produces:
// - application/json
//
// Responses:
//   200: rootResponse
func (c *Controller) getRoot() http.HandlerFunc {
	resp := rootResponse{
		Version:        c.version,
		GenerationDate: c.model.GenerationDate(),
		Links: []string{
			"/teams{?search}",
			"/leagues",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		serveJSON(w, http.StatusOK, resp)
	}
}

// Successful response
// swagger:response leaguesResponse
type leaguesResponse []*model.LeagueRecord

// swagger:route GET /leagues leagues getLeagues
//
// Returns a list of supported leagues
//
// This endpoint will return a list of all leagues in the system.
//
// Produces:
// - application/json
//
// Responses:
//   200: leaguesResponse
func (c *Controller) getLeagues() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveJSON(w, http.StatusOK, c.model.Leagues())
	}
}

// Successful response
// swagger:response teamsResponse
type teamsResponse model.Teams

// swagger:operation GET /teams getTeams
//
// Returns a list of teams
//
// By default, this endpoint will return all teams. You can search using the search query parameter.
//
// ---
// produces:
// - application/json
// parameters:
// - name: search
//   in: query
//   description: Search for the specified team
//   required: false
//   type: string
// responses:
//   '200':
//     '$ref': '#/responses/teamsResponse'
func (c *Controller) getTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s := r.FormValue("search"); len(s) > 0 {
			serveJSON(w, http.StatusOK, c.model.Search(s))
			return
		}

		serveJSON(w, http.StatusOK, c.model.AllTeams())
	}
}

// swagger:operation GET /leagues/{league} getTeamsByLeague
//
// Get teams in a league
//
// This endpoint returns a list of teams found in a provided league.
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: league
//   required: true
//   type: string
// responses:
//   '200':
//     '$ref': '#/responses/teamsResponse'
//   '404':
//     '$ref': '#/responses/errorResponse'
//   '500':
//     '$ref': '#/responses/errorResponse'
func (c *Controller) getLeaguesLeague() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		league := mux.Vars(r)["league"]
		teams, err := c.model.TeamsByLeague(league)
		if err != nil {
			if err == model.ErrLeagueNotFound {
				serveJSONError(w, http.StatusNotFound, errors.New("league not found"))
				return
			}

			serveJSONError(w, http.StatusInternalServerError, err)
			return
		}
		serveJSON(w, http.StatusOK, teams)
	}
}

// Successful response
// swagger:response teamResponse
type teamResponse *model.Team

// swagger:operation GET /leagues/{league}/{team} getTeam
//
// Get a single team in a provided league
//
// This endpoint returns a list of teams found in a provided league.
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: league
//   required: true
//   type: string
// - in: path
//   name: team
//   required: true
//   type: string
// responses:
//   '200':
//     '$ref': '#/responses/teamResponse'
//   '404':
//     '$ref': '#/responses/errorResponse'
//   '500':
//     '$ref': '#/responses/errorResponse'
func (c *Controller) getLeaguesLeagueTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		leagueName := mux.Vars(r)["league"]
		teamName := mux.Vars(r)["team"]
		team, err := c.model.TeamByLeagueAndName(leagueName, teamName)
		if err != nil {
			if err == model.ErrLeagueNotFound {
				serveJSONError(w, http.StatusNotFound, errors.New("league not found"))
				return
			} else if err == model.ErrTeamNotFound {
				serveJSONError(w, http.StatusNotFound, errors.New("team not found"))
				return
			}

			serveJSONError(w, http.StatusInternalServerError, err)
			return
		}
		serveJSON(w, http.StatusOK, team)
	}
}

func (c *Controller) getSwaggerJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.json")
	}
}

func serveJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logrus.WithError(err).Error("could not encode JSON")
	}
}

// An error response
// swagger:response errorResponse
type errorResponse struct {
	Message string `json:"message"`
}

func serveJSONError(w http.ResponseWriter, statusCode int, err error) {
	var msg string
	if err != nil {
		if statusCode/100 == 5 {
			logrus.Error(err)
		}

		msg = err.Error()
	} else {
		msg = http.StatusText(statusCode)
	}

	serveJSON(w, statusCode, errorResponse{Message: msg})
}

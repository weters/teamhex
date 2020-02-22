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

package controller

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/weters/teamhex/internal/model"
	"net/http"
	"time"
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
	router.Methods(http.MethodGet).Path("/teams").Handler(c.getTeams())
	router.Methods(http.MethodGet).Path("/leagues").Handler(c.getLeagues())
	router.Methods(http.MethodGet).Path("/leagues/{league:[^/]+}").Handler(c.getLeaguesLeague())
	router.Methods(http.MethodGet).Path("/leagues/{league:[^/]+}/{team:[^/]+}").Handler(c.getLeaguesLeagueTeam())

	return &c
}

func (c *Controller) getRoot() http.HandlerFunc {
	resp := struct {
		Version        string    `json:"version"`
		GenerationDate time.Time `json:"generationDate"`
		Links          []string  `json:"_links"`
	}{
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

func (c *Controller) getLeagues() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveJSON(w, http.StatusOK, c.model.Leagues())
	}
}

func (c *Controller) getTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s := r.FormValue("search"); len(s) > 0 {
			serveJSON(w, http.StatusOK, c.model.Search(s))
			return
		}

		serveJSON(w, http.StatusOK, c.model.AllTeams())
	}
}

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

func serveJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logrus.WithError(err).Error("could not encode JSON")
	}
}

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

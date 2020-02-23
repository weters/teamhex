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

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

//ErrLeagueNotFound represents an error when the league is not found
var ErrLeagueNotFound = errors.New("model: league not found")

//ErrTeamNotFound represents an error when the team is not found
var ErrTeamNotFound = errors.New("model: team not found")

//Model provides capabilities for finding team colors
type Model struct {
	raw           *DataFile
	leagues       []*LeagueRecord
	teamsByLeague map[string]*leagueData
}

//DataFile represents how the file is stored on disk
type DataFile struct {
	Generated time.Time `json:"generated"`
	Teams     Teams     `json:"teams"`
}

//New returns a new model instance
//An error is returned if the file cannot be found or parsed.
func New(dataFilename string) (*Model, error) {
	file, err := os.Open(dataFilename)
	if err != nil {
		return nil, err
	}

	var data DataFile
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	sort.Sort(data.Teams)

	teamsByLeague := make(map[string]*leagueData)
	uniqLeagues := make(map[string]bool)

	for _, team := range data.Teams {
		uniqLeagues[team.League] = true

		league := strings.ToLower(team.League)
		_, ok := teamsByLeague[league]
		if !ok {
			teamsByLeague[league] = &leagueData{
				sortedTeams: make(Teams, 0, 1),
				teamByName:  make(map[string]*Team),
			}
		}

		team.Link = fmt.Sprintf("/leagues/%s/%s", url.PathEscape(league), url.PathEscape(strings.ToLower(team.Name)))

		teamsByLeague[league].sortedTeams = append(teamsByLeague[league].sortedTeams, team)
		teamsByLeague[league].teamByName[strings.ToLower(team.Name)] = team
	}

	leagues := make([]*LeagueRecord, 0, len(uniqLeagues))
	for league := range uniqLeagues {
		leagues = append(leagues, &LeagueRecord{
			League: league,
			Link:   fmt.Sprintf("/leagues/%s", url.PathEscape(strings.ToLower(league))),
		})
	}

	sort.Sort(sortByLeagueRecord(leagues))

	return &Model{
		raw:           &data,
		leagues:       leagues,
		teamsByLeague: teamsByLeague,
	}, nil
}

//AllTeams returns all teams
func (m *Model) AllTeams() Teams {
	return m.raw.Teams
}

//Leagues returns a list of all the leagues
func (m *Model) Leagues() []*LeagueRecord {
	return m.leagues
}

//TeamByLeagueAndName returns a team by the league and team name
func (m *Model) TeamByLeagueAndName(leagueName, name string) (*Team, error) {
	lowerLeagueName := strings.ToLower(leagueName)
	lowerName := strings.ToLower(name)

	league, ok := m.teamsByLeague[lowerLeagueName]
	if !ok {
		return nil, ErrLeagueNotFound
	}

	team, ok := league.teamByName[lowerName]
	if !ok {
		return nil, ErrTeamNotFound
	}

	return team, nil
}

//TeamsByLeague returns a list of all teams in a given league
func (m *Model) TeamsByLeague(league string) (Teams, error) {
	teams, ok := m.teamsByLeague[strings.ToLower(league)]
	if !ok {
		return nil, ErrLeagueNotFound
	}

	return teams.sortedTeams, nil
}

//Search will search a team in by its name
func (m *Model) Search(match string) Teams {
	teams := make(Teams, 0)
	for _, team := range m.raw.Teams {
		if strings.Contains(strings.ToLower(team.Name), strings.ToLower(match)) {
			teams = append(teams, team)
		}
	}

	return teams
}

//GenerationDate returns the date the color data was generated
func (m *Model) GenerationDate() time.Time {
	return m.raw.Generated
}

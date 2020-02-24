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

import "strings"

//LeagueRecord represents an individual league
type LeagueRecord struct {
	// League is the name of the league
	League string `json:"league"`
	// Link is a link to retrieve teams for that league
	Link string `json:"_link"`
}

type sortByLeagueRecord []*LeagueRecord

func (s sortByLeagueRecord) Len() int {
	return len(s)
}

func (s sortByLeagueRecord) Less(i, j int) bool {
	return strings.Compare(s[i].League, s[j].League) < 0
}

func (s sortByLeagueRecord) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type leagueData struct {
	sortedTeams Teams
	teamByName  map[string]*Team
}

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
	"github.com/onsi/gomega"
	"testing"
	"time"
)

const testFile = "testdata/teamhex.json"

func TestLoadModel(t *testing.T) {
	g := gomega.NewWithT(t)
	m, err := New("badfile.json")
	g.Expect(m).Should(gomega.BeNil())
	g.Expect(err).ShouldNot(gomega.BeNil())

	m, err = New("testdata/invalid.json")
	g.Expect(m).Should(gomega.BeNil())
	g.Expect(err).ShouldNot(gomega.BeNil())

	m, err = New(testFile)
	g.Expect(m).ShouldNot(gomega.BeNil())
	g.Expect(err).Should(gomega.BeNil())
}

func TestGenerationDate(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)
	g.Expect(m.GenerationDate()).Should(gomega.Equal(time.Date(2020, 2, 22, 12, 0, 0, 0, time.UTC)))
}

func TestAllTeams(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)

	teams := m.AllTeams()
	g.Expect(len(teams)).Should(gomega.Equal(4))
	g.Expect(teams[0].Name).Should(gomega.Equal("Buffalo Bills"))
	g.Expect(teams[1].Name).Should(gomega.Equal("Buffalo Sabres"))
	g.Expect(teams[2].Name).Should(gomega.Equal("The Ohio State University"))
	g.Expect(teams[3].Name).Should(gomega.Equal("University At Buffalo, The State University Of New York"))
	g.Expect(teams[3].Link).Should(gomega.Equal("/leagues/ncaa/university%20at%20buffalo%2C%20the%20state%20university%20of%20new%20york"))
}

func TestLeagues(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)
	leagues := m.Leagues()
	g.Expect(len(leagues)).Should(gomega.Equal(3))
	g.Expect(leagues[0].League).Should(gomega.Equal("NCAA"))
	g.Expect(leagues[1].League).Should(gomega.Equal("NFL"))
	g.Expect(leagues[2].League).Should(gomega.Equal("NHL"))
}

func TestTeamByLeagueAndName(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)

	team, err := m.TeamByLeagueAndName("bad", "bdar")
	g.Expect(team).Should(gomega.BeNil())
	g.Expect(err).Should(gomega.MatchError(ErrLeagueNotFound))

	team, err = m.TeamByLeagueAndName("nfl", "bdar")
	g.Expect(team).Should(gomega.BeNil())
	g.Expect(err).Should(gomega.MatchError(ErrTeamNotFound))

	team, err = m.TeamByLeagueAndName("NfL", "bUfFaLo BiLlS")
	g.Expect(err).Should(gomega.BeNil())
	g.Expect(team).ShouldNot(gomega.BeNil())
	g.Expect(team.Name).Should(gomega.Equal("Buffalo Bills"))
	g.Expect(team.League).Should(gomega.Equal("NFL"))
	g.Expect(team.Link).Should(gomega.Equal("/leagues/nfl/buffalo%20bills"))
	g.Expect(team.Eras).Should(gomega.Equal([]*Era{
		{
			Year: 2011, Colors: []*Color{
				{Name: "Royal Blue", Hex: "#003087"},
				{Name: "Scarlet Red", Hex: "#C8102E"},
			},
		},
		{
			Year: 2002, Colors: []*Color{
				{Name: "Midnight Navy", Hex: "#091F2C"},
			},
		},
	}))
	g.Expect(team.Division).Should(gomega.Equal("AFC"))
}

func TestTeamsByLeague(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)

	teams, err := m.TeamsByLeague("bad")
	g.Expect(teams).Should(gomega.BeNil())
	g.Expect(err).Should(gomega.MatchError(ErrLeagueNotFound))

	teams, err = m.TeamsByLeague("nCaA")
	g.Expect(err).Should(gomega.BeNil())
	g.Expect(len(teams)).Should(gomega.Equal(2))

	teams, err = m.TeamsByLeague("nFl")
	g.Expect(err).Should(gomega.BeNil())
	g.Expect(len(teams)).Should(gomega.Equal(1))
}

func TestSearch(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)

	teams := m.Search("bad")
	g.Expect(len(teams)).Should(gomega.Equal(0))

	teams = m.Search("bill")
	g.Expect(len(teams)).Should(gomega.Equal(1))
	g.Expect(teams[0].Name).Should(gomega.Equal("Buffalo Bills"))

	teams = m.Search("buffalo")
	g.Expect(len(teams)).Should(gomega.Equal(3))
	g.Expect(teams[0].Name).Should(gomega.Equal("Buffalo Bills"))
	g.Expect(teams[1].Name).Should(gomega.Equal("Buffalo Sabres"))
	g.Expect(teams[2].Name).Should(gomega.Equal("University At Buffalo, The State University Of New York"))
}

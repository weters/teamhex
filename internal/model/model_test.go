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
	g.Expect(teams[0].Name).Should(gomega.Equal("Apples"))
	g.Expect(teams[1].Name).Should(gomega.Equal("Bananas"))
	g.Expect(teams[2].Name).Should(gomega.Equal("Pizzas"))
	g.Expect(teams[3].Name).Should(gomega.Equal("The Pears"))
	g.Expect(teams[3].Link).Should(gomega.Equal("/leagues/fruit/the%20pears"))
}

func TestLeagues(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New("testdata/teamhex-moreleagues.json")
	leagues := m.Leagues()
	g.Expect(len(leagues)).Should(gomega.Equal(3))
	g.Expect(leagues[0].League).Should(gomega.Equal("Berry"))
	g.Expect(leagues[1].League).Should(gomega.Equal("Hesperidium"))
	g.Expect(leagues[2].League).Should(gomega.Equal("Pome"))
}

func TestTeamByLeagueAndName(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)

	team, err := m.TeamByLeagueAndName("foo", "bar")
	g.Expect(team).Should(gomega.BeNil())
	g.Expect(err).Should(gomega.MatchError(ErrLeagueNotFound))

	team, err = m.TeamByLeagueAndName("fruit", "bar")
	g.Expect(team).Should(gomega.BeNil())
	g.Expect(err).Should(gomega.MatchError(ErrTeamNotFound))

	team, err = m.TeamByLeagueAndName("fRuIt", "aPpLeS")
	g.Expect(err).Should(gomega.BeNil())
	g.Expect(team).ShouldNot(gomega.BeNil())
	g.Expect(team.Name).Should(gomega.Equal("Apples"))
	g.Expect(team.League).Should(gomega.Equal("Fruit"))
	g.Expect(team.Link).Should(gomega.Equal("/leagues/fruit/apples"))
	g.Expect(team.Colors).Should(gomega.Equal([]string{"#f00", "#0f0"}))
	g.Expect(team.Conference).Should(gomega.Equal("Sweet"))
}

func TestTeamsByLeague(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)

	teams, err := m.TeamsByLeague("bad")
	g.Expect(teams).Should(gomega.BeNil())
	g.Expect(err).Should(gomega.MatchError(ErrLeagueNotFound))

	teams, err = m.TeamsByLeague("frUit")
	g.Expect(err).Should(gomega.BeNil())
	g.Expect(len(teams)).Should(gomega.Equal(3))

	teams, err = m.TeamsByLeague("italian FOOD")
	g.Expect(err).Should(gomega.BeNil())
	g.Expect(len(teams)).Should(gomega.Equal(1))
}

func TestSearch(t *testing.T) {
	g := gomega.NewWithT(t)
	m, _ := New(testFile)

	teams := m.Search("bad")
	g.Expect(len(teams)).Should(gomega.Equal(0))

	teams = m.Search("izz")
	g.Expect(len(teams)).Should(gomega.Equal(1))
	g.Expect(teams[0].Name).Should(gomega.Equal("Pizzas"))

	teams = m.Search("as")
	g.Expect(len(teams)).Should(gomega.Equal(2))
	g.Expect(teams[0].Name).Should(gomega.Equal("Bananas"))
	g.Expect(teams[1].Name).Should(gomega.Equal("Pizzas"))
}

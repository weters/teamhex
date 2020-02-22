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
	"github.com/onsi/gomega"
	"github.com/weters/teamhex/internal/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testFile = "./testdata/teamhex.json"

var g *gomega.GomegaWithT
var m *model.Model
var ts *httptest.Server

func runWithSetupAndTeardown(t *testing.T, tests func()) {
	g = gomega.NewWithT(t)
	m, _ = model.New(testFile)
	c := New(m, "v1.0.0")
	ts = httptest.NewServer(c)
	defer ts.Close()

	tests()
}

func TestGetRoot(t *testing.T) {
	type response struct {
		Version        string    `json:"version"`
		GenerationDate time.Time `json:"generationDate"`
		Links          []string  `json:"_links"`
	}
	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL)
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)
		g.Expect(string(body)).Should(gomega.Equal(toJSON(response{
			Version:        "v1.0.0",
			GenerationDate: time.Date(2020, 2, 22, 12, 0, 0, 0, time.UTC),
			Links: []string{
				"/teams{?search}",
				"/leagues",
			},
		})))
	})
}

func TestGetLeagues(t *testing.T) {
	expected := `[
	{  "league": "Fruit", "_link": "/leagues/fruit" },
	{  "league": "Italian Food", "_link": "/leagues/italian%20food" }
]`
	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/leagues")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var leagues []*model.LeagueRecord
		must(json.Unmarshal([]byte(expected), &leagues))
		g.Expect(string(body)).Should(gomega.Equal(toJSON(leagues)))
	})
}

func TestGetTeamsByLeague(t *testing.T) {
	expected := `[
  { "name": "Apples", "colors": [ "#f00", "#0f0" ], "league": "Fruit", "conference": "Sweet", "link": "/leagues/fruit/apples" },
  { "name": "Bananas", "colors": [ "#ff0", "#000" ], "league": "Fruit", "link": "/leagues/fruit/bananas"},
  { "name": "The Pears", "colors": [ "#af0", "#ff0" ], "league": "Fruit", "link": "/leagues/fruit/the%20pears"}
]`

	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/leagues/fruit")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var teams []*model.TeamRecord
		must(json.Unmarshal([]byte(expected), &teams))
		g.Expect(string(body)).Should(gomega.Equal(toJSON(teams)))
	})
}

func TestGetTeamsByLeagueWithLeagueNotFound(t *testing.T) {
	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/leagues/bad")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusNotFound))
		body, _ := ioutil.ReadAll(res.Body)
		g.Expect(string(body)).Should(gomega.Equal(`{"message":"league not found"}` + "\n"))
	})
}

func TestGetTeamByLeagueAndName(t *testing.T) {
	expected := `
{ "name": "The Pears", "colors": [ "#af0", "#ff0" ], "league": "Fruit", "link": "/leagues/fruit/the%20pears"}
`

	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/leagues/fruit/the%20pears")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var teams *model.TeamRecord
		must(json.Unmarshal([]byte(expected), &teams))
		g.Expect(string(body)).Should(gomega.Equal(toJSON(teams)))
	})
}

func TestGetTeamByLeagueAndNameWithLeagueNotFound(t *testing.T) {
	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/leagues/badleague/the%20pears")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusNotFound))
		body, _ := ioutil.ReadAll(res.Body)
		g.Expect(string(body)).Should(gomega.Equal(`{"message":"league not found"}` + "\n"))
	})
}

func TestGetTeamByLeagueAndNameWithTeamNotFound(t *testing.T) {
	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/leagues/italian%20food/the%20pears")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusNotFound))
		body, _ := ioutil.ReadAll(res.Body)
		g.Expect(string(body)).Should(gomega.Equal(`{"message":"team not found"}` + "\n"))
	})
}

func TestGetTeamsByAll(t *testing.T) {
	expected := `[
  { "name": "Apples", "colors": [ "#f00", "#0f0" ], "league": "Fruit", "conference": "Sweet", "link": "/leagues/fruit/apples" },
  { "name": "Bananas", "colors": [ "#ff0", "#000" ], "league": "Fruit", "link": "/leagues/fruit/bananas"},
  { "name": "Pizzas", "colors": [ "#b83", "#f00", "#000" ], "league": "Italian Food", "link": "/leagues/italian%20food/pizzas"},
  { "name": "The Pears", "colors": [ "#af0", "#ff0" ], "league": "Fruit", "link": "/leagues/fruit/the%20pears"}
]`

	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/teams")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var teams []*model.TeamRecord
		must(json.Unmarshal([]byte(expected), &teams))
		g.Expect(string(body)).Should(gomega.Equal(toJSON(teams)))
	})
}

func TestGetTeamsBySearch(t *testing.T) {
	expected := `[
  { "name": "Bananas", "colors": [ "#ff0", "#000" ], "league": "Fruit", "link": "/leagues/fruit/bananas"},
  { "name": "Pizzas", "colors": [ "#b83", "#f00", "#000" ], "league": "Italian Food", "link": "/leagues/italian%20food/pizzas"}
]`

	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/teams?search=as")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var teams []*model.TeamRecord
		must(json.Unmarshal([]byte(expected), &teams))
		g.Expect(string(body)).Should(gomega.Equal(toJSON(teams)))
	})
}

func TestGetTeamsBySearchNoResults(t *testing.T) {
	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/teams?search=bad+search")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		g.Expect(string(body)).Should(gomega.Equal("[]\n"))
	})
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func toJSON(o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return string(b) + "\n"
}

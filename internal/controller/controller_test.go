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
	{  "league": "NCAA", "_link": "/leagues/ncaa" },
	{  "league": "NFL", "_link": "/leagues/nfl" },
	{  "league": "NHL", "_link": "/leagues/nhl" }
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
    {
      "name": "The Ohio State University",
      "eras": [
        {
          "year": 2004,
          "colors": [
            { "name": "Scarlet", "hex": "#BA0C2F" }
          ]
        }
      ],
      "league": "NCAA",
      "division": "Big Ten Conference",
      "_link": "/leagues/ncaa/the%20ohio%20state%20university"
    },
    {
      "name": "University At Buffalo, The State University Of New York",
      "eras": [
        {
          "year": 2016,
          "colors": [
            { "name": "Royal Blue", "hex": "#0057B7" },
            { "name": "White", "hex": "#FFFFFF" }
          ]
        }
      ],
      "league": "NCAA",
      "division": "Mid-American Conference",
      "_link": "/leagues/ncaa/university%20at%20buffalo%2C%20the%20state%20university%20of%20new%20york"
    }
]`

	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/leagues/ncaa")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var teams []*model.Team
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
    {
      "name": "Buffalo Sabres",
      "eras": [
        {
          "year": 2010,
          "colors": [ { "name": "Navy", "hex": "#041E42" } ]
        }
      ],
      "league": "NHL",
	  "_link": "/leagues/nhl/buffalo%20sabres"
    }
`

	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/leagues/nhl/buffalo%20sabres")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var teams *model.Team
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
		res, err := http.Get(ts.URL + "/leagues/nfl/oakland%20raiders")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusNotFound))
		body, _ := ioutil.ReadAll(res.Body)
		g.Expect(string(body)).Should(gomega.Equal(`{"message":"team not found"}` + "\n"))
	})
}

func TestGetTeamsByAll(t *testing.T) {
	expected := `[
    {
      "name": "Buffalo Bills",
      "eras": [
        {
          "year": 2011,
          "colors": [
            { "name": "Royal Blue", "hex": "#003087" },
            { "name": "Scarlet Red", "hex": "#C8102E" }
          ]
        },
        {
          "year": 2002,
          "colors": [
            { "name": "Midnight Navy", "hex": "#091F2C" }
          ]
        }
      ],
      "league": "NFL",
      "division": "AFC",
	  "_link": "/leagues/nfl/buffalo%20bills"
    },
    {
      "name": "Buffalo Sabres",
      "eras": [
        {
          "year": 2010,
          "colors": [ { "name": "Navy", "hex": "#041E42" } ]
        }
      ],
      "league": "NHL",
	  "_link": "/leagues/nhl/buffalo%20sabres"
    },
    {
      "name": "The Ohio State University",
      "eras": [
        {
          "year": 2004,
          "colors": [
            { "name": "Scarlet", "hex": "#BA0C2F" }
          ]
        }
      ],
      "league": "NCAA",
      "division": "Big Ten Conference",
      "_link": "/leagues/ncaa/the%20ohio%20state%20university"
    },
    {
      "name": "University At Buffalo, The State University Of New York",
      "eras": [
        {
          "year": 2016,
          "colors": [
            { "name": "Royal Blue", "hex": "#0057B7" },
            { "name": "White", "hex": "#FFFFFF" }
          ]
        }
      ],
      "league": "NCAA",
      "division": "Mid-American Conference",
      "_link": "/leagues/ncaa/university%20at%20buffalo%2C%20the%20state%20university%20of%20new%20york"
    }
  ]`

	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/teams")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var teams []*model.Team
		must(json.Unmarshal([]byte(expected), &teams))
		g.Expect(string(body)).Should(gomega.Equal(toJSON(teams)))
	})
}

func TestGetTeamsBySearch(t *testing.T) {
	expected := `[
    {
      "name": "The Ohio State University",
      "eras": [
        {
          "year": 2004,
          "colors": [
            { "name": "Scarlet", "hex": "#BA0C2F" }
          ]
        }
      ],
      "league": "NCAA",
      "division": "Big Ten Conference",
      "_link": "/leagues/ncaa/the%20ohio%20state%20university"
    },
    {
      "name": "University At Buffalo, The State University Of New York",
      "eras": [
        {
          "year": 2016,
          "colors": [
            { "name": "Royal Blue", "hex": "#0057B7" },
            { "name": "White", "hex": "#FFFFFF" }
          ]
        }
      ],
      "league": "NCAA",
      "division": "Mid-American Conference",
      "_link": "/leagues/ncaa/university%20at%20buffalo%2C%20the%20state%20university%20of%20new%20york"
    }
]`

	runWithSetupAndTeardown(t, func() {
		res, err := http.Get(ts.URL + "/teams?search=univ")
		g.Expect(err).Should(gomega.BeNil())
		defer res.Body.Close()
		g.Expect(res.StatusCode).Should(gomega.Equal(http.StatusOK))
		body, _ := ioutil.ReadAll(res.Body)

		var teams []*model.Team
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

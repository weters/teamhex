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

//TeamRecord represents an individual team
type TeamRecord struct {
	Name       string   `json:"name"`
	Colors     []string `json:"colors"`
	League     string   `json:"league"`
	Conference string   `json:"conference,omitempty"`
	Link       string   `json:"link"`
}

//Teams is a colletion of teams
type Teams []*TeamRecord

func (t Teams) Len() int {
	return len(t)
}

func (t Teams) Less(i, j int) bool {
	return strings.Compare(t[i].Name, t[j].Name) < 0
}

func (t Teams) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

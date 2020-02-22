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

package main

import (
	"flag"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/weters/teamhex/internal/controller"
	"github.com/weters/teamhex/internal/model"
	"net/http"
	"os"
	"time"
)

const readTimeout = time.Second * 10
const writeTimeout = time.Second * 5

var Version = "v0.0.0"
var addr = flag.String("addr", ":5000", "address to listen on")
var dataFilename = flag.String("file", "teamhex.json", "path to JSON colors file")

func main() {
	flag.Parse()

	m, err := model.New(*dataFilename)
	if err != nil {
		logrus.WithError(err).Fatal("could not load model")
	}
	c := controller.New(m, Version)

	corsHandler := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet},
	})

	server := &http.Server{
		Addr:         *addr,
		Handler:      handlers.CombinedLoggingHandler(os.Stdout, corsHandler.Handler(c)),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	logrus.WithField("addr", server.Addr).Info("Server started")
	logrus.Fatal(server.ListenAndServe())
}

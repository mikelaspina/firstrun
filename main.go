/*
Copyright 2014 Michael S. LaSpina

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
	"log"
	"net/http"

	"github.com/mikelaspina/firstrun/pkg/server"
)

type Config struct {
	BindAddr string
	DocRoot  string
}

var defaultConfig = Config{
	BindAddr: ":8080",
	DocRoot:  "./public",
}

func main() {
	config := defaultConfig

	staticHandler := http.FileServer(http.Dir(config.DocRoot))
	http.Handle("/css/", staticHandler)
	http.Handle("/fonts/", staticHandler)
	http.Handle("/images/", staticHandler)
	http.Handle("/js/", staticHandler)

	scheduleHandler := &server.ScheduleHandler{}
	if err := scheduleHandler.Init(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", scheduleHandler)
	http.Handle("/schedule/", scheduleHandler)
	log.Fatal(http.ListenAndServe(config.BindAddr, nil))
}

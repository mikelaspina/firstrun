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

package tv

import (
	"time"
)

type Episode struct {
	Series  string    `json:"series"`
	Season  int       `json:"season"`
	Number  int       `json:"number"`
	Title   string    `json:"title"`
	AirDate time.Time `json:"airdate"`
}

func (self *Episode) Aired(year int, month time.Month, day int) {
	self.AirDate = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

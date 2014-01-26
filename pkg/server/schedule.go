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

package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/mikelaspina/firstrun/pkg/tv"
)

type ScheduleHandler struct {
	t     *template.Template
	sched tv.Schedule
}

func (self *ScheduleHandler) load(r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&self.sched.Episodes)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (self *ScheduleHandler) Init() error {
	var err error

	self.t, err = template.ParseFiles(
		"templates/upcoming.html",
		"templates/schedule.html",
		"templates/schedule-show-group.html")
	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Open("./data.json")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := self.load(f); err != nil {
		return err
	}

	log.Printf("Loaded %d episodes\n", len(self.sched.Episodes))
	return nil
}

func (h *ScheduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	switch r.URL.Path {
	case "/upcoming":
		h.upcoming(w, r)
	default:
		h.index(w, r)
	}
}

type indexPage struct {
	Title  string
	Series []indexGroup
}

type indexGroup struct {
	Title      string
	ShowBadge  bool
	BadgeCount int
	Episodes   []indexGroupItem
}

type indexGroupItem struct {
	Type    string
	Title   string
	AirDate string
	Link    string
}

type airDateSorter struct {
	Episodes []*tv.Episode
}

func (self airDateSorter) Len() int {
	return len(self.Episodes)
}

func (self airDateSorter) Less(i, j int) bool {
	return self.Episodes[i].AirDate.Before(self.Episodes[j].AirDate)
}

func (self airDateSorter) Swap(i, j int) {
	self.Episodes[i], self.Episodes[j] = self.Episodes[j], self.Episodes[i]
}

func byDate(eps []*tv.Episode) []*tv.Episode {
	s := airDateSorter{eps}
	sort.Sort(s)
	return s.Episodes
}

func (self *ScheduleHandler) index(w http.ResponseWriter, r *http.Request) {
	page := indexPage{Title: "TV Schedule"}
	for series, eps := range groupBySeries(self.sched.Episodes) {
		group := indexGroup{
			Title:      series,
			Episodes:   unwatched(eps, 0),
			BadgeCount: badges(eps),
		}

		if group.BadgeCount > 0 {
			group.ShowBadge = true
		}

		page.Series = append(page.Series, group)
	}

	if err := self.t.ExecuteTemplate(w, "schedule", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (self *ScheduleHandler) upcoming(w http.ResponseWriter, r *http.Request) {
	page := indexPage{Title: "Upcoming Episodes"}
	for series, eps := range groupBySeries(self.sched.Episodes) {
		group := indexGroup{
			Title:    series,
			Episodes: filterUpcoming(eps),
		}

		if group.BadgeCount > 0 {
			group.ShowBadge = true
		}

		page.Series = append(page.Series, group)
	}

	if err := self.t.ExecuteTemplate(w, "upcoming", page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func groupBySeries(eps []*tv.Episode) map[string][]*tv.Episode {
	groups := make(map[string][]*tv.Episode)
	for _, ep := range eps {
		groups[ep.Series] = append(groups[ep.Series], ep)
	}
	return groups
}

func filterUpcoming(eps []*tv.Episode) []indexGroupItem {
	items := make([]indexGroupItem, 0)
	cutoff := today()

	for _, ep := range byDate(eps) {
		if ep.AirDate.After(cutoff) {
			items = append(items, indexGroupItem{
				Type:    fmt.Sprintf("S%d : Ep. %d", ep.Season, ep.Number),
				Title:   ep.Title,
				AirDate: ep.AirDate.Format("01/02/2006"),
			})
		}
	}

	return items
}

var watchLinks = map[string]string{
	"Castle":                "http://www.hulu.com/profile/queue",
	"How I Met Your Mother": "http://www.cbs.com/shows/how_i_met_your_mother/",
	"New Girl":              "http://www.hulu.com/profile/queue",
	"The Blacklist":         "http://www.hulu.com/profile/queue",
	"The Mentalist":         "http://www.cbs.com/shows/the_mentalist/video/",
}

func link(series string) string {
	if link, ok := watchLinks[series]; ok {
		return link
	}
	return "#"
}

func unwatched(eps []*tv.Episode, maxUpcoming int) []indexGroupItem {
	items := make([]indexGroupItem, 0)
	cutoff := today()
	upcoming := 0
	for _, ep := range byDate(eps) {
		if ep.AirDate.After(cutoff) {
			if upcoming >= maxUpcoming {
				break
			}
			upcoming += 1
		}

		items = append(items, indexGroupItem{
			Type:    fmt.Sprintf("S%d : Ep. %d", ep.Season, ep.Number),
			Title:   ep.Title,
			AirDate: ep.AirDate.Format("01/02/2006"),
			Link:    link(ep.Series),
		})
	}

	return items
}

func badges(eps []*tv.Episode) int {
	count := 0
	cutoff := today()
	for _, ep := range eps {
		if ep.AirDate.Before(cutoff) {
			count += 1
		}
	}
	return count
}

func today() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

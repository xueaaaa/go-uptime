package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	sitemodel "github.com/xueaaaa/go-uptime/internal/site/model"
	"github.com/xueaaaa/go-uptime/internal/util"
)

type siteRow struct {
	URL       string `json:"url"`
	Class     string `json:"class"`
	Interval  string `json:"interval"`
	LastCheck string `json:"lastCheck"`
}

type indexData struct {
	Total       int       `json:"total"`
	Available   int       `json:"available"`
	Unavailable int       `json:"unavailable"`
	Unknown     int       `json:"unknown"`
	Sites       []siteRow `json:"sites"`
}

func (s *Server) getSites(w http.ResponseWriter, r *http.Request) (indexData, error) {
	sites, err := s.siteSvc.GetAll(r.Context())
	if err != nil {
		return indexData{}, err
	}

	rows := make([]siteRow, 0)
	available, unavailable, unknown := 0, 0, 0
	for _, site := range sites {
		row := siteRow{
			URL:       site.URL,
			Interval:  util.FormatDuration(site.Interval),
			LastCheck: "—",
		}
		switch site.Status {
		case sitemodel.Available:
			row.Class = "available"
			available++
		case sitemodel.Unavailable:
			row.Class = "unavailable"
			unavailable++
		default:
			row.Class = "unknown"
			unknown++
		}
		if site.LastCheckAt != nil {
			elapsed := time.Since(*site.LastCheckAt)

			switch {
			case elapsed < time.Minute:
				row.LastCheck = "just now"
			case elapsed < time.Hour:
				row.LastCheck = util.Pluralize(int(elapsed.Minutes()), "minute ago", "minutes ago")
			case elapsed < 24*time.Hour:
				row.LastCheck = util.Pluralize(int(elapsed.Hours()), "hour ago", "hours ago")
			default:
				row.LastCheck = site.LastCheckAt.Format("2006 Jan 02")
			}
		}
		rows = append(rows, row)
	}

	data := indexData{Total: len(sites), Available: available, Unavailable: unavailable, Unknown: unknown, Sites: rows}
	return data, nil
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	data, err := s.getSites(w, r)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&buf, "index.html", data); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

func (s *Server) sitesEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		data, err := s.getSites(w, r)
		if err != nil {
			return
		}
		jsn, err := json.Marshal(data)
		if err != nil {
			return
		}

		fmt.Fprintf(w, "data: %s\n\n", jsn)
		flusher.Flush()

		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
		}
	}
}

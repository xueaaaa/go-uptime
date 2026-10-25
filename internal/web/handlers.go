package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/google/uuid"
	errors2 "github.com/xueaaaa/go-uptime/internal/errors"
	sitemodel "github.com/xueaaaa/go-uptime/internal/site/model"
	"github.com/xueaaaa/go-uptime/internal/util"
)

type siteRow struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	Class     string    `json:"class"`
	Interval  string    `json:"interval"`
	LastCheck string    `json:"lastCheck"`
}

type checkRow struct {
	Class        string
	StatusCode   int
	Latency      int
	LatencyClass string
	Error        string
	CheckedAt    string
}

type indexData struct {
	Total       int       `json:"total"`
	Available   int       `json:"available"`
	Unavailable int       `json:"unavailable"`
	Unknown     int       `json:"unknown"`
	Sites       []siteRow `json:"sites"`
}

type siteData struct {
	ID            uuid.UUID
	URL           string
	FormattedURL  string
	Interval      string
	IntervalError bool
	Total         int
	Success       int
	Failure       int
	Unknown       int
	Checks        []checkRow
}

func formatCheckTime(checkedAt time.Time) string {
	elapsed := time.Since(checkedAt)
	switch {
	case elapsed < time.Minute:
		return "just now"
	case elapsed < time.Hour:
		return util.Pluralize(int(elapsed.Minutes()), "minute ago", "minutes ago")
	case elapsed < 24*time.Hour:
		return util.Pluralize(int(elapsed.Hours()), "hour ago", "hours ago")
	default:
		return checkedAt.Format("2006 Jan 02")
	}
}

func (s *Server) getSites(w http.ResponseWriter, r *http.Request) (indexData, error) {
	sites, err := s.siteSvc.GetAll(r.Context())
	if err != nil && !errors.Is(err, errors2.NoSites) {
		return indexData{}, err
	}

	rows := make([]siteRow, 0)
	available, unavailable, unknown := 0, 0, 0
	for _, site := range sites {
		row := siteRow{
			ID:        site.ID,
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
			row.LastCheck = formatCheckTime(*site.LastCheckAt)
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

func (s *Server) site(w http.ResponseWriter, r *http.Request) {
	siteIdVal := r.PathValue("site-id")
	if siteIdVal == "" {
		http.Error(w, "empty site id", http.StatusBadRequest)
		return
	}

	siteId, err := uuid.Parse(siteIdVal)
	if err != nil {
		http.Error(w, "invalid site id", http.StatusBadRequest)
		return
	}

	site, err := s.siteSvc.Get(r.Context(), siteId)
	if err != nil {
		if errors.Is(err, errors2.NotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	u, err := url.Parse(site.URL)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	checks, err := s.checkSvc.GetBySiteID(r.Context(), siteId)
	if err != nil && !errors.Is(err, errors2.NotFound) {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	sort.Slice(checks, func(i, j int) bool { return checks[i].CheckedAt.After(checks[j].CheckedAt) })

	rows := make([]checkRow, 0)
	success, failure, unknown := 0, 0, 0
	for _, check := range checks {
		row := checkRow{
			StatusCode: check.StatusCode,
			Latency:    check.Latency,
			Error:      check.Error,
			CheckedAt:  formatCheckTime(check.CheckedAt),
		}
		switch {
		case check.Latency < 500:
			row.LatencyClass = "great"
		case check.Latency < 2000:
			row.LatencyClass = "medium"
		default:
			row.LatencyClass = "bad"
		}
		switch check.Status {
		case sitemodel.Available:
			row.Class = "success"
			success++
		case sitemodel.Unavailable:
			row.Class = "failure"
			failure++
		default:
			row.Class = "unknown"
			unknown++
		}
		rows = append(rows, row)
	}

	data := siteData{ID: site.ID, URL: site.URL, FormattedURL: u.Hostname(), Interval: site.Interval.String(), IntervalError: r.URL.Query().Get("intervalError") != "", Total: len(checks), Success: success, Failure: failure, Unknown: unknown, Checks: rows}

	var buf bytes.Buffer
	if err := s.tmpl.ExecuteTemplate(&buf, "site.html", data); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

func (s *Server) updateInterval(w http.ResponseWriter, r *http.Request) {
	siteID, err := uuid.Parse(r.PathValue("site-id"))
	if err != nil {
		http.Error(w, "invalid site id", http.StatusBadRequest)
		return
	}

	interval, err := time.ParseDuration(r.FormValue("interval"))
	if err != nil || interval < 10*time.Second || interval%time.Second != 0 || interval/time.Second > (1<<31)-1 {
		http.Redirect(w, r, "/site/"+siteID.String()+"?intervalError=1", http.StatusSeeOther)
		return
	}

	site, err := s.siteSvc.Get(r.Context(), siteID)
	if err != nil {
		if errors.Is(err, errors2.NotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	site.Interval = interval
	site.NextCheckAt = time.Now().Add(interval)
	if err := s.siteSvc.Update(r.Context(), site); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/site/"+siteID.String(), http.StatusSeeOther)
}

func (s *Server) deleteSite(w http.ResponseWriter, r *http.Request) {
	siteID, err := uuid.Parse(r.PathValue("site-id"))
	if err != nil {
		http.Error(w, "invalid site id", http.StatusBadRequest)
		return
	}
	if err := s.siteSvc.Delete(r.Context(), siteID); err != nil {
		if errors.Is(err, errors2.NotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

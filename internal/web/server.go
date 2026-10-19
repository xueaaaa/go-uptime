package web

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	service2 "github.com/xueaaaa/go-uptime/internal/check/service"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

type Server struct {
	addr     string
	tmpl     *template.Template
	siteSvc  service.SiteService
	checkSvc service2.CheckService
}

func NewServer(addr string, siteSvc service.SiteService, checkSvc service2.CheckService) *Server {
	return &Server{
		addr:     addr,
		tmpl:     template.Must(template.ParseFS(templatesFS, "templates/*.html")),
		siteSvc:  siteSvc,
		checkSvc: checkSvc,
	}
}

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()

	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(static)))
	r.Get("/", s.index)
	r.Route("/events", func(r chi.Router) {
		r.Get("/sites", s.sitesEvents)
	})

	return r
}

func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{Addr: s.addr, Handler: s.routes()}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

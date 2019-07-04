package server

import (
	"github.com/asahasrabuddhe/rest-api/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

var r chi.Router
var log *logrus.Logger

func init() {
	log = logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		PrettyPrint:      true,
	})

	r = chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logger.NewStructuredLogger(log))
	r.Use(middleware.Compress(1))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
}

func Mount(pattern string, router chi.Router) {
	r.Mount(pattern, router)
}

func Serve() {
	log.Fatal(http.ListenAndServe(":8080", r))
}

func AddLogHook(hook logrus.Hook) {
	log.AddHook(hook)
}

func ReplaceLogHook(hooks logrus.LevelHooks) logrus.LevelHooks {
	return log.ReplaceHooks(hooks)
}

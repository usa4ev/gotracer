package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	"github.com/usa4ev/gotracer/internal/model"
	"github.com/usa4ev/gotracer/internal/resources"
)

type (
	ctxKey string
	provider interface{
		Count(date time.Time) (map[string][]model.Entry, error)
	}
)

const keyResource = ctxKey("resource")

type (
	tracker interface {
		TrackSlowest() (resources.Resource, time.Duration, error)
		TrackFastest() (resources.Resource, time.Duration, error)
		TrackCall(resource resources.Resource) (time.Duration, error)
	}

	router struct {
		*chi.Mux
		tracker tracker
		provider provider
	}
)

func New(t tracker, l *logrus.Logger, p provider) *router {
	r := &router{
		Mux:     chi.NewRouter(),
		tracker: t,
		provider: p,
	}

	r.Use(middleware.RequestID)
	r.Use(NewStructuredLogger(l))

	r.With(r.MwValuesToCtx).Method(http.MethodGet,"/track", http.HandlerFunc(r.trackOne))
	r.Handle("/fastest", http.HandlerFunc(r.trackFastest))
	r.Handle("/slowest", http.HandlerFunc(r.trackSlowest))
	r.Handle("/stats", http.HandlerFunc(r.stats))
	return r
}

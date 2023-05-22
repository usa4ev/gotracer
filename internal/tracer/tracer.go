package tracer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/usa4ev/gotracer/internal/resources"
	"golang.org/x/sync/singleflight"
)

type (
	tracer struct {
		sfg       *singleflight.Group
		resources []resources.Resource
	}
)

func New(resources []resources.Resource) *tracer {

	return &tracer{
		sfg:       new(singleflight.Group),
		resources: resources}
}

type track struct {
	idx int
	time.Duration
	err error
}

func (t *tracer) trackAll() ([]track, error) {

	res := make([]track, len(t.resources))

	wg := sync.WaitGroup{}
	wg.Add(len(t.resources))

	for i, v := range t.resources {

		go func(idx int, resource resources.Resource) {
			defer wg.Done()

			d, err := t.trackOne(resources.Resource(resource))
			res[idx] = track{idx, d, err}
		}(i, v)
	}

	wg.Wait()

	return res, nil
}

// sortTracks sorts slice in duration ascending order
func sortTracks(t []track) {
	sort.Slice(t, func(i, j int) bool {
		return t[i].Duration < t[j].Duration
	})
}

func (t *tracer) TrackSlowest() (resources.Resource, time.Duration, error) {
	res, _ := t.trackAll()

	sortTracks(res)

	for i := len(t.resources)-1; i >= 0; i--{
		if res[i].err == nil{

			return t.resources[res[i].idx], res[i].Duration, nil
		}
	}
	
	return "", 0, fmt.Errorf("no resource is available")
}

func (t *tracer) TrackFastest() (resources.Resource, time.Duration, error) {
	res, _ := t.trackAll()

	sortTracks(res)

	for i := 0; i < len(t.resources); i++{
		if res[i].err == nil{

			return t.resources[res[i].idx], res[i].Duration, nil
		}
	}
	
	return "", 0, fmt.Errorf("no resource is available")
}

func (t *tracer) TrackCall(resource resources.Resource) (time.Duration, error) {
	res, err, _ := t.sfg.Do(string(resource), func() (any, error) {
		return t.trackOne(resource)
	})
	if err != nil {

		return 0, err
	}

	return res.(time.Duration), nil
}

func (t *tracer) trackOne(resource resources.Resource) (d time.Duration, err error) {
	var start time.Time

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var url string
	if strings.HasPrefix("http://", string(resource)) || 
			strings.HasPrefix("https://", string(resource)){
		url = string(resource)
	}else{
		url = fmt.Sprintf("https://%s", resource)
	}

	req, err := http.NewRequestWithContext(ctxTimeout, "GET", url, nil)
	if err != nil {

		return
	}

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			d = time.Since(start)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	_, err = http.DefaultTransport.RoundTrip(req)

	return
}

package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/usa4ev/gotracer/internal/charts"
	"github.com/usa4ev/gotracer/internal/resources"
)


type trackResponse struct{
	Resource string `json:"url"`
	TTFB string `json:"ttfb"`
	Error string`json:"error"`
}

func (rr * router) trackOne(w http.ResponseWriter, r *http.Request){
	resource := r.Context().Value(keyResource)
	if resource == nil{
		http.Error(w, "context missing resource", http.StatusInternalServerError)

		return
	}

	d,err := rr.tracker.TrackCall(resources.Resource(resource.(string)))

	res := trackResponse{
		Resource: resource.(string),
		TTFB: d.String(),
	}
	
	if err != nil{
		res.Error = err.Error()
	}

	enc := json.NewEncoder(w)

	if err := enc.Encode(res); err != nil {
		http.Error(
			w, 
			fmt.Sprintf("failed to encode message: %v", err.Error()), 
			http.StatusInternalServerError)

		return
	}
}

func (rr * router) trackFastest(w http.ResponseWriter, r *http.Request){
	resource, d, err := rr.tracker.TrackFastest()

	res := trackResponse{
		Resource: string(resource),
		TTFB: d.String(),
	}
	
	if err != nil{
		res.Error = err.Error()
	}

	enc := json.NewEncoder(w)

	if err := enc.Encode(res); err != nil {
		http.Error(
			w, 
			fmt.Sprintf("failed to encode message: %v", err.Error()), 
			http.StatusInternalServerError)

		return
	}
}

func (rr * router) trackSlowest(w http.ResponseWriter, r *http.Request){
	resource, d, err := rr.tracker.TrackSlowest()
	res := trackResponse{
		Resource: string(resource),
		TTFB: d.String(),
	}
	
	if err != nil{
		res.Error = err.Error()
	}

	enc := json.NewEncoder(w)

	if err := enc.Encode(res); err != nil {
		http.Error(
			w, 
			fmt.Sprintf("failed to encode message: %v", err.Error()), 
			http.StatusInternalServerError)

		return
	}
}

func (rr * router) stats(w http.ResponseWriter, r *http.Request){
	var date time.Time

	year, month, day := time.Now().Date()
    
	param := r.URL.Query().Get("date")
	if param == ""{
		date = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())
	}else{
		
		if parsed, err := time.Parse("20060102",param); err == nil{
			date = parsed 
		}else{
			http.Error(w,
				fmt.Sprintf("bad query parameter: %v, expected date YYYYMMDD", param),
				http.StatusBadRequest,
				)

			return
		}
	}

	res, err := rr.provider.Count(date)
	if err != nil{
		http.Error(w,
			fmt.Sprintf("failed to get data: %v", err),
			http.StatusInternalServerError,
			)

		return
	}	

	title := charts.NewTitle()
	title.Title = "Hourly request rate"
	title.Subtitle = fmt.Sprintf("Filtered by date: %s", date.Format("2006.01.02"))
	
	charts.DrawChart(w, res, title)
}
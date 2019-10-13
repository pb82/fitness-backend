package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pb82/fitness-backend/model"
	"github.com/pb82/fitness-backend/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	WorkoutSpeed     = "speed"
	WorkoutHeartrate = "heartrate"
	WorkoutAltitude  = "altitude"
)

var port string
var path string
var db store.WorkoutStore

// Probe endpoint: let grafana know that the backend is responding
func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Return the time series that this backend can return
func Search(w http.ResponseWriter, r *http.Request) {
	timeseries := []string{WorkoutHeartrate, WorkoutSpeed, WorkoutAltitude}
	json.NewEncoder(w).Encode(timeseries)
}

// Return the keys that can be used for filters
func TagKeys(w http.ResponseWriter, r *http.Request) {
	keys := []model.GrafanaTagKey{
		{
			Type: "string",
			Text: model.GrafanaWorkoutKey,
		},
	}
	json.NewEncoder(w).Encode(keys)
}

// Return the values that can be used for filters
func TagValues(w http.ResponseWriter, r *http.Request) {
	var tags []model.GrafanaTagValue
	for _, workout := range db.Workouts {
		tags = append(tags, model.GrafanaTagValue{Text: workout.ToHuman()})
	}
	json.NewEncoder(w).Encode(tags)
}

func QueryHeartrate(workout *model.Workout, request *model.GrafanaRequest) *model.GrafanaResponse {
	response := model.GrafanaResponse{
		Target:     WorkoutHeartrate,
		Datapoints: []model.Datapoint{},
	}

	for _, heartrate := range workout.Heartrate {
		if model.MatchesTime(heartrate.Timestamp, request.Range.FromMillis(), request.Range.ToMillis()) {
			response.Datapoints = append(response.Datapoints, model.Datapoint{
				float64(heartrate.Bpm), float64(heartrate.Timestamp),
			})
		}
	}

	log.Printf("QueryHeartrate: returning %v result(s)", len(response.Datapoints))
	return &response
}

func QuerySpeed(workout *model.Workout, request *model.GrafanaRequest) *model.GrafanaResponse {
	response := model.GrafanaResponse{
		Target:     WorkoutSpeed,
		Datapoints: []model.Datapoint{},
	}

	for _, location := range workout.Location {
		if model.MatchesTime(location.Timestamp, request.Range.FromMillis(), request.Range.ToMillis()) {
			response.Datapoints = append(response.Datapoints, model.Datapoint{
				float64(location.Speed), float64(location.Timestamp),
			})
		}
	}

	log.Printf("QuerySpeed: returning %v result(s)", len(response.Datapoints))
	return &response
}

func QueryAltitude(workout *model.Workout, request *model.GrafanaRequest) *model.GrafanaResponse {
	response := model.GrafanaResponse{
		Target:     WorkoutAltitude,
		Datapoints: []model.Datapoint{},
	}

	for _, location := range workout.Location {
		if model.MatchesTime(location.Timestamp, request.Range.FromMillis(), request.Range.ToMillis()) {
			response.Datapoints = append(response.Datapoints, model.Datapoint{
				float64(location.Altitude), float64(location.Timestamp),
			})
		}
	}

	log.Printf("QueryAltitude: returning %v result(s)", len(response.Datapoints))
	return &response
}

func Query(w http.ResponseWriter, r *http.Request) {
	if db.Empty() {
		json.NewEncoder(w).Encode([]string{})
		return
	}

	request := model.GrafanaRequest{}
	json.NewDecoder(r.Body).Decode(&request)

	workout := db.Filter(request.AdhocFilters)
	if workout == nil {
		json.NewEncoder(w).Encode([]string{})
		return
	}

	var response []*model.GrafanaResponse
	for _, target := range request.Targets {
		switch (target.Target) {
		case WorkoutHeartrate:
			response = append(response, QueryHeartrate(workout, &request))
		case WorkoutSpeed:
			response = append(response, QuerySpeed(workout, &request))
		case WorkoutAltitude:
			response = append(response, QueryAltitude(workout, &request))
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	json.NewEncoder(w).Encode(response)
}

// Upload new workout
func Push(w http.ResponseWriter, r *http.Request) {
	workout := model.Workout{}
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, existingWorkout := range db.Workouts {
		if workout.Timestamp == existingWorkout.Timestamp {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}

	db.Add(&workout)
	log.Printf("imported workout %s", workout.ToHuman())
}

func init() {
	flag.StringVar(&port, "port", "3000", "http server port")
	flag.StringVar(&path, "path", "./", "upload directory")
	flag.Parse()
}

func main() {
	r := mux.NewRouter()

	// Grafana json datastore api
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/search", Search).Methods("GET", "POST")
	r.HandleFunc("/query", Query).Methods("POST", "GET")
	r.HandleFunc("/tag-keys", TagKeys).Methods("POST", "GET")
	r.HandleFunc("/tag-values", TagValues).Methods("POST", "GET")

	// Push new workouts
	r.HandleFunc("/push", Push).Methods("POST")

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"accep", "content-type"}),
	)
	r.Use(cors)

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	server.Shutdown(ctx)
	os.Exit(0)
}

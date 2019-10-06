package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Heartrate struct {
	Timestamp int `json:"timestamp"`
	Bpm       int `json:"bpm"`
}

type Location struct {
	Timestamp int     `json:"timestamp"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Speed     float64 `json:"speed"`
}

type Workout struct {
	Timestamp int         `json:"timestamp"`
	Heartrate []Heartrate `json:"heartrate"`
	Location  []Location  `json:"location"`
}

type GrafanaRequestTarget struct {
	Target string `json:"target"`
	RefId  string `json:"refId"`
	Type   string `json:"type"`
}

type GrafanaRequest struct {
	Targets []GrafanaRequestTarget `json:"targets"`
}

type GrafanaResponse struct {
	Target     string      `json:"target"`
	Datapoints [][]float32 `json:"datapoints"`
}

var workout = Workout{
	Timestamp: int(time.Now().Unix()),
	Heartrate: []Heartrate{
		{
			Timestamp: int(time.Now().Unix()),
			Bpm:       100,
		},
		{
			Timestamp: int(time.Now().Unix()) + 1000,
			Bpm:       105,
		},
		{
			Timestamp: int(time.Now().Unix()) + 2000,
			Bpm:       90,
		},
	},
	Location: []Location{
		{
			Timestamp: int(time.Now().Unix()),
			Lat:       0,
			Lon:       0,
			Speed:     10,
		},
	},
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func Search(w http.ResponseWriter, r *http.Request) {
	timeseries := []string{"heartrate", "location", "speed"}
	json.NewEncoder(w).Encode(timeseries)
}

func Query(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var q GrafanaRequest
	err := decoder.Decode(&q)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response []GrafanaResponse
	for _, target := range q.Targets {
		if target.Target == "heartrate" {
			heartrateResponse := GrafanaResponse{
				Target: "heartrate",
			}

			for _, heartrate := range workout.Heartrate {
				heartrateResponse.Datapoints = append(heartrateResponse.Datapoints, []float32{
					float32(heartrate.Bpm), float32(heartrate.Timestamp),
				})
			}

			response = append(response, heartrateResponse)
		}
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "http server port")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/search", Search).Methods("GET", "POST")
	r.HandleFunc("/query", Query).Methods("POST")

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"accep", "content-type"}),
	)

	r.Use(cors)

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
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

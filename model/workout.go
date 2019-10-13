package model

import (
	"fmt"
	"time"
)

var timezone, _ = time.LoadLocation("Europe/Berlin")

type Heartrate struct {
	Timestamp int64 `json:"timestamp"`
	Bpm       int   `json:"bpm"`
}

type Location struct {
	Timestamp int64   `json:"timestamp"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Speed     float64 `json:"speed"`
	Altitude  float64 `json:"altitude"`
}

type Workout struct {
	Timestamp int64       `json:"timestamp"`
	Heartrate []Heartrate `json:"heartrate"`
	Location  []Location  `json:"location"`
}

func (w *Workout) ToHuman() string {
	t := time.Unix(0, w.Timestamp*int64(time.Millisecond)).In(timezone)
	return fmt.Sprintf("%v/%v/%v (%v:%.2v)", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func (w *Workout) MatchesFilter(filter *GrafanaFilter) bool {
	return w.ToHuman() == filter.Value
}

func MatchesTime(test, from, to int64) bool {
	return test >= from && test <= to
}

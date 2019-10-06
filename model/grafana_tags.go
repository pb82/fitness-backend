package model

const (
	GrafanaWorkoutKey = "Workout"
)

type GrafanaTagKey struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type GrafanaTagValue struct {
	Text string `json:"text"`
}

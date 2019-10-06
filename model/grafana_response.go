package model

type Datapoint []float32

type GrafanaResponse struct {
	Target     string      `json:"target"`
	Datapoints []Datapoint `json:"datapoints"`
}

func (r *GrafanaResponse) AddDatapoint(metric, timestamp float32) {
	r.Datapoints = append(r.Datapoints, Datapoint{metric, timestamp})
}

func ResponseForTarget(target string) *GrafanaResponse {
	return &GrafanaResponse{
		Target:     target,
		Datapoints: []Datapoint{},
	}
}

package model

import "time"

const (
	TimeLayout = "2006-01-02T15:04:05.000Z"
)

type GrafanaRequestTarget struct {
	Target string `json:"target"`
	RefId  string `json:"refId"`
	Type   string `json:"type"`
}

type GrafanaRequestRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type GrafanaRequest struct {
	Range        GrafanaRequestRange    `json:"range"`
	Targets      []GrafanaRequestTarget `json:"targets"`
	AdhocFilters []GrafanaFilter        `json:"adhocFilters"`
}

type GrafanaFilter struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

func toMillis(timestamp string) int64 {
	t, err := time.Parse(TimeLayout, timestamp)
	if err != nil {
		return -1
	}

	return t.UnixNano() / int64(time.Millisecond)
}

func (r *GrafanaRequestRange) FromMillis() int64 {
	return toMillis(r.From)
}

func (r *GrafanaRequestRange) ToMillis() int64 {
	return toMillis(r.To)
}

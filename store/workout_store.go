package store

import "github.com/pb82/fitness-backend/model"

type WorkoutStore struct {
	Workouts []*model.Workout
}

func (s *WorkoutStore) Add(w *model.Workout) {
	s.Workouts = append(s.Workouts, w)
}

func (s *WorkoutStore) Filter(filters []model.GrafanaFilter) *model.Workout {
	for _, workout := range s.Workouts {
		for _, filter := range filters {
			if !workout.MatchesFilter(&filter) {
				continue
			}
		}
		return workout
	}

	return nil
}

func (s *WorkoutStore) Empty() bool {
	return len(s.Workouts) == 0
}

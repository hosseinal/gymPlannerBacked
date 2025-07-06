package models

type PlanDetails struct {
	PlanID     int64  `json:"plan_id" db:"plan_id"`
	Move       string `json:"move" db:"move"`
	NumOfReps  int    `json:"num_of_reps" db:"num_of_reps"`
} 
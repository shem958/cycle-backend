package models

import "gorm.io/gorm"

type Cycle struct {
	gorm.Model
	StartDate string `json:"startDate" binding:"required"`     // required
	Length    int    `json:"length" binding:"required,min=1"`  // required, at least 1
	Mood      string `json:"mood"`                             // optional
	Symptoms  string `json:"symptoms"`                         // optional
}

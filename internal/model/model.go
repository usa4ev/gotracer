package model

import "time"

// Entry represents single query rate entry
type Entry struct{
	Time time.Time
	Count int
}
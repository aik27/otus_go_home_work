package storage

import "time"

type Event struct {
	Id          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	UserId      int64     `db:"user_id"`
	RemindIn    string    `db:"remind_in"`
}

type EventList []Event

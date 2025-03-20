package sqlstorage

import (
	"context"
	"log"

	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db  *sqlx.DB
	ctx context.Context
}

func New(dsn string) (*Storage, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	err := s.db.PingContext(ctx)
	if err != nil {
		return err
	}

	s.ctx = ctx

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	err := s.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Save(event storage.Event) (storage.Event, error) {
	if event.Id == 0 {
		event, err := s.create(event)
		return event, err
	}

	err := s.update(event)

	return event, err
}

func (s *Storage) Delete(event storage.Event) error {
	_, err := s.db.ExecContext(s.ctx, "DELETE FROM events where id=$1", event.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetByUserId(userId int64) (storage.EventList, error) {
	events := make(storage.EventList, 0)
	params := map[string]interface{}{
		"user_id": userId,
	}

	rows, err := s.db.NamedQueryContext(s.ctx, "SELECT * FROM events WHERE user_id=:user_id", params)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	for rows.Next() {
		var event storage.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) GetAll() (storage.EventList, error) {
	events := make(storage.EventList, 0)
	params := map[string]interface{}{}

	rows, err := s.db.NamedQueryContext(s.ctx, `SELECT * FROM events`, params)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	for rows.Next() {
		var event storage.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) create(event storage.Event) (storage.Event, error) {
	err := s.db.QueryRowContext(
		s.ctx,
		`INSERT INTO events(title, description, start_date, end_date, remind_in, user_id) 
				VALUES($1,$2,$3,$4,$5,$6) RETURNING id`,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.RemindIn,
		event.UserId,
	).Scan(&event.Id)
	if err != nil {
		return event, err
	}

	return event, nil
}

func (s *Storage) update(event storage.Event) error {
	params := map[string]interface{}{
		"id":          event.Id,
		"title":       event.Title,
		"description": event.Description,
		"start_date":  event.StartDate,
		"end_date":    event.EndDate,
		"remind_in":   event.RemindIn,
		"user_id":     event.UserId,
	}
	_, err := s.db.NamedExecContext(
		s.ctx,
		`UPDATE events SET 
                  title=:title, 
                  description=:description, 
                  start_date=:start_date, 
                  end_date=:end_date, 
                  remind_in=:remind_in, 
                  user_id=:user_id 
              WHERE id=:id`,
		params,
	)
	if err != nil {
		return err
	}

	return nil
}

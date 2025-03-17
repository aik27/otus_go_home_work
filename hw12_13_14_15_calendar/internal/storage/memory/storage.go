package memorystorage

import (
	"context"
	"sync"

	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[int64]storage.Event
	mu     sync.RWMutex
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func New() *Storage {
	return &Storage{
		events: map[int64]storage.Event{},
		mu:     sync.RWMutex{},
	}
}

func (s *Storage) Save(event storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.Id == 0 {
		event.Id = int64(len(s.events) + 1)
	}

	s.events[event.Id] = event

	return event, nil
}

func (s *Storage) Delete(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, event.Id)

	return nil
}

func (s *Storage) GetById(id int64) (storage.Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.events[id]; ok {
		return s.events[id], true
	}

	return storage.Event{}, false
}

func (s *Storage) GetByUserId(userId int64) (storage.EventList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(storage.EventList, 0)

	for _, ev := range s.events {
		if ev.UserId == userId {
			result = append(result, ev)
		}
	}

	return result, nil
}

func (s *Storage) GetAll() (storage.EventList, error) {
	eventList := make(storage.EventList, 0, len(s.events))

	for _, v := range s.events {
		eventList = append(eventList, v)
	}

	return eventList, nil
}

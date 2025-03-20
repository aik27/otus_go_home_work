package memorystorage

import (
	"testing"
	"time"

	"github.com/aik27/otus_go_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	event1 := storage.Event{
		Title:       "Test Event",
		Description: "Test",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Second * 10),
		UserId:      1,
		RemindIn:    "1",
	}
	event2 := storage.Event{
		Title:       "Test Event",
		Description: "Test",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Second * 10),
		UserId:      1,
		RemindIn:    "1",
	}
	event3 := storage.Event{
		Title:       "Test Event",
		Description: "Test",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Second * 10),
		UserId:      1,
		RemindIn:    "1",
	}

	t.Run("Save and GetById", func(t *testing.T) {
		s := New()
		savedEvent, err := s.Save(event1)
		require.NoError(t, err)
		require.NotZero(t, savedEvent.Id)

		fetchedEvent, found := s.GetById(savedEvent.Id)
		require.True(t, found)
		require.Equal(t, savedEvent, fetchedEvent)
	})

	t.Run("Delete", func(t *testing.T) {
		s := New()
		savedEvent, err := s.Save(event1)
		require.NoError(t, err)

		err = s.Delete(savedEvent)
		require.NoError(t, err)

		_, found := s.GetById(savedEvent.Id)
		require.False(t, found)
	})

	t.Run("GetByUserId", func(t *testing.T) {
		s := New()
		_, err := s.Save(event1)
		require.NoError(t, err)

		_, err = s.Save(event2)
		require.NoError(t, err)

		_, err = s.Save(event3)
		require.NoError(t, err)

		events, err := s.GetByUserId(1)
		require.NoError(t, err)
		require.Len(t, events, 3)
	})

	t.Run("GetAll", func(t *testing.T) {
		s := New()

		_, err := s.Save(event1)
		require.NoError(t, err)

		_, err = s.Save(event2)
		require.NoError(t, err)

		events, err := s.GetAll()
		require.NoError(t, err)
		require.Len(t, events, 2)
	})
}

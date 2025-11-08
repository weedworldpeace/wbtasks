package repository

import (
	"calendar/internal/models"
	"calendar/pkg/data"
	"slices"
	"testing"
	"time"
)

const (
	defaultNotCreatedId = "-1"
	defaultMessage      = "ok its the hardest... i swear to god"
	defaultDate         = "2006-01-02"
)

type testCase struct {
	input    any
	expected any
}

func testCreateEvent(repo *Repository, t *testing.T) {
	defaultId := "1"
	testCases := []testCase{
		{input: &models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: defaultDate}},
			expected: nil,
		},
	}

	for i := range testCases {
		res := repo.CreateEvent(testCases[i].input.(*models.UserEvent))
		if res != testCases[i].expected {
			t.Fail()
		}
	}
}

func testUpdateEvent(repo *Repository, t *testing.T) {
	defaultId := "2"
	testCases := []testCase{
		{input: &models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: defaultDate}},
			expected: nil,
		},
		{input: &models.UserEvent{UserId: defaultNotCreatedId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: defaultDate}},
			expected: ErrNonExistUserId,
		},
		{input: &models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultNotCreatedId, Message: defaultMessage, Date: defaultDate}},
			expected: ErrNonExistEventId,
		},
	}

	err := repo.CreateEvent(&models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: defaultDate}})
	if err != nil {
		t.Fatal()
	}

	for i := range testCases {
		res := repo.UpdateEvent(testCases[i].input.(*models.UserEvent))
		if res != testCases[i].expected {
			t.Fail()
		}
	}
}

func testDeleteEvent(repo *Repository, t *testing.T) {
	defaultId := "3"
	testCases := []testCase{
		{input: &models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: defaultDate}},
			expected: nil,
		},
		{input: &models.UserEvent{UserId: defaultNotCreatedId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: defaultDate}},
			expected: ErrNonExistUserId,
		},
		{input: &models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultNotCreatedId, Message: defaultMessage, Date: defaultDate}},
			expected: ErrNonExistEventId,
		},
	}

	for i := range testCases {
		err := repo.CreateEvent(&models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: defaultDate}})
		if err != nil {
			t.Fatal()
		}
		res := repo.DeleteEvent(testCases[i].input.(*models.UserEvent))
		if res != testCases[i].expected {
			t.Fail()
		}
	}
}

type costyl1 struct {
	date string
	from int64
	to   int64
}
type costyl2 struct {
	result []models.Event
	err    error
}

func testReadEvents(repo *Repository, t *testing.T) {
	defaultId := "4"

	firstDate := "2006-01-02"
	firstEvent := models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: firstDate}}
	secondDate := "2006-01-08"
	secondEvent := models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: secondDate}}
	thirdDate := "2006-02-01"
	thirdEvent := models.UserEvent{UserId: defaultId, Event: models.Event{EventId: defaultId, Message: defaultMessage, Date: thirdDate}}

	firstDateTime, err := time.Parse("2006-01-02", firstDate)
	if err != nil {
		t.Fatal()
	}

	testCases := []testCase{
		{input: costyl1{defaultNotCreatedId, 0, 0},
			expected: costyl2{[]models.Event{}, ErrNonExistUserId},
		},
		{input: costyl1{defaultId, firstDateTime.Unix(), firstDateTime.AddDate(0, 0, 1).Unix()},
			expected: costyl2{[]models.Event{firstEvent.Event}, nil},
		},
		{input: costyl1{defaultId, firstDateTime.Unix(), firstDateTime.AddDate(0, 0, 7).Unix()},
			expected: costyl2{[]models.Event{firstEvent.Event, secondEvent.Event}, nil},
		},
		{input: costyl1{defaultId, firstDateTime.Unix(), firstDateTime.AddDate(0, 1, 0).Unix()},
			expected: costyl2{[]models.Event{firstEvent.Event, secondEvent.Event, thirdEvent.Event}, nil},
		},
	}

	err = repo.CreateEvent(&firstEvent)
	if err != nil {
		t.Fatal()
	}
	err = repo.CreateEvent(&secondEvent)
	if err != nil {
		t.Fatal()
	}
	err = repo.CreateEvent(&thirdEvent)
	if err != nil {
		t.Fatal()
	}

	for i := range testCases {
		res, err := repo.ReadEvents(testCases[i].input.(costyl1).date, testCases[i].input.(costyl1).from, testCases[i].input.(costyl1).to)
		if err != testCases[i].expected.(costyl2).err || !slices.Equal(res, testCases[i].expected.(costyl2).result) {
			t.Fail()
		}
	}
}

func TestMain(t *testing.T) {
	repo := New(data.New())

	t.Run("test CreateEvent", func(t *testing.T) { // change mb
		t.Parallel()
		testCreateEvent(repo, t)
	})

	t.Run("test UpdateEvents", func(t *testing.T) {
		t.Parallel()
		testUpdateEvent(repo, t)
	})

	t.Run("test DeleteEvents", func(t *testing.T) {
		t.Parallel()
		testDeleteEvent(repo, t)
	})

	t.Run("test ReadEvents", func(t *testing.T) {
		t.Parallel()
		testReadEvents(repo, t)
	})

}

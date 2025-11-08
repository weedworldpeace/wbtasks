package service

import (
	"calendar/internal/mocks"
	"calendar/internal/models"
	"errors"
	"slices"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestValidUserId(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    string
		expected error
	}{
		{input: "invalid user id", expected: ErrInvalidUserId},
		{input: uuid.NewString(), expected: nil},
	}

	for _, v := range testCases {
		err := validUserId(v.input)
		if !errors.Is(err, v.expected) {
			t.Logf("input: %s, expected:%v, got: %v", v.input, v.expected, err)
			t.Fail()
		}
	}
}

func TestValidEventId(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    string
		expected error
	}{
		{input: "invalid event id", expected: ErrInvalidEventId},
		{input: uuid.NewString(), expected: nil},
	}

	for _, v := range testCases {
		err := validEventId(v.input)
		if !errors.Is(err, v.expected) {
			t.Logf("input: %s, expected:%v, got: %v", v.input, v.expected, err)
			t.Fail()
		}
	}
}

func TestValidDate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    string
		expected error
	}{
		{input: "invalid date", expected: ErrInvalidDate},
		{input: "2006-01-03", expected: nil},
	}

	for _, v := range testCases {
		err := validDate(v.input)
		if !errors.Is(err, v.expected) {
			t.Logf("input: %s, expected:%v, got: %v", v.input, v.expected, err)
			t.Fail()
		}
	}
}

func TestReadEvents(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	srv := New(mocks.NewMockRepositoryInterface(ctrl))

	type input struct {
		userId  string
		rawDate string
		genre   string
	}
	type expected struct {
		events []models.Event
		err    error
	}
	type testCase struct {
		input    input
		expected expected
	}

	testCases := []testCase{
		{input: input{userId: uuid.NewString(), rawDate: "2006-01-03", genre: "invalid time period"}, expected: expected{events: []models.Event{}, err: ErrInvalidTimePeriod}},
	}

	for _, v := range testCases {
		res, err := srv.ReadEvents(v.input.userId, v.input.rawDate, v.input.genre)
		if !errors.Is(err, v.expected.err) || !slices.Equal(res, v.expected.events) {
			t.Logf("input: %v, expected:%v, got: %v, %v", v.input, v.expected, res, err)
			t.Fail()
		}
	}
}

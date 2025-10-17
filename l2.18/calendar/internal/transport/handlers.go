package transport

import (
	"calendar/internal/models"
	"calendar/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	errBadMethod = errors.New("bad method")
	errNoUserId  = errors.New("no user id")
	errNoDate    = errors.New("no date")
)

type handlers struct {
	ctx     context.Context
	service ServiceInterface
}

func (h *handlers) writeBadResponse(w http.ResponseWriter, r *http.Request, err error) {
	lg := logger.LoggerFromCtx(h.ctx).Lg
	lg.Error(err.Error())

	res, err := json.Marshal(models.NewBadResponse(err))
	if err != nil {
		lg.Error(err.Error())
	}
	_, err = w.Write(res)
	if err != nil {
		lg.Error(err.Error())
	}
}

func (h *handlers) writeGoodPostResponse(w http.ResponseWriter, r *http.Request, message string) {
	lg := logger.LoggerFromCtx(h.ctx).Lg

	res, err := json.Marshal(models.NewGoodPostResponse(message))
	if err != nil {
		lg.Error(err.Error())
	}

	w.WriteHeader(200)
	_, err = w.Write(res)
	if err != nil {
		lg.Error(err.Error())
	}
}

func (h *handlers) writeGoodGetResponse(w http.ResponseWriter, r *http.Request, events []models.Event) {
	lg := logger.LoggerFromCtx(h.ctx).Lg

	res, err := json.Marshal(models.NewGoodGetResponse(events))
	if err != nil {
		lg.Error(err.Error())
	}

	_, err = w.Write(res)
	if err != nil {
		lg.Error(err.Error())
	}
}

func (h *handlers) middleware(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := logger.LoggerFromCtx(h.ctx).Lg
		lg.Info(fmt.Sprintf("received request %s %s", r.Method, r.URL))
		next.ServeHTTP(w, r)
	}
}

func (h *handlers) createEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeBadResponse(w, r, errBadMethod)
		return
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	userEvent := models.NewUserEvent()

	err = json.Unmarshal(raw, userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	err = h.service.CreateEvent(userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}
	h.writeGoodPostResponse(w, r, "event created successfully")
}

func (h *handlers) updateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeBadResponse(w, r, errBadMethod)
		return
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	userEvent := models.NewUserEvent()

	err = json.Unmarshal(raw, userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	err = h.service.UpdateEvent(userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodPostResponse(w, r, "event updated successfully")
}

func (h *handlers) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeBadResponse(w, r, errBadMethod)
		return
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	userEvent := models.NewUserEvent()

	err = json.Unmarshal(raw, userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	err = h.service.DeleteEvent(userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodPostResponse(w, r, "event deleted successfully")
}

func (h *handlers) eventsForDay(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if !values.Has("user_id") {
		h.writeBadResponse(w, r, errNoUserId)
		return
	}
	if !values.Has("date") {
		h.writeBadResponse(w, r, errNoDate)
		return
	}

	events, err := h.service.ReadEvents(values.Get("user_id"), values.Get("date"), "day")
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodGetResponse(w, r, events)
}

func (h *handlers) eventsForWeek(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if !values.Has("user_id") {
		h.writeBadResponse(w, r, errNoUserId)
		return
	}
	if !values.Has("date") {
		h.writeBadResponse(w, r, errNoDate)
		return
	}

	events, err := h.service.ReadEvents(values.Get("user_id"), values.Get("date"), "week")
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodGetResponse(w, r, events)
}

func (h *handlers) eventsForMonth(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if !values.Has("user_id") {
		h.writeBadResponse(w, r, errNoUserId)
		return
	}
	if !values.Has("date") {
		h.writeBadResponse(w, r, errNoDate)
		return
	}

	events, err := h.service.ReadEvents(values.Get("user_id"), values.Get("date"), "month")
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodGetResponse(w, r, events)
}

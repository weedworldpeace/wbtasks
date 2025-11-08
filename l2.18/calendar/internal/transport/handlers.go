package transport

import (
	"calendar/internal/models"
	"calendar/internal/repository"
	"calendar/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
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
	lg := logger.LoggerFromCtx(r.Context()).Lg
	lg.Error(err.Error())

	if errors.Is(err, errBadMethod) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else if errors.Is(err, repository.ErrNonExistEventId) || (errors.Is(err, repository.ErrNonExistUserId)) {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	res, err := json.Marshal(models.NewBadResponse(err))
	if err != nil {
		lg.Error(err.Error())
		return
	}

	_, err = w.Write(res)
	if err != nil {
		lg.Error(err.Error())
	}
}

func (h *handlers) writeGoodPostResponse(w http.ResponseWriter, r *http.Request, message string) {
	lg := logger.LoggerFromCtx(r.Context()).Lg

	res, err := json.Marshal(models.NewGoodPostResponse(message))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lg.Error(err.Error())
		return
	}

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lg.Error(err.Error())
	}
}

func (h *handlers) writeGoodGetResponse(w http.ResponseWriter, r *http.Request, events []models.Event) {
	lg := logger.LoggerFromCtx(r.Context()).Lg

	res, err := json.Marshal(models.NewGoodGetResponse(events))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lg.Error(err.Error())
		return
	}

	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lg.Error(err.Error())
	}
}

func (h *handlers) processPost(w http.ResponseWriter, r *http.Request) (*models.UserEvent, bool) {
	if r.Method != http.MethodPost {
		h.writeBadResponse(w, r, errBadMethod)
		return nil, false
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return nil, false
	}

	userEvent := models.NewUserEvent()

	err = json.Unmarshal(raw, userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return nil, false
	}

	return userEvent, true
}

func (h *handlers) processGet(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	if r.Method != http.MethodGet {
		h.writeBadResponse(w, r, errBadMethod)
		return "", "", false
	}

	values := r.URL.Query()
	if !values.Has("user_id") {
		h.writeBadResponse(w, r, errNoUserId)
		return "", "", false
	}
	if !values.Has("date") {
		h.writeBadResponse(w, r, errNoDate)
		return "", "", false
	}
	return values.Get("user_id"), values.Get("date"), true
}

func (h *handlers) middleware(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lg := logger.LoggerFromCtx(h.ctx)
		requestId := uuid.NewString()
		lgWithReqId := lg.LoggerWithRequestId(requestId)

		lgWithReqId.Lg.Info("received request", "method", r.Method, "url", r.URL.String())

		r = r.WithContext(context.WithValue(r.Context(), logger.LoggerKey, lgWithReqId))
		next.ServeHTTP(w, r)
	}
}

func (h *handlers) createEvent(w http.ResponseWriter, r *http.Request) {
	userEvent, toContinue := h.processPost(w, r)
	if !toContinue {
		return
	}

	err := h.service.CreateEvent(userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}
	h.writeGoodPostResponse(w, r, "event created successfully")
}

func (h *handlers) updateEvent(w http.ResponseWriter, r *http.Request) {
	userEvent, toContinue := h.processPost(w, r)
	if !toContinue {
		return
	}

	err := h.service.UpdateEvent(userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodPostResponse(w, r, "event updated successfully")
}

func (h *handlers) deleteEvent(w http.ResponseWriter, r *http.Request) {
	userEvent, toContinue := h.processPost(w, r)
	if !toContinue {
		return
	}

	err := h.service.DeleteEvent(userEvent)
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodPostResponse(w, r, "event deleted successfully")
}

func (h *handlers) eventsForDay(w http.ResponseWriter, r *http.Request) {
	userId, date, toContinue := h.processGet(w, r)
	if !toContinue {
		return
	}

	events, err := h.service.ReadEvents(userId, date, "day")
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodGetResponse(w, r, events)
}

func (h *handlers) eventsForWeek(w http.ResponseWriter, r *http.Request) {
	userId, date, toContinue := h.processGet(w, r)
	if !toContinue {
		return
	}

	events, err := h.service.ReadEvents(userId, date, "week")
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodGetResponse(w, r, events)
}

func (h *handlers) eventsForMonth(w http.ResponseWriter, r *http.Request) {
	userId, date, toContinue := h.processGet(w, r)
	if !toContinue {
		return
	}

	events, err := h.service.ReadEvents(userId, date, "month")
	if err != nil {
		h.writeBadResponse(w, r, err)
		return
	}

	h.writeGoodGetResponse(w, r, events)
}

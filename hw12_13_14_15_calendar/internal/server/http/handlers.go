package internalhttp

import (
	"encoding/json"
	"io"
	"net/http"
)

type Handler struct {
	app    Application
	logger Logger
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	event, err := decodeBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())

		return
	}

	err = h.app.Create(r.Context(), EventToModel(event))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	event, err := decodeBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	e, err := h.app.Get(r.Context(), EventToModel(event).ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	event, err := decodeBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	err = h.app.Update(r.Context(), EventToModel(event))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	event, err := decodeBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	err = h.app.Delete(r.Context(), EventToModel(event).ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) listDayEvents(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	event, err := decodeBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	e, err := h.app.ListDayEvents(r.Context(), EventToModel(event).StartDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) listWeekEvents(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	event, err := decodeBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	e, err := h.app.ListWeekEvents(r.Context(), EventToModel(event).StartDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) listMonthEvents(w http.ResponseWriter, r *http.Request) { //nolint:dupl
	event, err := decodeBody(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	e, err := h.app.ListMonthEvents(r.Context(), EventToModel(event).StartDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.logger.Error(err.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
}

func decodeBody(body io.ReadCloser) (Event, error) {
	dec := json.NewDecoder(body)
	e := Event{}

	err := dec.Decode(&e)
	if err != nil {
		return Event{}, err
	}

	return e, nil
}

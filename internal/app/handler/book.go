package handler

import (
	"encoding/json"
	"http-rest-api-go/internal/app/model"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) handleBooksCreate() http.HandlerFunc {
	type request struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			h.error(w, r, http.StatusBadRequest, err)
			return
		}

		b := &model.Book{
			Title:  req.Title,
			Author: req.Author,
		}
		if err := h.service.Create(b); err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusCreated, b)
	}
}

func (h *Handler) handleBooksGetAll() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		books, err := h.service.GetAll()

		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusOK, books)

	}
}

func (h *Handler) handleBooksGet() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		book, err := h.service.GetById(id)

		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusOK, book)

	}
}

func (h *Handler) handleBooksPut() http.HandlerFunc {
	type request struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		b := &model.UpdateBookInput{}

		if r.Body != nil {
			jsonData, _ := io.ReadAll(r.Body)

			err := json.Unmarshal(jsonData, &b)
			if err != nil {
				// handle error
				h.error(w, r, http.StatusBadRequest, err)
				return
			}

		}

		if err := h.service.Update(id, b); err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusOK, nil)
	}
}

func (h *Handler) handleBooksDelete() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = h.service.Delete(id)
		if err != nil {
			h.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		h.respond(w, r, http.StatusOK, nil)
	}
}

func (h *Handler) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	h.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (h *Handler) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

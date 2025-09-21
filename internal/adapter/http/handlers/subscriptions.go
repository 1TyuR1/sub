package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"crud_ef/internal/domain"
	"crud_ef/internal/usecase/subscription"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubscriptionDTO struct {
	ID           uuid.UUID `json:"id"`
	ServiceName  string    `json:"service_name"`
	MonthlyPrice string    `json:"monthly_price"`
	UserID       uuid.UUID `json:"user_id"`
	StartMonth   string    `json:"start_month"`
	EndMonth     *string   `json:"end_month,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateRequest struct {
	ServiceName  string  `json:"service_name"`
	MonthlyPrice string  `json:"monthly_price"`
	UserID       string  `json:"user_id"`
	StartMonth   string  `json:"start_month"`
	EndMonth     *string `json:"end_month,omitempty"`
}

type UpdateRequest struct {
	ServiceName  *string `json:"service_name,omitempty"`
	MonthlyPrice *string `json:"monthly_price,omitempty"`
	StartMonth   *string `json:"start_month,omitempty"`
	EndMonth     *string `json:"end_month,omitempty"`
}

type SubscriptionRoutes struct {
	svc *subscription.Service
}

func NewSubscriptionRoutes(svc *subscription.Service) *SubscriptionRoutes {
	return &SubscriptionRoutes{svc: svc}
}

func (h *SubscriptionRoutes) Register(r chi.Router) {
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/{id}", h.get)
		r.Get("/", h.list)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// @Summary      Create subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request  body  CreateRequest  true  "payload"
// @Success      201  {object}  SubscriptionDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [post]
func (h *SubscriptionRoutes) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
		return
	}
	in := domain.CreateInput{
		ServiceName:  req.ServiceName,
		MonthlyPrice: req.MonthlyPrice,
		UserID:       uid,
		StartMonth:   req.StartMonth,
		EndMonth:     req.EndMonth,
	}
	s, err := h.svc.Create(r.Context(), in)
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid") {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	writeJSON(w, http.StatusCreated, toDTO(s))
}

// @Summary      Get subscription by id
// @Tags         subscriptions
// @Produce      json
// @Param        id   path  string  true  "Subscription ID"
// @Success      200  {object}  SubscriptionDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionRoutes) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	s, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, toDTO(s))
}

// @Summary      List subscriptions
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query  string  false  "Filter by user UUID"
// @Param        service_name  query  string  false  "Filter by service (ILIKE)"
// @Param        limit         query  int     false  "Limit (1..200)"
// @Param        offset        query  int     false  "Offset"
// @Success      200  {array}   SubscriptionDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions [get]
func (h *SubscriptionRoutes) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var f domain.ListFilter
	if v := q.Get("user_id"); v != "" {
		uid, err := uuid.Parse(v)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
			return
		}
		f.UserID = &uid
	}
	if v := q.Get("service_name"); v != "" {
		f.ServiceName = &v
	}
	limit := 50
	offset := 0
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	f.Limit = limit
	f.Offset = offset

	items, err := h.svc.List(r.Context(), f)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	out := make([]SubscriptionDTO, 0, len(items))
	for _, s := range items {
		out = append(out, toDTO(s))
	}
	writeJSON(w, http.StatusOK, out)
}

// @Summary      Update subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id       path   string          true  "Subscription ID"
// @Param        request  body   UpdateRequest   true  "payload"
// @Success      200  {object}  SubscriptionDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionRoutes) update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	in := domain.UpdateInput{
		ServiceName:  req.ServiceName,
		MonthlyPrice: req.MonthlyPrice,
		StartMonth:   req.StartMonth,
		EndMonth:     req.EndMonth,
	}
	s, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid") {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		}
		return
	}
	writeJSON(w, http.StatusOK, toDTO(s))
}

// @Summary      Delete subscription
// @Tags         subscriptions
// @Param        id  path  string  true  "Subscription ID"
// @Success      204  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionRoutes) delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	ok, err := h.svc.Delete(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusNoContent, map[string]string{"status": "deleted"})
}

func toDTO(s domain.Subscription) SubscriptionDTO {
	return SubscriptionDTO{
		ID:           s.ID,
		ServiceName:  s.ServiceName,
		MonthlyPrice: s.MonthlyPrice,
		UserID:       s.UserID,
		StartMonth:   s.StartMonth,
		EndMonth:     s.EndMonth,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

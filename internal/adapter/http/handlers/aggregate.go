package handlers

import (
	"net/http"

	"crud_ef/internal/usecase/subscription"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TotalResponse struct {
	From        string  `json:"from"`
	To          string  `json:"to"`
	UserID      *string `json:"user_id,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
	Total       string  `json:"total"`
}

type AggregateRoutes struct {
	svc *subscription.Service
}

func NewAggregateRoutes(svc *subscription.Service) *AggregateRoutes {
	return &AggregateRoutes{svc: svc}
}

func (h *AggregateRoutes) Register(r chi.Router) {
	r.Get("/subscriptions/total", h.total)
}

// @Summary      Total cost for period
// @Tags         subscriptions
// @Produce      json
// @Param        from          query  string  true   "YYYY-MM"
// @Param        to            query  string  true   "YYYY-MM"
// @Param        user_id       query  string  false  "User UUID"
// @Param        service_name  query  string  false  "Service filter (ILIKE)"
// @Success      200  {object}  TotalResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /subscriptions/total [get]
func (h *AggregateRoutes) total(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fromStr := q.Get("from")
	toStr := q.Get("to")
	if fromStr == "" || toStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "from and to are required (YYYY-MM)"})
		return
	}
	var uid *uuid.UUID
	if s := q.Get("user_id"); s != "" {
		u, err := uuid.Parse(s)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
			return
		}
		uid = &u
	}
	var service *string
	if s := q.Get("service_name"); s != "" {
		service = &s
	}

	total, err := h.svc.Total(r.Context(), fromStr, toStr, uid, service)
	if err != nil {
		if err.Error() == "to must be >= from" || err.Error()[:7] == "invalid" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}
	resp := TotalResponse{
		From:  fromStr,
		To:    toStr,
		Total: total,
	}
	if uid != nil {
		s := uid.String()
		resp.UserID = &s
	}
	if service != nil {
		resp.ServiceName = service
	}
	writeJSON(w, http.StatusOK, resp)
}

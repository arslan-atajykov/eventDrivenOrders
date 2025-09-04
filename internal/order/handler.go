package order

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo     *Repository
	producer *Producer
}

func NewHandler(repo *Repository, producer *Producer) *Handler {
	return &Handler{repo: repo, producer: producer}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var o Order
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&o); err != nil {
		http.Error(w, "invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation
	if o.Customer == "" {
		http.Error(w, "customer is required", http.StatusBadRequest)
		return
	}

	o.Status = "new"

	if err := h.repo.CreateOrder(r.Context(), &o); err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}
	if err := h.producer.PublishOrder(r.Context(), &o); err != nil {
		http.Error(w, "failed to publish even:"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

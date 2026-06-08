package deployments

import (
	"CloudHub/internal/types"
	"CloudHub/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/deployments", func(r chi.Router) {
		r.Post("/", h.handleNewDeployment)
	})
}
func (h *Handler) handleNewDeployment(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateNewDeploymentPayload

	err := utils.ParseJson(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	newDeployment, err := h.store.CreateNewDeployment(r.Context(), payload.GitUrl)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = utils.WriteJson(w, http.StatusCreated, newDeployment)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

package deployments

import (
	"CloudHub/internal/queue"
	"CloudHub/internal/types"
	"CloudHub/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	uuid2 "github.com/google/uuid"
)

type Handler struct {
	store      *Store
	RedisStore *queue.RedisStore
}

func NewHandler(store *Store, rdb *queue.RedisStore) *Handler {
	return &Handler{
		store:      store,
		RedisStore: rdb,
	}
}
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/deployments", func(r chi.Router) {
		r.Post("/", h.handleNewDeployment)
		r.Get("/", h.GetDeployments)
		r.Get("/{deploymentID}", h.GetDeploymentById)
	})
}
func (h *Handler) handleNewDeployment(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateNewDeploymentPayload
	ctx := r.Context()

	err := utils.ParseJson(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	newDeployment, err := h.store.CreateNewDeployment(ctx, payload.GitUrl)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = h.RedisStore.PushDeployment(ctx, newDeployment.ID.String())
	if err != nil {
		return
	}

	err = utils.WriteJson(w, http.StatusCreated, newDeployment)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) GetDeployments(w http.ResponseWriter, r *http.Request) {
	allDeployments, err := h.store.GetAllDeployments(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, allDeployments)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) GetDeploymentById(w http.ResponseWriter, r *http.Request) {

	uuid, err := uuid2.Parse(chi.URLParam(r, "deploymentID"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	deployment, err := h.store.GetDeploymentById(r.Context(), uuid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = utils.WriteJson(w, http.StatusOK, deployment)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

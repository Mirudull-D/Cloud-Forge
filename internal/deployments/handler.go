package deployments

import (
	"CloudHub/internal/docker"
	"CloudHub/internal/queue"
	"CloudHub/internal/types"
	"CloudHub/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	uuid2 "github.com/google/uuid"
)

type Handler struct {
	store       *Store
	RedisStore  *queue.RedisStore
	DockerStore *docker.Client
}

func NewHandler(store *Store, rdb *queue.RedisStore, dockerStore *docker.Client) *Handler {
	return &Handler{
		store:       store,
		RedisStore:  rdb,
		DockerStore: dockerStore,
	}
}
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/deployments", func(r chi.Router) {
		r.Post("/", h.handleNewDeployment)
		r.Get("/", h.GetDeployments)
		r.Get("/{deploymentID}", h.GetDeploymentById)
		r.Delete("/{deploymentID}", h.DeleteDeploymentById)
		r.Delete("/{deploymentID}/container", h.DeleteContainer)
		r.Post("/{deploymentID}/restart", h.RestartDeployment)
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

func (h *Handler) DeleteDeploymentById(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid2.Parse(chi.URLParam(r, "deploymentID"))
	ctx := r.Context()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	deployment, err := h.store.GetDeploymentById(ctx, uuid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if deployment.ContainerID.Valid {
		_, err := h.DockerStore.InspectContainer(ctx, deployment.ContainerID.String)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		err = h.DockerStore.RemoveContainer(ctx, deployment.ContainerID.String)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if deployment.ImageName.Valid {
			err = h.DockerStore.RemoveImage(
				ctx,
				deployment.ImageName.String,
			)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
		}

	}

	err = h.store.DeleteDeployment(ctx, uuid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusNoContent, nil)
	return
}
func (h *Handler) DeleteContainer(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid2.Parse(chi.URLParam(r, "deploymentID"))
	ctx := r.Context()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	deployment, err := h.store.GetDeploymentById(ctx, uuid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if deployment.ContainerID.Valid {
		_, err := h.DockerStore.InspectContainer(ctx, deployment.ContainerID.String)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		err = h.DockerStore.RemoveContainer(ctx, deployment.ContainerID.String)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}
	updatedDeployment, err := h.store.DeleteContainer(ctx, deployment.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusNoContent, updatedDeployment)
	return
}
func (h *Handler) RestartDeployment(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid2.Parse(chi.URLParam(r, "deploymentID"))
	ctx := r.Context()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	deployment, err := h.store.GetDeploymentById(ctx, uuid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if deployment.ContainerID.Valid {
		inspect, err := h.DockerStore.InspectContainer(ctx, deployment.ContainerID.String)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if !inspect.State.Running {
			utils.WriteJson(w, http.StatusBadRequest, "container is not running")
			return
		}
		err = h.DockerStore.RestartContainer(ctx, deployment.ContainerID.String)
		if err != nil {
			return
		}
		utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

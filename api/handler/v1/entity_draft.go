package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/e-space-uz/backend/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Security ApiKeyAuth
// @Router /v1/entity-draft [post]
// @Summary Create entity draft
// @Description API for creating entity draft
// @Tags entity-draft
// @Accept json
// @Produce json
// @Param entity body models.EntityDraftSwag true "entity"
// @Success 201 {object} models.CreateResponse

func (h *handlerV1) CreateEntityDraft(c *gin.Context) {
	var (
		entityDraft   models.CreateEntityDraft
		userInfo, err = h.UserInfo(c, true)
	)

	if err != nil {
		return
	}
	if err := c.ShouldBindJSON(&entityDraft); HandleHTTPError(c, http.StatusBadRequest, "EntityService.CreateEntityDraft.BindingJson", err) {
		return
	}

	fmt.Println(entityDraft.ID)

	if (entityDraft.Region.Soato) == 0 {
		HandleHTTPError(c, http.StatusConflict, "district soato required", errors.New("district soato required"))
		return
	}

	_, err = h.storage.Applicant().Get(context.Background(), userInfo.ID)
	if HandleHTTPError(c, http.StatusBadRequest, "EntityService.GetApplicant", err) {
		return
	}
	soato := strconv.Itoa(int(entityDraft.Region.Soato))
	entityDraft.EntityDraftSoato = soato

	resp, err := h.storage.EntityDraft().Create(
		context.Background(),
		&entityDraft,
	)

	if HandleHTTPError(c, http.StatusBadRequest, "EntityService.CreateEntityDraft", err) {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Router /v1/entity-draft/{entity_draft_id} [get]
// @Summary Get entity draft
// @Tags entity-draft
// @Accept json
// @Produce json
// @Param entity_draft_id path string true "entity_draft_id"
// @Success 200 {object} models.EntityDraft

func (h *handlerV1) GetEntityDraft(c *gin.Context) {
	ID := c.Param("entity_draft_id")

	_, err := primitive.ObjectIDFromHex(ID)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing entity id", err) {
		return
	}

	entity, err := h.storage.EntityDraft().Get(
		context.Background(),
		ID,
	)

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting entity", err) {
		return
	}

	c.JSON(http.StatusOK, entity)
}

package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
		entityDraft   models.CreateUpdateEntityDraft
		userInfo, err = h.UserInfo(c, true)
	)

	if err != nil {
		return
	}
	if err := c.ShouldBindJSON(&entityDraft); HandleHTTPError(c, http.StatusBadRequest, "EntityService.CreateEntityDraft.BindingJson", err) {
		return
	}

	entityDraft.Id = primitive.NewObjectID().Hex()
	fmt.Println(entityDraft.Id)
	objectID, err := primitive.ObjectIDFromHex(entityDraft.Id)
	if err != nil {
		fmt.Println("id")
	}
	fmt.Println(objectID)
	statusParent, err := h.storage.Status().GetParentStatus(context.Background(), &models.ASGetParentStatusRequest{
		ParentStatusId: models.ParentStatus,
	})
	if HandleHTTPError(c, http.StatusBadRequest, "EntityService.CreateEntityDraft.GetParentStatus", err) {
		return
	}
	entityDraft.StatusId = statusParent.Id

	if (entityDraft.Region.Soato) == 0 {
		HandleHTTPError(c, http.StatusConflict, "district soato required", errors.New("district soato required"))
		return
	}

	applicant, err := h.storage.Applicant().Get(context.Background(), &models.GetRequest{Id: userInfo.ID})
	if HandleHTTPError(c, http.StatusBadRequest, "EntityService.GetApplicant", err) {
		return
	}
	soato := strconv.Itoa(int(entityDraft.Region.Soato))
	entityDraft.EntityDraftSoato = soato
	entityDraft.Applicant = &models.ApplicantEntity{
		Name:        userInfo.FullName,
		UserId:      userInfo.ID,
		PhoneNumber: applicant.PhoneNumber,
	}

	resp, err := h.storage.EntityDraft().Create(
		context.Background(),
		&entityDraft,
	)

	if HandleHTTPError(c, http.StatusBadRequest, "EntityService.CreateEntityDraft", err) {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// @Router /v1/entity-draft/{entity_draft_id} [patch]
// @Summary Update entity draft comment
// @Description API for updating entity draft comment
// @Tags entity-draft
// @Accept json
// @Produce json
// @Param entity_draft_id path string true "entity_draft_id"
// @Param entity-draft body models.ConfirmEntityDraftSwag true "entity-draft"
// @Success 200 {object} models.EmptyResponse

func (h *handlerV1) ConfirmEntityDraft(c *gin.Context) {
	var (
		entityDraft   models.ConfirmEntityDraftSwag
		id            = c.Param("entity_draft_id")
		userInfo, err = h.UserInfo(c, true)
	)
	if err != nil {
		return
	}
	_, err = primitive.ObjectIDFromHex(id)
	if HandleHTTPError(c, http.StatusBadRequest, "error while entity_draft_id is incorrect format", err) {
		return
	}

	if err := c.BindJSON(&entityDraft); HandleHTTPError(c, http.StatusBadRequest, "error while binding json", err) {
		return
	}

	_, err = primitive.ObjectIDFromHex(entityDraft.StatusID)
	if HandleHTTPError(c, http.StatusBadRequest, "error while status_id is incorrect format", err) {
		return
	}
	resp, err := h.storage.EntityDraft().ConfirmEntityDraft(
		context.Background(),
		&models.ConfirmEntityDraftRequest{
			EntityDraftId: id,
			EntityId:      entityDraft.EntityID,
			StatusId:      entityDraft.StatusID,
			Comment:       entityDraft.Comment,
		})

	if HandleHTTPError(c, http.StatusBadRequest, "error while updating entity draft comment", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
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
		&models.ASGetRequest{
			Id: ID,
		})

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting entity", err) {
		return
	}

	c.JSON(http.StatusOK, entity)
}

// @Router /v1/entity-draft-expired [get]
// @Summary Getting expired draft entities
// @Description API for getting expired entity drafts
// @Tags entity-draft
// @Accept json
// @Produce json
// @Success 200 {object} models.GetAllEntityDraftsResponse

func (h *handlerV1) GetExpired(c *gin.Context) {
	var (
		response *models.GetAllEntityDraftsResponse
	)

	response, err := h.storage.EntityDraft().GetExpired(
		context.Background(),
		&models.GetExpiredDraftsRequest{
			Limit: 10,
		})

	if HandleHTTPError(c, http.StatusBadRequest, "Error while getting expired entities", err) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Router /v1/applicant-entity-draft/{user_id} [get]
// @Summary Getting All draft entities by user id
// @Description API for getting all entity drafts by user id
// @Tags entity-draft
// @Accept json
// @Produce json
// @Param user_id path string  true "user id"
// @Success 200 {object} models.GetAllEntityDraftsResponse

func (h *handlerV1) GetAllEntityDraftByUserID(c *gin.Context) {
	var (
		response *models.GetAllEntityDraftsResponse
		userID   = c.Param("user_id")
	)

	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	response, err = h.storage.EntityDraft().GetAll(
		context.Background(),
		&models.GetAllEntityDraftsRequest{
			UserId: userID,
			Page:   uint32(page),
			Limit:  uint32(limit),
		})

	if HandleHTTPError(c, http.StatusBadRequest, "Error while getting all entities", err) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Router /v1/entity-draft/{entity_draft_id} [put]
// @Summary Update entity draft
// @Description API for updating entity
// @Tags entity-draft
// @Accept json
// @Produce json
// @Param entity_draft_id path string  true "entity_draft_id"
// @Param entity-draft body models.EntityDraftSwag true "entity-draft"
// @Success 200 {object} models.CreateResponse

func (h *handlerV1) UpdateEntityDraft(c *gin.Context) {
	var (
		entity models.CreateUpdateEntityDraft
	)
	entityID := c.Param("entity_draft_id")

	_, err := primitive.ObjectIDFromHex(entityID)
	if HandleHTTPError(c, http.StatusBadRequest, "error while entity_draft_id is incorrect format", err) {
		return
	}
	err = c.ShouldBindJSON(&entity)
	if HandleHTTPError(c, http.StatusBadRequest, "error while updating entity", err) {
		return
	}

	if (entity.District.Code) == 0 {
		HandleHTTPError(c, http.StatusConflict, "district code required", errors.New("district code required"))
		return
	}

	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}

	_, err = h.storage.ActionHistory().Create(
		context.Background(),
		&models.ActionHistory{
			Id:            primitive.NewObjectID().Hex(),
			UserId:        userInfo.ID,
			Action:        models.EntityDraftUpdated,
			EntityId:      entityID,
			EntityName:    "entity-draft",
			UpdatedFields: []*models.UpdatedFields{},
		},
	)
	if HandleHTTPError(c, http.StatusBadRequest, "error while creating action history", err) {
		return
	}

	entity.EntityDraftSoato = strconv.FormatUint(entity.District.Soato, 10)
	entity.Id = entityID
	entity.Applicant = &models.ApplicantEntity{
		Name:   userInfo.FullName,
		UserId: userInfo.ID,
	}

	resp, err := h.storage.EntityDraft().Update(
		context.Background(),
		&entity)

	if HandleHTTPError(c, http.StatusBadRequest, "error while updating entity", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Router /v1/entity-draft/{entity_draft_id} [delete]
// @Summary Delete entity draft
// @Description API for deleting entity draft
// @Tags entity-draft
// @Accept json
// @Produce json
// @Param entity_draft_id path string  true "entity_draft_id"
// @Success 200 {object} models.SuccessResponse

func (h *handlerV1) DeleteEntityDraft(c *gin.Context) {
	entityID := c.Param("entity_draft_id")

	_, err := primitive.ObjectIDFromHex(entityID)
	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing uuid, entity id incorrect format", err) {
		return
	}

	entity, err := h.storage.EntityDraft().Delete(
		context.Background(),
		&models.ASDeleteRequest{Id: entityID})

	if HandleHTTPError(c, http.StatusBadRequest, "error while deleting entity draft", err) {
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (h *handlerV1) DeleteEntityDraftFromDB(c *gin.Context) {
	entityID := c.Param("entity_draft_id")

	_, err := primitive.ObjectIDFromHex(entityID)
	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing uuid, entity id incorrect format", err) {
		return
	}

	entity, err := h.storage.EntityDraft().DeleteFromDB(
		context.Background(),
		&models.ASDeleteRequest{Id: entityID})

	if HandleHTTPError(c, http.StatusBadRequest, "error while deleting entity draft from db", err) {
		return
	}

	c.JSON(http.StatusOK, entity)
}

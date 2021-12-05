package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/e-space-uz/backend/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Security ApiKeyAuth
// @Router /v1/entity [post]
// @Summary Create entity
// @Description API for creating entity
// @Tags entity
// @Accept json
// @Produce json
// @Param entity body models.CreateUpdateEntitySwag true "entity"
// @Success 201 {object} models.CreateResponse
func (h *handlerV1) CreateEntity(c *gin.Context) {
	var (
		entity     models.CreateUpdateEntity
		entitySwag models.CreateUpdateEntitySwag
	)

	if err := c.BindJSON(&entitySwag); HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.Create.BindingAction", err) {
		return
	}
	// TODO: MarshalUnmarshal make one function
	arrayOfByte, err := json.Marshal(entitySwag)

	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.Create.Marshalling", err) {
		return
	}
	if err = json.Unmarshal(arrayOfByte, &entity); HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.Create.Marshalling", err) {
		return
	}

	if (entity.District.Soato) == 0 {
		HandleHTTPError(c, http.StatusConflict, "Entity.Entity.Create", errors.New("district soato required"))
		return
	}

	resp, err := h.storage.Entity().Create(
		context.Background(),
		&entity)
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.Create", err) {
		return
	}
	// actionHistory := &ek_analytic_service.ActionHistory{
	// 	UserID:         userInfo.ID,
	// 	UserUniqueName: userInfo.Login,
	// 	Action:         models.EntityCreated,
	// 	EntityID:       id,
	// 	EntityName:     "entity",
	// }
	// h.CreateActionHistory(c, actionHistory)
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.CreateEntityDraft.ActionHistoryCreate", err) {
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// @Router /v1/entity/{entity_id} [get]
// @Summary Get entity
// @Tags entity
// @Accept json
// @Produce json
// @Param entity_id path string true "entity_id"
// @Success 200 {object} models.GetEntity
func (h *handlerV1) GetEntity(c *gin.Context) {
	var (
		ID     = c.Param("entity_id")
		_, err = primitive.ObjectIDFromHex(ID)
	)
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.ParseId", err) {
		return
	}

	entity, err := h.storage.Entity().Get(
		context.Background(),
		ID,
	)

	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.GetEntity", err) {
		return
	}

	c.JSON(http.StatusOK, entity)
}

// @Router /v1/entity-properties [get]
// @Summary Getting All entity `with` properties
// @Description API for getting all entity with properties
// @Tags entity
// @Accept json
// @Produce json
// @Param city_id query string  true "city_id"
// @Param region_id query string  true "region_id"
// @Param entity_type_code query string  false "entity_type_code"
// @Param entity_number query string  false "entity_number"
// @Param status_id query string  false "status_id"
// @Success 200 {object} models.GetAllEntitiesResponse
func (h *handlerV1) GetAllEntitiesWithProperties(c *gin.Context) {
	var (
		entityNumber = c.Query("entity_number")
		request      = &models.GetAllEntitiesRequest{
			EntityNumber: entityNumber,
		}
	)

	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}
	request.Page = uint32(page)
	request.Limit = uint32(limit)
	response, err := h.storage.Entity().GetAllWithProperties(
		context.Background(),
		request)
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.GetAllWithProperties", err) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Router /v1/entity-staff/{staff_id} [get]
// @Summary Getting All entities by staff id
// @Description API for getting all entities by staff id
// @Tags entity
// @Accept json
//@Produce json
//@Param staff_id path string  true "staff_id"
//@Success 200 {object} models.GetAllEntitiesResponse
func (h *handlerV1) GetAllByStaffID(c *gin.Context) {
	// var (
	// 	response = &models.GetAllEntitiesResponse{}
	// 	staffId  string
	// )

	// staffId = c.Param("staff_id")

	// page, err := ParseQueryParam(c, h.log, "page", "1")
	// if err != nil {
	// 	return
	// }

	// limit, err := ParseQueryParam(c, h.log, "limit", "20")
	// if err != nil {
	// 	return
	// }
	// entities, err := h.storage.Entity().GetAllByStaffID(
	// 	context.Background(),
	// 	&models.GetAllByStaffIDRequest{
	// 		StaffId: staffId,
	// 		Page:    uint32(page),
	// 		Limit:   uint32(limit),
	// 	})

	// if HandleHTTPError(c, http.StatusBadRequest, "Error while getting all entities by staff id", err) {
	// 	return
	// }
	// if err = ProtoToStructNumeric(&response, entities); HandleHTTPError(c, http.StatusInternalServerError, "error while parsing entities response", err) {
	// 	return
	// }
	// c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Router /v1/entity-status-update [put]
// @Summary Update entity status
// @Description API for updating entity status
// @Tags entity
// @Accept json
// @Produce json
// @Param entity body models.UpdateEntityStatus true "entity"
// @Success 200 {object} models.EmptyResponse
func (h *handlerV1) UpdateEntityStatus(c *gin.Context) {
	var (
		entity models.UpdateEntityStatus
	)

	if err := c.ShouldBindJSON(&entity); HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.UpdateEntityStatus", err) {
		return
	}

	c.JSON(http.StatusOK, "resp")
}

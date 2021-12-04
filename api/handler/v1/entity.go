package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
		entity     es.CreateUpdateEntity
		entitySwag models.CreateUpdateEntitySwag
	)

	if err := c.BindJSON(&entitySwag); HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.Create.BindingAction", err) {
		return
	}
	// TODO: MarshalUnmarshal make one function
	arrayOfByte, err := json.Marshal(entitySwag)

	if HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.Create.Marshalling", err) {
		return
	}
	if err = json.Unmarshal(arrayOfByte, &entity); HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.Create.Marshalling", err) {
		return
	}

	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}

	id := primitive.NewObjectID().Hex()

	if (entity.District.Soato) == 0 {
		HandleHTTPError(c, http.StatusBadRequest, http.StatusConflict, "Entity.Entity.Create", errors.New("district soato required"))
		return
	}

	entity.StaffId = userInfo.ID
	entity.Id = id

	if entity.EntityTypeCode == 2 {
		entity.StatusId = models.NewEntityStatusTypeTwo
	} else {
		entity.StatusId = models.NewEntityStatusTypeOne
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
	_, err = h.storage.ActionHistoryService().Create(
		context.Background(),
		&us.ActionHistory{
			Id:            primitive.NewObjectID().Hex(),
			UserId:        userInfo.ID,
			Action:        fmt.Sprintf("%s tomonidan %s", userInfo.Login, models.EntityCreated),
			EntityId:      entity.Id,
			EntityName:    "entity",
			UpdatedFields: []*us.UpdatedFields{},
		},
	)
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
// @Success 200 {object} entity_service.GetEntity

func (h *handlerV1) GetEntity(c *gin.Context) {
	var (
		ID     = c.Param("entity_id")
		_, err = primitive.ObjectIDFromHex(ID)
	)
	if HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.ParseId", err) {
		return
	}

	entity, err := h.storage.Entity().Get(
		context.Background(),
		&es.ASGetRequest{
			Id: ID,
		})

	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.GetEntity", err) {
		return
	}

	c.JSON(http.StatusOK, entity)
}

//// @Router /v1/entity-city/{city_id} [get]
//// @Summary Get entities by city id
//// @Tags entity
//// @Accept json
//// @Produce json
//// @Param city_id path string true "city_id"
//// @Success 200 {object} entity_service.GetAllEntitiesResponse
//// @Failure 400 {object} models.FailureResponse
//// @Failure 404 {object} models.FailureResponse
//// @Failure 500 {object} models.FailureResponse
//// @Failure 503 {object} models.FailureResponse
//func (h *handlerV1) GetEntitiesByCityID(c *gin.Context) {
//	//cityID := c.Param("city_id")
//	//
//	//_, err := primitive.ObjectIDFromHex(cityID)
//	//if HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "error while parsing entity id", err) {
//	//	return
//	//}
//	//
//	//entity, err := h.storage.Entity().GetByCityID(
//	//	context.Background(),
//	//	&es.GetEntitiesByCityIdRequest{
//	//		CityId: cityID,
//	//	})
//	//
//	//if HandleHTTPError(c, http.StatusBadRequest, "error while getting entities by city id", err) {
//	//	return
//	//}
//	//
//	//c.JSON(http.StatusOK, entity)
//}

// @Security ApiKeyAuth
// @Router /v1/entity [get]
// @Summary Getting All entities
// @Description API for getting all entities
// @Tags entity
// @Accept json
// @Produce json
// @Param find query models.GetAllEntitiesRequestSwag false "filters"
// @Success 200 {object} entity_service.GetAllEntitiesResponse

func (h *handlerV1) GetAllEntities(c *gin.Context) {
	var (
		cityID       = c.Query("city_id")
		regionID     = c.Query("region_id")
		statusID     = c.Query("status_id")
		fromDate     = c.Query("from_date")
		toDate       = c.Query("to_date")
		entityNumber = c.Query("entity_number")
		request      = &es.GetAllEntitiesRequest{
			CityId:   cityID,
			RegionId: regionID,
			StatusId: statusID,
			FromDate: fromDate,
			ToDate:   toDate,
		}
		userInfo, err = h.UserInfo(c, false)
	)

	if err != nil {
		return
	}
	// TODO: write a function for page and limit
	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	typeCode, err := ParseQueryParam(c, h.log, "entity_type_code", "0")
	if err != nil {
		return
	}

	if fromDate != "" {
		if err = ValidateTime(fromDate); HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.GetAllEntities.ParseDate", err) {
			return
		}
	}
	if toDate != "" {
		if err = ValidateTime(toDate); HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.GetAllEntities.ParseDate", err) {
			return
		}
	}

	if userInfo.UserType == "staff" {
		request.EntitySoato = userInfo.Soato
	}

	request.Page = uint32(page)
	request.Limit = uint32(limit)
	request.TypeCode = uint64(typeCode)
	request.EntityNumber = entityNumber
	//err = h.redisCache.Get(util.GenerateString(userInfo.Soato, cityID, regionID, statusID, date, strconv.Itoa(typeCode), strconv.Itoa(entityNumber)), &response)
	//if err == models.ErrCacheMiss {
	response, err := h.storage.Entity().GetAll(
		context.Background(),
		request)
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.GetAllEntities", err) {
		return
	}
	//fmt.Println(response)
	//	if err = h.redisCache.Set(util.GenerateString(userInfo.Soato, cityID, regionID, statusID, date, strconv.Itoa(typeCode), strconv.Itoa(entityNumber)), response); HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.CacheGetAllEntity", err) {
	//		return
	//	}
	//} else if HandleHTTPError(c, http.StatusBadRequest, http.StatusInternalServerError, "Entity.Entity.GetAll.SetCaching", err) {
	//	return
	//}
	c.JSON(http.StatusOK, response)
}

// @Router /v1/entity-collection [get]
// @Summary Getting All entity collections
// @Description API for getting all entity collections
// @Tags entity
// @Accept json
// @Produce json
// @Param collection_name query string  true "collection_name"
// @Param search query string  false "search"
// @Success 200 {object} entity_service.GetCollectionResponse

func (h *handlerV1) GetCollections(c *gin.Context) {
	var (
		collectionName = c.Query("collection_name")
		search         = c.DefaultQuery("search", "")
		request        = &es.GetCollectionRequest{
			CollectionName: collectionName,
			Search:         search,
		}
	)
	response, err := h.storage.Entity().GetCollection(context.Background(), request)
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.GetAllEntitieCollections", err) {
		return
	}

	c.JSON(http.StatusOK, response)
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
// @Success 200 {object} entity_service.GetAllEntitiesResponse

func (h *handlerV1) GetAllEntitiesWithProperties(c *gin.Context) {
	var (
		entityTypeCode int
		cityID         = c.Query("city_id")
		regionID       = c.Query("region_id")
		statusID       = c.Query("status_id")
		entityNumber   = c.Query("entity_number")
		request        = &es.GetAllEntitiesRequest{
			CityId:       cityID,
			RegionId:     regionID,
			StatusId:     statusID,
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
	request.TypeCode = uint64(entityTypeCode)
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
//@Success 200 {object} entity_service.GetAllEntitiesResponse

func (h *handlerV1) GetAllByStaffID(c *gin.Context) {
	// var (
	// 	response = &es.GetAllEntitiesResponse{}
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
	// 	&es.GetAllByStaffIDRequest{
	// 		StaffId: staffId,
	// 		Page:    uint32(page),
	// 		Limit:   uint32(limit),
	// 	})

	// if HandleHTTPError(c, http.StatusBadRequest, "Error while getting all entities by staff id", err) {
	// 	return
	// }
	// if err = ProtoToStructNumeric(&response, entities); HandleHTTPError(c, http.StatusBadRequest, http.StatusInternalServerError, "error while parsing entities response", err) {
	// 	return
	// }
	// c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Router /v1/entity/{entity_id} [put]
// @Summary Update entity
// @Description API for updating entity
// @Tags entity
// @Accept json
// @Produce json
// @Param entity_id path string  true "entity_id"
// @Param entity body models.CreateUpdateEntitySwag true "entity"
// @Success 200 {object} models.CreateResponse

func (h *handlerV1) UpdateEntity(c *gin.Context) {
	var (
		entity   es.CreateUpdateEntity
		entityID = c.Param("entity_id")
	)

	_, err := primitive.ObjectIDFromHex(entityID)
	if HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.UpdateEntity", err) {
		return
	}

	if err = c.ShouldBindJSON(&entity); HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.UpdateEntity", err) {
		return
	}

	if (entity.District.Soato) == 0 {
		HandleHTTPError(c, http.StatusBadRequest, http.StatusConflict, "Entity.Entity.UpdateEntity", errors.New("district soato required"))
		return
	}
	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}

	entityId := primitive.NewObjectID().Hex()
	_, err = h.storage.ActionHistoryService().Create(
		context.Background(),
		&us.ActionHistory{
			Id:            primitive.NewObjectID().Hex(),
			UserId:        userInfo.ID,
			Action:        models.EntityUpdated,
			EntityId:      entityId,
			EntityName:    "entity",
			UpdatedFields: []*us.UpdatedFields{},
		},
	)
	if HandleHTTPError(c, http.StatusBadRequest, "UserService.ActionHistory.Create", err) {
		return
	}

	soato := strconv.Itoa(int(entity.District.Soato))
	entity.EntitySoato = soato
	entity.Id = entityID
	entity.StaffId = userInfo.ID

	resp, err := h.storage.Entity().Update(
		context.Background(),
		&entity)
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.UpdateEntity", err) {
		return
	}

	// getResponse, err := h.storage.Entity().GetAll(
	// 	context.Background(),
	// 	&es.GetAllEntitiesRequest{
	// 		EntitySoato: userInfo.Soato,
	// 		Page:        1,
	// 		Limit:       20,
	// 	})
	// if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.GetAll", err) {
	// 	return
	// }

	// if err = h.redisCache.Set(util.GenerateString(userInfo.Soato, "00"), getResponse); HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.Create.SetCache", err) {
	// 	return
	// }

	c.JSON(http.StatusOK, resp)
}

// @Security ApiKeyAuth
// @Router /v1/entity-parent-status [put]
// @Summary Update entity with parent status
// @Description API for updating entity with parent statuss
// @Tags entity
// @Accept json
// @Produce json
// @Param entity body models.UpdateEntityStatus true "entity"
// @Success 200 {object} models.EmptyResponse

func (h *handlerV1) UpdateEntityParentStatus(c *gin.Context) {
	var (
		entity          es.UpdateEntityStatusRequest
		actionHistoryID = primitive.NewObjectID().Hex()
	)

	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}

	if err := c.ShouldBindJSON(&entity); HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.UpdateEntityParentStatus", err) {
		return
	}
	_, err = primitive.ObjectIDFromHex(entity.EntityId)
	if HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.UpdateEntityParentStatus", err) {
		return
	}
	_, err = primitive.ObjectIDFromHex(entity.StatusId)
	if HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.UpdateEntityParentStatus", err) {
		return
	}

	nextStatus, err := h.storage.StatusService().GetParentStatus(
		context.Background(),
		&es.ASGetParentStatusRequest{
			ParentStatusId: entity.StatusId,
		})
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.UpdateEntityParentStatus", err) {
		return
	}
	entity.StatusId = nextStatus.Id
	resp, err := h.storage.Entity().UpdateEntityStatus(
		context.Background(),
		&entity)

	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.UpdateEntityParentStatus", err) {
		return
	}
	_, err = h.storage.ActionHistoryService().Create(
		context.Background(),
		&us.ActionHistory{
			Id:            actionHistoryID,
			UserId:        userInfo.ID,
			Action:        "Ariza keyingi qadamga o'tkazildi", // Todo: make global variables
			EntityId:      entity.EntityId,
			EntityName:    "entity",
			UpdatedFields: []*us.UpdatedFields{},
		})
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.UpdateEntityStatusByActionID", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
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
		entity          es.UpdateEntityStatusRequest
		actionHistoryID = primitive.NewObjectID().Hex()
	)

	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}

	if err := c.ShouldBindJSON(&entity); HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.UpdateEntityStatus", err) {
		return
	}
	_, err = primitive.ObjectIDFromHex(entity.EntityId)
	if HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.UpdateEntityStatus", err) {
		return
	}
	_, err = primitive.ObjectIDFromHex(entity.StatusId)
	if HandleHTTPError(c, http.StatusBadRequest, http.StatusBadRequest, "Entity.Entity.UpdateEntityStatus", err) {
		return
	}

	resp, err := h.storage.Entity().UpdateEntityStatus(
		context.Background(),
		&entity)

	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.UpdateEntityStatus", err) {
		return
	}
	_, err = h.storage.ActionHistoryService().Create(
		context.Background(),
		&us.ActionHistory{
			Id:            actionHistoryID,
			UserId:        userInfo.ID,
			Action:        "Ariza statusi o'zgardi",
			EntityId:      entity.EntityId,
			EntityName:    "entity",
			UpdatedFields: []*us.UpdatedFields{},
		})
	if HandleHTTPError(c, http.StatusBadRequest, "Entity.Entity.UpdateEntityStatus", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

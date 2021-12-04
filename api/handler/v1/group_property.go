package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/group-property [post]
// @Summary Create group_property
// @Description API for creating group_property
// @Tags group_property
// @Accept json
// @Produce json
// @Param group_property body models.GroupPropertySwag true "group property"
// @Success 201 {object} ek_variables.CreateResponse

func (h *handlerV1) CreateGroupProperty(c *gin.Context) {
	var (
		groupProperty models.CreateGroupProperty
	)
	err := c.ShouldBindJSON(&groupProperty)
	if HandleHTTPError(c, http.StatusBadRequest, "DiscussionLogicService.Action.Create.BindingAction", err) {
		return
	}
	groupProperty.Id = primitive.NewObjectID().Hex()

	resp, err := h.storage.GroupProperty().Create(
		context.Background(),
		&groupProperty,
	)

	if HandleHTTPError(c, http.StatusBadRequest, "error while creating group property", err) {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Router /v1/group-property/{group_property_id} [get]
// @Summary Get group property
// @Tags group_property
// @Accept json
// @Produce json
// @Param group_property_id path string true "group_property_id"
// @Success 200 {object} models.GroupProperty

func (h *handlerV1) GetGroupProperty(c *gin.Context) {
	var (
		response        models.GroupProperty
		groupPropertyID = c.Param("group_property_id")
	)
	_, err := primitive.ObjectIDFromHex(groupPropertyID)
	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing objectID, incorrect format", err) {
		return
	}
	property, err := h.storage.GroupProperty().Get(
		context.Background(),
		&models.ASGetRequest{
			Id: groupPropertyID,
		})

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting group property", err) {
		return
	}
	err = ProtoToStruct(&response, property)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing proto to struct", err) {
		return
	}
	c.JSON(http.StatusOK, response)
}

// @Router /v1/group-property-status/{status_id} [get]
// @Summary Get group property
// @Tags group_property
// @Accept json
// @Produce json
// @Param status_id path string true "status_id"
// @Success 200 {object} models.GetAllGroupPropertyByStatusIdResponse

func (h *handlerV1) GetGroupPropertyByStatusID(c *gin.Context) {
	var (
		statusID = c.Param("status_id")
		response *models.GetAllGroupPropertyByStatusIdResponse
		redisKey = ek_variables.RedisGroupPropertyKey + statusID
	)
	_, err := primitive.ObjectIDFromHex(statusID)
	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing objectID, incorrect format", err) {
		return
	}

	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}
	// userInfo := models.LoginInfo{
	// 	ID:             "6142c2046074f7fa21292a28",
	// 	RoleID:         "6142c2046074f7fa21292a28",
	// 	OrganizationID: "61027e35772d476a220673ec",
	// 	Login:          "6142c2046074f7fa21292a28",
	// }
	// err = h.redisCache.Get(redisKey, &response)
	// if err == ek_variables.ErrCacheMiss {
	response, err = h.storage.GroupProperty().GetByStatusId(
		context.Background(),
		&models.ASGetByStatusIdRequest{
			StatusId:       statusID,
			OrganizationId: userInfo.OrganizationID,
		})
	if HandleHTTPError(c, http.StatusBadRequest, "error while getting group property by status id", err) {
		return
	}
	fmt.Println(len(response.GroupProperties))
	if err = h.redisCache.SetWithDeadline(redisKey, response, time.Minute*4); HandleHTTPError(c, http.StatusBadRequest, "EntityService.Entity.CacheGetAllEntity", err) {
		return
	}
	// } else if HandleHTTPError(c, http.StatusInternalServerError, "EntityService.GroupProperty.GetAll.GetCaching", err) {
	// 	return
	// }

	c.JSON(http.StatusOK, response)
}

// @Router /v1/group-property [get]
// @Summary Getting All Group Properties
// @Description API for getting all group properties
// @Tags group_property
// @Accept json
// @Produce json
// @Param search query string false "search"
// @Param page query integer false "page"
// @Param limit query integer false "limit"
// @Success 200 {object} models.GetAllGroupPropertiesResponse

func (h *handlerV1) GetAllGroupProperties(c *gin.Context) {
	var (
		search   = c.Query("search")
		response models.GetAllGroupPropertiesResponse
	)

	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	groupProperties, err := h.storage.GroupProperty().GetAll(
		context.Background(),
		&models.GetAllGroupPropertiesRequest{
			Search: search,
			Page:   uint32(page),
			Limit:  uint32(limit),
		})

	if HandleHTTPError(c, http.StatusBadRequest, "Erro while getting all group properties", err) {
		return
	}
	err = ProtoToStruct(&response, groupProperties)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing proto to struct", err) {
		return
	}
	fmt.Println(response)
	c.JSON(http.StatusOK, response)
}

// @Router /v1/group-property/{group_property_id} [put]
// @Summary Update Group Property
// @Description API for updating group property
// @Tags group_property
// @Accept json
// @Produce json
// @Param group_property_id path string  true "group_property_id"
// @Param property body models.GroupPropertySwag true "group property"
// @Success 200 {object} ek_variables.CreateResponse

func (h *handlerV1) UpdateGroupProperty(c *gin.Context) {
	var (
		groupProperty   models.CreateGroupProperty
		groupPropertyID = c.Param("group_property_id")
		_, err          = primitive.ObjectIDFromHex(groupPropertyID)
	)
	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing objectID, incorrect format", err) {
		return
	}

	if err = c.ShouldBindJSON(&groupProperty); HandleHTTPError(c, http.StatusBadRequest, "error while binding model to json", err) {
		return
	}

	groupProperty.Id = groupPropertyID
	resp, err := h.storage.GroupProperty().Update(
		context.Background(),
		&groupProperty)

	if HandleHTTPError(c, http.StatusBadRequest, "error while updating group property", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Router /v1/group-property-type [get]
// @Summary Getting All Group Properties By Type
// @Description API for getting all group properties by type
// @Tags group_property
// @Accept json
// @Produce json
// @Param step query integer true "step"
// @Param type query integer true "type"
// @Success 200 {object} models.GetAllGroupPropertiesByTypeResponse

func (h *handlerV1) GetAllGroupPropertiesByType(c *gin.Context) {
	var (
		stepQuery   = c.Query("step")
		typeOfQuery = c.Query("type")
		response    models.GetAllGroupPropertiesByTypeResponse
	)
	step, err := strconv.Atoi(stepQuery)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing query to request", err) {
		return
	}
	typeOf, err := strconv.Atoi(typeOfQuery)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing query to request", err) {
		return
	}

	groupProperties, err := h.storage.GroupProperty().GetAllByType(
		context.Background(),
		&models.GetAllGroupPropertiesByTypeRequest{
			Step: uint32(step),
			Type: uint32(typeOf),
		})

	if HandleHTTPError(c, http.StatusBadRequest, "Erro while getting all group properties by type", err) {
		return
	}
	if err = ProtoToStructNumeric(&response, groupProperties); HandleHTTPError(c, http.StatusInternalServerError, "error while parsing proto to struct", err) {
		return
	}
	c.JSON(http.StatusOK, response)
}

// // @Router /v1/group-property/{group_property_id} [delete]
// // @Summary Delete Property
// // @Description API for deleting property
// // @Tags property
// // @Accept json
// // @Produce json
// // @Param group_property_id path string  true "group_property_id"
// // @Success 200 {object} ek_variables.SuccessResponse
// // @Failure 400 {object} ek_variables.FailureResponse
// // @Failure 404 {object} ek_variables.FailureResponse
// // @Failure 500 {object} ek_variables.FailureResponse
// // @Failure 503 {object} ek_variables.FailureResponse
// func (h *handlerV1) DeleteGroupProperty(c *gin.Context) {
// 	groupPropertyID := c.Param("group_property_id")

// 	_, err := primitive.ObjectIDFromHex(groupPropertyID)
// 	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing objectID, property id incorrect format", err) {
// 		return
// 	}

// 	resp, err := h.storage.GroupProperty().Delete(
// 		context.Background(),
// 		&models.ASDeleteRequest{Id: groupPropertyID})

// 	if HandleHTTPError(c,"error while deleting group property", err) {
// 		return
// 	}

// 	c.JSON(http.StatusOK, resp)
// }

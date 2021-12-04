package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/property [post]
// @Summary Create property
// @Description API for creating property
// @Tags property
// @Accept json
// @Produce json
// @Param property body models.PropertySwag true "property"
// @Success 201 {object} ek_variables.CreateResponse
func (h *handlerV1) CreateProperty(c *gin.Context) {
	var (
		property models.CreateUpdateProperty
		err      = c.ShouldBindJSON(&property)
	)
	if HandleHTTPError(c, http.StatusBadRequest, "DiscussionLogicService.Action.Create.BindingAction", err) {
		return
	}
	property.Id = primitive.NewObjectID().Hex()

	resp, err := h.storage.Property().Create(
		context.Background(),
		&property,
	)

	if HandleHTTPError(c, http.StatusBadGateway, "error while creating Property", err) {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Router /v1/property/{property_id} [get]
// @Summary Get Property
// @Tags property
// @Accept json
// @Produce json
// @Param property_id path string true "property_id"
// @Success 200 {object} entity_service.Property

func (h *handlerV1) GetProperty(c *gin.Context) {
	var (
		propertyID = c.Param("property_id")
		response   *models.Property
	)
	_, err := primitive.ObjectIDFromHex(propertyID)
	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing objectID, incorrect format", err) {
		return
	}
	property, err := h.storage.Property().Get(
		context.Background(),
		propertyID,
	)
	if HandleHTTPError(c, http.StatusBadGateway, "error while getting property", err) {
		return
	}

	err = ProtoToStruct(&response, property)
	if HandleHTTPError(c, http.StatusInternalServerError, "error while parsing property response", err) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Router /v1/property [get]
// @Summary Getting All Properties
// @Description API for getting all properties
// @Tags property
// @Accept json
// @Produce json
// @Param name query string false "name"
// @Param page query integer false "page"
// @Param limit query integer false "limit"
// @Success 200 {object} entity_service.GetAllPropertiesResponse

func (h *handlerV1) GetAllProperties(c *gin.Context) {
	var (
		response *models.GetAllPropertiesResponse
		name     = c.Query("name")
	)

	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	properties, err := h.storage.Property().GetAll(
		context.Background(),
		&models.GetAllPropertiesRequest{
			Name:  name,
			Page:  uint32(page),
			Limit: uint32(limit),
		})

	if HandleHTTPError(c, http.StatusBadGateway, "Erro while getting all properties", err) {
		return
	}

	if err = ProtoToStruct(&response, properties); HandleHTTPError(c, http.StatusInternalServerError, "error while parsing properties response", err) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Router /v1/property/{property_id} [put]
// @Summary Update Property
// @Description API for updating property
// @Tags property
// @Accept json
// @Produce json
// @Param property_id path string  true "property_id"
// @Param property body models.PropertySwag true "property"
// @Success 200 {object} ek_variables.CreateResponse

func (h *handlerV1) UpdateProperty(c *gin.Context) {
	var (
		property models.CreateUpdateProperty
	)
	propertyID := c.Param("property_id")

	_, err := primitive.ObjectIDFromHex(propertyID)
	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing objectID, incorrect format", err) {
		return
	}

	err = c.ShouldBindJSON(&property)
	if HandleHTTPError(c, http.StatusBadRequest, "error while binding model to json", err) {
		return
	}

	property.Id = propertyID
	resp, err := h.storage.Property().Update(
		context.Background(),
		&property)

	if HandleHTTPError(c, http.StatusBadGateway, "error while updating property", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Router /v1/property/{property_id} [delete]
// @Summary Delete Property
// @Description API for deleting property
// @Tags property
// @Accept json
// @Produce json
// @Param property_id path string  true "property_id"
// @Success 200 {object} ek_variables.SuccessResponse

func (h *handlerV1) DeleteProperty(c *gin.Context) {
	propertyID := c.Param("property_id")

	_, err := primitive.ObjectIDFromHex(propertyID)
	if HandleHTTPError(c, http.StatusBadRequest, "Error while parsing objectID, property id incorrect format", err) {
		return
	}

	resp, err := h.storage.Property().Delete(
		context.Background(),
		&models.ASDeleteRequest{Id: propertyID})

	if HandleHTTPError(c, http.StatusBadGateway, "error while deleting property", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

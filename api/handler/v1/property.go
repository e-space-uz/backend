package v1

import (
	"context"
	"net/http"

	"github.com/e-space-uz/backend/models"
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
	property.ID = primitive.NewObjectID()

	resp, err := h.storage.Property().Create(
		context.Background(),
		&property,
	)

	if HandleHTTPError(c, http.StatusBadGateway, "error while creating Property", err) {
		return
	}

	c.JSON(http.StatusCreated, resp)
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

	err = h.storage.Property().Update(
		context.Background(),
		&property)

	if HandleHTTPError(c, http.StatusBadGateway, "error while updating property", err) {
		return
	}

	c.JSON(http.StatusOK, "resp")
}

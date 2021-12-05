package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	groupProperties, _, err := h.storage.GroupProperty().GetAll(
		context.Background(),
		uint32(page),
		uint32(limit),
	)

	if HandleHTTPError(c, http.StatusBadRequest, "Erro while getting all group properties", err) {
		return
	}
	c.JSON(http.StatusOK, groupProperties)
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
	)
	step, err := strconv.Atoi(stepQuery)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing query to request", err) {
		return
	}
	typeOf, err := strconv.Atoi(typeOfQuery)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing query to request", err) {
		return
	}

	groupProperties, _, err := h.storage.GroupProperty().GetAllByType(
		context.Background(),
		uint32(step),
		uint32(typeOf),
	)

	if HandleHTTPError(c, http.StatusBadRequest, "Erro while getting all group properties by type", err) {
		return
	}
	c.JSON(http.StatusOK, groupProperties)
}

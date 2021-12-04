package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/region/{region_id} [get]
// @Summary Get region
// @Description API for getting region
// @Tags region
// @Accept json
// @Produce json
// @Param region_id path string  true "region_id"

func (h *handlerV1) GetRegion(c *gin.Context) {
	regionID := c.Param("region_id")

	_, err := primitive.ObjectIDFromHex(regionID)

	if HandleHTTPError(c, http.StatusBadRequest, "error while creating regionID", err) {
		return
	}
	region, err := h.storage.Region().Get(context.Background(), regionID)

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting region ", err) {
		return
	}

	c.JSON(http.StatusOK, region)
}

// @Security ApiKeyAuth
// @Router /v1/region [get]
// @Summary Getting All regions
// @Description API for getting all regiones
// @Tags region
// @Accept json
// @Produce json
// @Param name query string  false "name"
// @Param soato query string  false "soato"
// @Param page query integer false "page"
// @Param limit query integer false "limit"
// @Success 200 {object} models.GetAllRegionsResponse

func (h *handlerV1) GetAllRegions(c *gin.Context) {
	var (
		soatoQuery = c.Query("soato")
		soato      int
		_, err     = h.UserInfo(c, false)
	)
	if HandleHTTPError(c, http.StatusUnauthorized, "SettingService.Region.GetAllRegion", err) {
		return
	}
	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}
	if soatoQuery != "" {
		soato, err = ParseQueryParam(c, h.log, "soato", "0")
		if err != nil {
			return
		}
	}
	fmt.Println(soato)

	regions, count, err := h.storage.Region().GetAll(
		context.Background(),
		uint32(page),
		uint32(limit),
	)

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting all regions", err) {
		return
	}

	c.JSON(http.StatusOK, regions)
}

// @Router /v1/regions/{city_id} [get]
// @Summary Getting All regions by city ID
// @Description API for getting all regiones by city ID
// @Tags region
// @Accept json
// @Produce json
// @Param city_id path string true "city_id"
// @Param name query string false "name"
// @Success 200 {object} models.GetAllRegionsResponse

func (h *handlerV1) GetAllRegionsByCityID(c *gin.Context) {
	var (
		cityID = c.Param("city_id")
		_, err = primitive.ObjectIDFromHex(cityID)
		name   = c.Query("name")
		// redisKey = ek_variables.RedisRegionKey + c.Request.URL.Query().Encode()
	)
	if HandleHTTPError(c, http.StatusBadRequest, "SettingService.GetAllRegionsByCityID.ParsingCityID", err) {
		return
	}
	response, count, err := h.storage.Region().GetAllByCity(
		context.Background(),
		cityID,
		name,
	)
	if HandleHTTPError(c, http.StatusBadRequest, "error while getting all regions", err) {
		return
	}
	c.JSON(http.StatusOK, response)
}

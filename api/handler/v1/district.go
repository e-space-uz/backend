package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/district/{district_id} [get]
// @Summary Get district
// @Description API for getting district
// @Tags district
// @Accept json
// @Produce json
// @Param district_id path string  true "district_id"

func (h *handlerV1) GetDistrict(c *gin.Context) {
	var (
		districtID = c.Param("district_id")
	)
	_, err := primitive.ObjectIDFromHex(districtID)

	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing district uuid ", err) {
		return
	}
	district, err := h.storage.District().Get(context.Background(), districtID)

	if HandleHTTPError(c, http.StatusInternalServerError, "error while getting district ", err) {
		return
	}

	c.JSON(http.StatusOK, district)
}

// @Router /v1/district [get]
// @Summary Getting All districts
// @Description API for getting all districtes
// @Tags district
// @Accept json
// @Produce json
// @Param name query string  false "name"
// @Param soato query string  false "soato"
// @Param page query integer false "page"
// @Param limit query integer false "limit"
// @Success 200 {object} models.GetAllDistrictsResponse

func (h *handlerV1) GetAllDistricts(c *gin.Context) {

	var (
		name       = c.Query("name")
		soatoQuery = c.Query("soato")
		soato      int
	)
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

	districts, err := h.storage.District().GetAll(
		context.Background(),
		&models.GetAllDistrictsRequest{
			Name:  name,
			Soato: uint32(soato),
			Page:  uint32(page),
			Limit: uint32(limit),
		})

	if HandleHTTPError(c, "error while getting all districts", err) {
		return
	}

	c.JSON(http.StatusOK, districts)
}

// @Router /v1/districts/{city_id}/{region_id} [get]
// @Summary Getting All districts by region ID
// @Description API for getting all districts by region ID
// @Tags district
// @Accept json
// @Produce json
// @Param city_id path string true "city_id"
// @Param region_id path string true "region_id"
// @Param name query string false "name"
// @Success 200 {object} models.GetAllDistrictsResponse

func (h *handlerV1) GetAllDistrictsByRegionID(c *gin.Context) {
	var (
		regionID = c.Param("region_id")
		cityID   = c.Param("city_id")
		name     = c.Query("name")
		_, err   = primitive.ObjectIDFromHex(cityID)
		response = &models.GetAllDistrictsResponse{}
		// redisKwwey = ek_variables.RedisDistrictKey + c.Request.URL.Query().Encode()
	)
	if HandleHTTPError(c, http.StatusBadRequest, "SettingService.GetAllDistrictsByRegionID.ParsingCityObjectID ", err) {
		return
	}
	_, err = primitive.ObjectIDFromHex(regionID)

	if HandleHTTPError(c, http.StatusBadRequest, "SettingService.GetAllDistrictsByRegionID.ParsingRegionObjectID ", err) {
		return
	}
	// err = h.redisCache.Get(redisKey, &response)
	// if err == ek_variables.ErrCacheMiss {
	response, err = h.storage.District().GetAllByCityRegion(
		context.Background(),
		regionID,
		cityID,
		name,
	)

	if HandleHTTPError(c, "SettingService.GetAllDistrictsByRegionID.ServiceInternal", err) {
		return
	}

	c.JSON(http.StatusOK, response)
}

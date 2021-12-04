package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/city/{city_id} [get]
// @Summary Get City
// @Description API for getting city
// @Tags city
// @Accept json
// @Produce json
// @Param city_id path string  true "city_id"
// @Success 200 {object} models.City

func (h *handlerV1) GetCity(c *gin.Context) {
	var (
		cityID = c.Param("city_id")
		_, err = primitive.ObjectIDFromHex(cityID)
	)
	if HandleHTTPError(c, http.StatusBadRequest, "SettingService.GetCity.ParseCityID", err) {
		return
	}
	city, err := h.storage.City().Get(context.Background(), &models.GetReq{
		Id: cityID,
	})

	if HandleHTTPError(c, http.StatusBadRequest, "SettingService.GetCity.InterService", err) {
		return
	}

	c.JSON(http.StatusOK, city)
}

// @Router /v1/city [get]
// @Summary Getting All cities
// @Description API for getting all cityes
// @Tags city
// @Accept json
// @Produce json
// @Param name query string  false "name"
// @Param soato query string  false "soato"
// @Param page query integer false "page"
// @Param limit query integer false "limit"
// @Success 200 {object} models.GetAllCitiesResponse
func (h *handlerV1) GetAllCities(c *gin.Context) {
	var (
		name       = c.Query("name")
		soatoQuery = c.Query("soato")
	)
	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}
	cities, count, err := h.storage.City().GetAll(
		context.Background(),
		uint32(page),
		uint32(limit),
	)

	if HandleHTTPError(c, http.StatusBadRequest, "SettingService.GetAllCities.InternalService", err) {
		return
	}

	c.JSON(http.StatusOK, cities)
}

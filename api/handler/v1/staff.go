package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/staff/{staff_id} [get]
// @Summary Get Staff
// @Description API for getting staff
// @Tags staff
// @Accept json
// @Produce json
// @Param staff_id path string  true "staff_id"
// @Success 200 {object} user_service.Staff
func (h *handlerV1) GetStaff(c *gin.Context) {
	id := c.Param("staff_id")

	_, err := primitive.ObjectIDFromHex(id)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing id", err) {
		return
	}

	staff, err := h.storage.Staff().Get(context.Background(), id)

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting staff ", err) {
		return
	}

	c.JSON(http.StatusOK, staff)
}

// @Security ApiKeyAuth
// @Router /v1/staff-by-token [get]
// @Summary Get Staff by token
// @Description API for getting staff
// @Tags staff
// @Accept json
// @Produce json
// @Success 200 {object} user_service.Staff
func (h *handlerV1) GetStaffByToken(c *gin.Context) {
	userInfo, err := h.UserInfo(c, true)
	if HandleHTTPError(c, http.StatusUnauthorized, "UserService.StaffByToken", err) {
		return
	}
	staff, err := h.storage.Staff().Get(context.Background(), userInfo.ID)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.StaffByToken.GetStaff", err) {
		return
	}

	c.JSON(http.StatusOK, staff)
}

// @Router /v1/staff [get]
// @Summary Getting All Staffs
// @Description API for getting all staff
// @Tags staff
// @Accept json
// @Produce json
// @Param find query ek_user_service.GetAllStaffsRequestSwag false "filters"
// @Success 200 {object} user_service.GetAllStaffsResponse
func (h *handlerV1) GetAllStaffs(c *gin.Context) {

	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	_, err = ParseBoolParam(c, h.log, "status", "false")
	if err != nil {
		return
	}

	staffs, _, err := h.storage.Staff().GetAll(
		context.Background(),
		uint32(page),
		uint32(limit),
	)

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting all staffs", err) {
		return
	}

	c.JSON(http.StatusOK, staffs)
}

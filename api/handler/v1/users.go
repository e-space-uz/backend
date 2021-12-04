package v1

import (
	"context"
	"net/http"

	"github.com/e-space-uz/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/staff/{user_id} [get]
// @Summary Get Staff
// @Description API for getting staff
// @Tags staff
// @Accept json
// @Produce json
// @Param user_id path string  true "user_id"
// @Success 200 {object} user_service.Staff

func (h *handlerV1) GetUser(c *gin.Context) {
	id := c.Param("user_id")

	_, err := primitive.ObjectIDFromHex(id)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing id", err) {
		return
	}

	staff, err := h.storage.User().Get(context.Background(), id)

	if HandleHTTPError(c, http.StatusBadGateway, "error while getting staff ", err) {
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
	staff, err := h.storage.User().Get(context.Background(), &models.GetRequest{})

	if HandleHTTPError(c, http.StatusBadGateway, "UserService.StaffByToken.GetStaff", err) {
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
	var (
		response       *models.GetAllStaffsResponse
		phoneNumber    = c.Query("phone_number")
		organizationId = c.Query("organization_id")
		searchString   = c.Query("search_string")
		roleId         = c.Query("role_id")
		soato          = c.Query("soato")
	)

	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	status, err := ParseBoolParam(c, h.log, "status", "false")
	if err != nil {
		return
	}

	staffs, err := h.storage.User().GetAll(
		context.Background(),
		&models.GetAllStaffsRequest{
			Page:           uint32(page),
			Limit:          uint32(limit),
			Soato:          soato,
			PhoneNumber:    phoneNumber,
			OrganizationId: organizationId,
			SearchString:   searchString,
			RoleId:         roleId,
			Status:         status)

	if HandleHTTPError(c, http.StatusBadGateway, "error while getting all staffs", err) {
		return
	}

	c.JSON(http.StatusOK, staffs)
}

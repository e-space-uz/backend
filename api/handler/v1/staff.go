package v1

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/pkg/security"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/staff [post]
// @Summary staff
// @Description API for creating staff
// @Tags staff
// @Accept json
// @Produce json
// @Param staff body ek_user_service.CreateUpdateStaffSwag  true "Staff"
// @Success 201 {object} user_service.CreateResponse
func (h *handlerV1) CreateStaff(c *gin.Context) {
	var (
		staff *models.CreateUpdateStaff
	)

	if err := c.ShouldBindJSON(&staff); HandleHTTPError(c, http.StatusBadRequest, "DiscussionLogicService.Action.Create.BindingAction", err) {
		return
	}
	_, err := primitive.ObjectIDFromHex(staff.RoleId)
	if HandleHTTPError(c, http.StatusBadRequest, "error while role id is invalid format", err) {
		return
	}
	loginExistanceResponse, err := h.storage.User().LoginExists(
		context.Background(),
		&models.LoginExistsRequest{
			Login: staff.Login,
		})
	if HandleHTTPError(c, http.StatusBadRequest, "error while checking existence of login", err) || HandleHTTPValidationError(c, loginExistanceResponse.Exist, http.StatusConflict, errors.New("login already exists")) {
		return
	}
	hashedPassword, err := security.HashPassword(staff.Password)
	if HandleHTTPError(c, http.StatusBadRequest, "error while hashing the password", err) {
		return
	}

	staff.Password = hashedPassword
	staff.Id = primitive.NewObjectID().Hex()
	staff.Policy = 1
	staff.UserType = "staff"
	staff.UniqueName, err = h.UniqueName(c, staff)
	if HandleHTTPError(c, http.StatusBadRequest, "error while getting unique name", err) {
		return
	}

	resp, err := h.storage.User().Create(
		context.Background(),
		staff,
	)

	if HandleHTTPError(c, http.StatusBadRequest, "error while creating staff", err) {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

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

	staff, err := h.storage.User().Get(context.Background(), &models.GetRequest{
		Id: id,
	})

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
	staff, err := h.storage.User().Get(context.Background(), &models.GetRequest{
		Id: userInfo.ID,
	})

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
		uint32(page),
		uint32(limit),
	)

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting all staffs", err) {
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Router /v1/staff/{staff_id} [put]
// @Summary Update Staff
// @Description API for updating staff
// @Tags staff
// @Accept json
// @Produce json
// @Param staff_id path string  true "staff_id"
// @Param staff body ek_user_service.CreateUpdateStaffSwag true "staff"
// @Success 200 {object} user_service.CreateResponse
func (h *handlerV1) UpdateStaff(c *gin.Context) {
	var (
		staff   *models.CreateUpdateStaff
		staffID = c.Param("staff_id")
	)

	_, err := primitive.ObjectIDFromHex(staffID)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing id, id is incorrect format", err) {
		return
	}

	if err = c.ShouldBindJSON(&staff); HandleHTTPError(c, http.StatusBadRequest, "error while updating staffs", err) {
		return
	}

	if staff.Password != "" && len(staff.Password) >= 8 {
		hashedPassword, err := security.HashPassword(staff.Password)
		if HandleHTTPError(c, http.StatusBadRequest, "error while hashing the password", err) {
			return
		}
		staff.Password = hashedPassword
	} else {
		staff.Password = ""
	}
	// staff.UniqueName, err = h.UniqueName(c, staff)
	if staff.Region != nil {
		staff.Soato = strconv.Itoa(int(staff.Region.Soato))
	} else if staff.City != nil {
		staff.Soato = strconv.Itoa(int(staff.City.Soato))
	} else {
		staff.Soato = "17"
	}
	if HandleHTTPError(c, http.StatusBadRequest, "error while getting unique name", err) {
		return
	}
	staff.Id = staffID
	resp, err := h.storage.User().Update(
		context.Background(),
		staff)

	if HandleHTTPError(c, http.StatusBadRequest, "error while updating staff", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

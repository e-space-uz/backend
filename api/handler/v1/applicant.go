package v1

import (
	"net/http"

	"github.com/e-space-uz/backend/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/applicant [post]
// @Summary Create applicant
// @Description API for creating applicant
// @Tags applicant
// @Accept json
// @Produce json
// @Param applicant body models.CreateUpdateApplicantSwag  true "applicant"
// @Success 201 {object} models.CreateResponse
func (h *handlerV1) CreateApplicant(c *gin.Context) {
	var (
		applicant models.Applicant
	)

	if err := c.ShouldBindJSON(&applicant); HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.Create.BindingApplicant", err) {
		return
	}

	applicant.ID = primitive.NewObjectID().Hex()
	resp, err := h.storage.Applicant().Create(
		c.Request.Context(),
		&applicant,
	)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.Create.ServiceInternal", err) {
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Router /v1/applicant/{applicant_id} [get]
// @Summary Get Applicant
// @Description API for getting applicant
// @Tags applicant
// @Accept json
// @Produce json
// @Param applicant_id path string  true "applicant_id"
// @Success 200 {object} models.Applicant
func (h *handlerV1) GetApplicant(c *gin.Context) {
	var (
		applicantID = c.Param("applicant_id")

		_, err = primitive.ObjectIDFromHex(applicantID)
	)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GET.ParsingApplicantObjectID", err) {
		return
	}

	applicant, err := h.storage.Applicant().Get(c.Request.Context(), applicantID)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GET.ServiceInternal", err) {
		return
	}

	c.JSON(http.StatusOK, applicant)
}

// @Router /v1/applicant [get]
// @Summary Getting All Applicants
// @Description API for getting all applicants
// @Tags applicant
// @Accept json
// @Produce json
// @Param find query models.GetAllApplicantsRequestSwag false "filters"
// @Success 200 {object} models.GetAllApplicantsResponse
func (h *handlerV1) GetAllApplicants(c *gin.Context) {
	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	applicants, count, err := h.storage.Applicant().GetAll(
		c.Request.Context(),
		uint32(page),
		uint32(limit),
	)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETALL.ServiceInternal", err) {
		return
	}

	c.JSON(http.StatusOK, models.GetAllApplicantsResponse{
		Applicants: applicants,
		Count:      int64(count),
	})
}

// @Router /v1/applicant/{applicant_id} [put]
// @Summary Update Applicant
// @Description API for updating applicant
// @Tags applicant
// @Accept json
// @Produce json
// @Param applicant_id path string  true "applicant_id"
// @Param applicant body models.CreateUpdateApplicantSwag true "applicant"
// @Success 200 {object} models.CreateResponse
func (h *handlerV1) UpdateApplicant(c *gin.Context) {
	var (
		applicant   models.Applicant
		applicantID = c.Param("applicant_id")
		_, err      = primitive.ObjectIDFromHex(applicantID)
	)

	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing object id", err) {
		return
	}

	if err = c.ShouldBindJSON(&applicant); HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.UPDATE.BindingApplicant", err) {
		return
	}
	applicant.ID = applicantID

	err = h.storage.Applicant().Update(c.Request.Context(), &applicant)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.UPDATE.ServiceInternal", err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// GetApplicantByToken
// @Security ApiKeyAuth
// @Router /v1/applicant-by-token [get]
// @Summary Get Applicant by token
// @Description API for getting applicant by token
// @Tags applicant
// @Accept json
// @Produce json
// @Success 200 {object} models.Applicant
func (h *handlerV1) GetApplicantByToken(c *gin.Context) {
	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}

	applicant, err := h.storage.Applicant().Get(c.Request.Context(), userInfo.ID)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETBYTOKEN.ServiceInternal ", err) {
		return
	}

	c.JSON(http.StatusOK, applicant)
}

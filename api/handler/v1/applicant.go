package v1

import (
	"context"
	"net/http"

	"github.com/e-space-uz/backend/pkg/security"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Router /v1/applicant [post]
// @Summary Create applicant
// @Description API for creating applicant
// @Tags applicant
// @Accept json
// @Produce json
// @Param applicant body ek_user_service.CreateUpdateApplicantSwag  true "applicant"
// @Success 201 {object} user_service.CreateResponse
func (h *handlerV1) CreateApplicant(c *gin.Context) {
	var (
		applicant us.Applicant
	)

	if err := c.ShouldBindJSON(&applicant); HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.Create.BindingApplicant", err) {
		return
	}

	applicant.Id = primitive.NewObjectID().Hex()

	resp, err := h.storage.User().Create(
		context.Background(),
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
// @Success 200 {object} user_service.Applicant
func (h *handlerV1) GetApplicant(c *gin.Context) {
	var (
		applicantID = c.Param("applicant_id")

		_, err = primitive.ObjectIDFromHex(applicantID)
	)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GET.ParsingApplicantObjectID", err) {
		return
	}

	applicant, err := h.storage.User().Get(applicantID)

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
// @Param find query ek_user_service.GetAllApplicantsRequestSwag false "filters"
// @Success 200 {object} user_service.GetAllApplicantsResponse
func (h *handlerV1) GetAllApplicants(c *gin.Context) {
	var (
		_              ek_variables.CreateResponse
		fullName       = c.Query("full_name")
		userType       = c.Query("user_type")
		phoneNumber    = c.Query("phone_number")
		passportNumber = c.Query("passport_number")
		pinfl          = c.Query("pinfl")
	)

	page, err := ParseQueryParam(c, h.log, "page", "1")
	if err != nil {
		return
	}

	limit, err := ParseQueryParam(c, h.log, "limit", "20")
	if err != nil {
		return
	}

	applicants, err := h.storage.User().GetAll(
		context.Background(),
		&us.GetAllApplicantsRequest{
			FullName:       fullName,
			PhoneNumber:    phoneNumber,
			UserType:       userType,
			PassportNumber: passportNumber,
			Pinfl:          pinfl,
			Page:           uint32(page),
			Limit:          uint32(limit),
		})

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETALL.ServiceInternal", err) {
		return
	}

	c.JSON(http.StatusOK, applicants)
}

// @Router /v1/applicant/{applicant_id} [put]
// @Summary Update Applicant
// @Description API for updating applicant
// @Tags applicant
// @Accept json
// @Produce json
// @Param applicant_id path string  true "applicant_id"
// @Param applicant body ek_user_service.CreateUpdateApplicantSwag true "applicant"
// @Success 200 {object} user_service.CreateResponse
func (h *handlerV1) UpdateApplicant(c *gin.Context) {
	var (
		applicant   us.Applicant
		applicantID = c.Param("applicant_id")
		_, err      = primitive.ObjectIDFromHex(applicantID)
	)

	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing object id", err) {
		return
	}

	if err = c.ShouldBindJSON(&applicant); HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.UPDATE.BindingApplicant", err) {
		return
	}
	applicant.Id = applicantID

	resp, err := h.storage.User().Update(
		context.Background(),
		&applicant)

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.UPDATE.ServiceInternal", err) {
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOneIDAccessToken swagger
// @Router /v1/applicant/one-id/{code} [get]
// @Summary Get Applicant one-id token
// @Description API for get applicant one-id token
// @Tags applicant
// @Accept json
// @Produce json
// @Param code path string  true "code"
// @Success 200 {object} user_service.CreateResponse
func (h *handlerV1) GetOneIDAccessToken(c *gin.Context) {
	var (
		Gender                     string
		OneIDCode                  = c.Param("code")
		TokenBody                  ek_user_service.TokenResponse
		UserBody                   ek_user_service.UserResponse
		applicantObject            *us.Applicant
		accessTokenExpireDuration  = ek_variables.AccessTokenExpireDuration
		refreshTokenExpireDuration = ek_variables.RefreshTokenExpireDuration
	)
	// h.GetOneIDTokens is in one_id.go file
	TokenBody, err := h.GetOneIDTokens(OneIDCode)
	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETONEIDACCESSTOKEN.GetTokenByCode -> "+TokenBody.Message, err) {
		return
	}

	// h.GetOneIDUser is in one_id.go file
	UserBody, err = h.GetOneIDUser(TokenBody.AccessToken)
	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETONEIDACCESSTOKEN.GetUserByToken -> "+TokenBody.Message, err) {
		return
	}
	// check user exists or not
	userIDExist, err := h.storage.User().Exists(
		context.Background(),
		&us.USExistsRequest{
			Id: UserBody.UserID,
		},
	)
	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETONEIDACCESSTOKEN.ExistsMethod", err) {
		return
	}

	if UserBody.Gd == "1" {
		Gender = "male"
	} else {
		Gender = "female"
	}
	//
	if userIDExist.Exist {
		applicantObject, err = h.storage.User().GetByUserId(
			context.Background(),
			&us.GetRequest{
				Id: UserBody.UserID,
			},
		)
		if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETONEIDACCESSTOKEN.GetUserID", err) {
			return
		}
	} else {
		applicantObject, err = h.storage.User().Create(
			context.Background(),
			&us.Applicant{
				Id:                 primitive.NewObjectID().Hex(),
				BirthDate:          UserBody.BirthDate,
				BirthPlace:         UserBody.BirthPlace,
				Citizenship:        UserBody.Ctzn,
				PermanentAddress:   UserBody.PerAdr,
				PassportIssuePlace: UserBody.PportIssuePlace,
				LastName:           UserBody.SurName,
				Gender:             Gender,
				Nationality:        UserBody.Natn,
				PassportIssueDate:  UserBody.PportIssueDate,
				PassportExpiryDate: UserBody.PportExprDate,
				PassportNumber:     UserBody.PportNo,
				Pin:                UserBody.Pin,
				PhoneNumber:        UserBody.MobPhoneNo,
				Login:              UserBody.UserID,
				Email:              UserBody.Email,
				MiddleName:         UserBody.MidName,
				UserType:           UserBody.UserType,
				FirstName:          UserBody.FirstName,
				FullName:           UserBody.FullName,
			},
		)
		if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETONEIDACCESSTOKEN.Create", err) {
			return
		}
	}
	m := map[string]interface{}{
		"id":        applicantObject.Id,
		"login":     applicantObject.Login,
		"user_type": applicantObject.UserType,
		"full_name": applicantObject.FullName,
	}

	accessToken, err := security.GenerateJWT(m, accessTokenExpireDuration, h.cfg.LoginSecretAccessKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "UserService.Applicant.GETONEIDACCESSTOKEN.GenerateAccessToken", err) {
		return
	}

	refreshToken, err := security.GenerateJWT(m, refreshTokenExpireDuration, h.cfg.LoginSecretRefreshKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "UserService.Applicant.GETONEIDACCESSTOKEN.GenerateRefreshToken", err) {
		return
	}

	c.JSON(http.StatusOK, &us.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// GetApplicantByToken
// @Security ApiKeyAuth
// @Router /v1/applicant-by-token [get]
// @Summary Get Applicant by token
// @Description API for getting applicant by token
// @Tags applicant
// @Accept json
// @Produce json
// @Success 200 {object} user_service.Applicant
func (h *handlerV1) GetApplicantByToken(c *gin.Context) {
	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}

	applicant, err := h.storage.User().Get(context.Background(), &us.GetRequest{
		Id: userInfo.ID,
	})

	if HandleHTTPError(c, http.StatusBadRequest, "UserService.Applicant.GETBYTOKEN.ServiceInternal ", err) {
		return
	}

	c.JSON(http.StatusOK, applicant)
}

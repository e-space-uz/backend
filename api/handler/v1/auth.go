package v1

import (
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/pkg/security"
	"github.com/gin-gonic/gin"
)

// @Router /v1/login [post]
// @Summary login
// @Description API to singin
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest  true "Login"
// @Success 201 {object} models.LoginResponse

func (h *handlerV1) Login(c *gin.Context) {
	var (
		login                      models.LoginRequest
		accessTokenExpireDuration  = config.AccessTokenExpireDuration
		refreshTokenExpireDuration = config.RefreshTokenExpireDuration
	)

	if err := c.ShouldBindJSON(&login); HandleHTTPError(c, http.StatusBadRequest, "DiscussionLogicService.Action.Create.BindingAction", err) {
		return
	}
	if len(login.Login) < 6 {
		HandleHTTPError(c, http.StatusBadRequest, "error validating request",
			errors.New("please, provide valid login"))
		return
	}
	if len(login.Password) < 8 {
		HandleHTTPError(c, http.StatusBadRequest, "error validating request",
			errors.New("please, provide valid password"))
		return
	}

	loginResponse, err := h.storage.Staff().Login(context.Background(), &models.LoginRequest{
		Login: login.Login,
	})
	if HandleHTTPError(c, http.StatusBadRequest, "error while getting login info", err) {
		return
	}
	match, err := security.ComparePassword(loginResponse.Password, login.Password)
	if err != nil {
		HandleHTTPError(c, http.StatusUnauthorized, "password does not match", err)
		return
	}
	if !match {
		HandleHTTPError(c, http.StatusUnauthorized, "password does not match", errors.New("provided password does not match"))
		return
	}

	m := map[string]interface{}{
		"id":              loginResponse.Id,
		"login":           loginResponse.Login,
		"user_type":       loginResponse.UserType,
		"full_name":       loginResponse.FullName,
		"role_id":         loginResponse.RoleId,
		"soato":           loginResponse.Soato,
		"organization_id": loginResponse.Role.Organization.Id,
	}

	accessToken, err := security.GenerateJWT(m, accessTokenExpireDuration, h.cfg.LoginSecretAccessKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "Error while generating access token", err) {
		return
	}

	refreshToken, err := security.GenerateJWT(m, refreshTokenExpireDuration, h.cfg.LoginSecretRefreshKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "Error while generating refresh token", err) {
		return
	}
	response := &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Verified:     loginResponse.Verified,
		Role:         loginResponse.Role,
		Soato:        loginResponse.Soato,
	}
	c.JSON(http.StatusOK, response)
}

// @Router /v1/login-exists [post]
// @Summary login
// @Description API to singin
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.LoginExistsRequest  true "Login"
// @Success 201 {object} models.LoginExistsResponse

func (h *handlerV1) LoginExist(c *gin.Context) {
	var (
		loginExist models.LoginExistsRequest
	)
	if err := c.ShouldBindJSON(&loginExist); HandleHTTPError(c, http.StatusBadRequest, "DiscussionLogicService.Action.Create.BindingAction", err) {
		return
	}
	loginExistanceResponse, err := h.storage.Staff().LoginExists(
		context.Background(),
		&models.LoginExistsRequest{
			Login: loginExist.Login,
		})
	if HandleHTTPError(c, http.StatusBadRequest, "error while checking login", err) {
		return
	}
	c.JSON(http.StatusOK, loginExistanceResponse)
}

// @Router /v1/login-refresh [post]
// @Summary if access-token expired, get your access token with refresh
// @Description API to get your access token with refresh
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token query string  true "refresh_token"
// @Param is_applicant query string  false "is_applicant"
// @Success 201 {object} models.LoginResponse

func (h *handlerV1) LoginRefresh(c *gin.Context) {
	var (
		token                      = c.Query("refresh_token")
		isApplicantQuery           = c.Query("is_applicant")
		accessTokenExpireDuration  = models.AccessTokenExpireDuration
		refreshTokenExpireDuration = models.RefreshTokenExpireDuration
		m                          map[string]interface{}
		response                   = &models.LoginResponse{}
	)
	claims, err := security.ExtractClaims(token, h.cfg.LoginSecretRefreshKey)
	if err != nil {
		HandleHTTPError(c, http.StatusBadRequest, "please provide token", errors.New("incorrect token format"))
		return
	}
	if isApplicantQuery != "" {
		m = map[string]interface{}{
			"id":        claims["id"],
			"login":     claims["login"],
			"user_type": claims["user_type"],
			"full_name": claims["full_name"],
		}
	} else {
		loginResponse, err := h.storage.Staff().Login(context.Background(), &models.LoginRequest{
			Login: claims["login"].(string),
		})
		if HandleHTTPError(c, http.StatusBadRequest, "error while getting login info", err) {
			return
		}
		m = map[string]interface{}{
			"id":        claims["id"],
			"login":     claims["login"],
			"user_type": claims["user_type"],
			"full_name": claims["full_name"],
			"role_id":   claims["role_id"],
			"soato":     claims["soato"],
		}
		response = &models.LoginResponse{
			Verified: loginResponse.Verified,
			Role:     loginResponse.Role,
			Soato:    loginResponse.Soato,
		}
	}

	accessToken, err := security.GenerateJWT(m, accessTokenExpireDuration, h.cfg.LoginSecretAccessKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "Error while generating token", err) {
		return
	}

	refreshToken, err := security.GenerateJWT(m, refreshTokenExpireDuration, h.cfg.LoginSecretRefreshKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "Error while generating token", err) {
		return
	}
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Router /v1/update-password/{user_id} [post]
// @Summary update user password
// @Description API to update user password
// @Tags auth
// @Accept json
// @Produce json
// @Param user_id path string true "User id"
// @Param password body models.UpdatePassword true "Update password"
// @Success 201 {object} models.EmptyResponse

func (h *handlerV1) UpdatePassword(c *gin.Context) {
	var (
		UpdateBody                 models.UpdatePassword
		accessTokenExpireDuration  = models.AccessTokenExpireDuration
		refreshTokenExpireDuration = models.RefreshTokenExpireDuration
		userID                     = c.Param("user_id")
	)
	if err := c.BindJSON(&UpdateBody); HandleHTTPError(c, http.StatusBadRequest, "UserService.Auth.UpdatePassword", err) {
		return
	}

	_, err := primitive.ObjectIDFromHex(userID)
	if HandleHTTPError(c, http.StatusBadRequest, "error while parsing id", err) {
		return
	}

	if len(UpdateBody.NewPassword) < 8 {
		HandleHTTPError(c, http.StatusBadRequest, "error validating request",
			errors.New("please, provide valid password"))
		return
	}

	staffResponse, err := h.storage.Staff().Get(context.Background(), &models.GetRequest{
		Id: userID,
	})

	if HandleHTTPError(c, http.StatusBadRequest, "error while getting login info", err) {
		return
	}
	match, err := security.ComparePassword(staffResponse.Password, UpdateBody.OldPassword)
	if err != nil {
		HandleHTTPError(c, http.StatusInternalServerError, "UserService.Auth.UpdatePassword", err)
		return
	}
	if !match {
		HandleHTTPError(c, http.StatusConflict, "old password does not match", errors.New("old provided password does not match"))
		return
	}

	hashedPassword, err := security.HashPassword(UpdateBody.NewPassword)
	if HandleHTTPError(c, http.StatusBadRequest, "error while hashing the password", err) {
		return
	}
	UpdateBody.NewPassword = hashedPassword

	_, err = h.storage.Staff().UpdatePassword(
		context.Background(),
		&models.PasswordUpdateRequest{
			OldPassword: UpdateBody.OldPassword,
			NewPassword: UpdateBody.NewPassword,
			UserId:      staffResponse.Id,
		})
	if HandleHTTPError(c, http.StatusBadRequest, "error while updating password", err) {
		return
	}

	m := map[string]interface{}{
		"id":        staffResponse.Id,
		"login":     staffResponse.Login,
		"user_type": staffResponse.UserType,
		"role_id":   staffResponse.Role.Id,
		"soato":     staffResponse.Soato,
	}

	accessToken, err := security.GenerateJWT(m, accessTokenExpireDuration, h.cfg.LoginSecretAccessKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "Error while generating token", err) {
		return
	}

	refreshToken, err := security.GenerateJWT(m, refreshTokenExpireDuration, h.cfg.LoginSecretRefreshKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "Error while generating token", err) {
		return
	}

	response := &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Verified:     staffResponse.Verified,
		Soato:        staffResponse.Soato,
	}
	c.JSON(http.StatusOK, response)
}

// @Security ApiKeyAuth
// @Router /v1/update-password [post]
// @Summary update user password
// @Description API to update user password
// @Tags auth
// @Accept json
// @Produce json
// @Param password body models.UpdatePassword true "Update password"
// @Success 201 {object} models.EmptyResponse

func (h *handlerV1) UpdatePasswordFromToken(c *gin.Context) {
	var (
		UpdateBody                 models.UpdatePassword
		accessTokenExpireDuration  = models.AccessTokenExpireDuration
		refreshTokenExpireDuration = models.RefreshTokenExpireDuration
	)
	if err := c.BindJSON(&UpdateBody); HandleHTTPError(c, http.StatusBadRequest, "UserService.Auth.UpdatePassword", err) {
		return
	}

	if len(UpdateBody.NewPassword) < 8 {
		HandleHTTPError(c, http.StatusBadRequest, "error validating request",
			errors.New("please, provide valid password"))
		return
	}

	userInfo, err := h.UserInfo(c, true)
	if err != nil {
		return
	}

	loginResponse, err := h.storage.Staff().Login(context.Background(), &models.LoginRequest{
		Login: userInfo.Login,
	})
	if HandleHTTPError(c, http.StatusBadRequest, "error while getting login info", err) {
		return
	}
	match, err := security.ComparePassword(loginResponse.Password, UpdateBody.OldPassword)
	if err != nil {
		HandleHTTPError(c, http.StatusInternalServerError, "UserService.Auth.UpdatePassword", err)
		return
	}
	if !match {
		HandleHTTPError(c, http.StatusConflict, "old password does not match", errors.New("old provided password does not match"))
		return
	}

	hashedPassword, err := security.HashPassword(UpdateBody.NewPassword)
	if HandleHTTPError(c, http.StatusBadRequest, "error while hashing the password", err) {
		return
	}
	UpdateBody.NewPassword = hashedPassword

	_, err = h.storage.Staff().UpdatePassword(
		context.Background(),
		&models.PasswordUpdateRequest{
			OldPassword: UpdateBody.OldPassword,
			NewPassword: UpdateBody.NewPassword,
			UserId:      userInfo.ID,
		})
	if HandleHTTPError(c, http.StatusBadRequest, "error while updating password", err) {
		return
	}

	m := map[string]interface{}{
		"id":        loginResponse.Id,
		"login":     loginResponse.Login,
		"user_type": loginResponse.UserType,
		"full_name": loginResponse.FullName,
		"role_id":   loginResponse.RoleId,
		"soato":     loginResponse.Soato,
	}

	accessToken, err := security.GenerateJWT(m, accessTokenExpireDuration, h.cfg.LoginSecretAccessKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "Error while generating token", err) {
		return
	}

	refreshToken, err := security.GenerateJWT(m, refreshTokenExpireDuration, h.cfg.LoginSecretRefreshKey)
	if HandleHTTPError(c, http.StatusInternalServerError, "Error while generating token", err) {
		return
	}

	response := &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Verified:     loginResponse.Verified,
		Role:         loginResponse.Role,
		Soato:        loginResponse.Soato,
	}

	c.JSON(http.StatusOK, response)
}

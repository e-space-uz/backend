package v1

import (
	"bufio"
	"crypto/rand"
	"errors"
	"math/big"
	"net/http"
	"net/url"
	"strconv"

	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/models"
	"github.com/e-space-uz/backend/pkg/logger"
	"github.com/e-space-uz/backend/pkg/security"
	"github.com/e-space-uz/backend/storage"
	"github.com/gin-gonic/gin"
)

var (
	ErrAlreadyExists       = "ALREADY_EXISTS"
	ErrBadRequest          = "BAD_REQUEST"
	ErrNotFound            = "NOT_FOUND"
	ErrInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable  = "SERVICE_UNAVAILABLE"
	log                    = logger.New("DEBUG", "ek_admin_api_gateway")
)

const (
	//ErrorBadRequest ...
	ErrorBadRequest = "BAD_REQUEST"
)

type handlerV1 struct {
	cfg     config.Config
	log     logger.Logger
	storage storage.StorageI
}

type HandlerV1Options struct {
	Cfg     config.Config
	Log     logger.Logger
	Storage storage.StorageI
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		log:     options.Log,
		cfg:     options.Cfg,
		storage: options.Storage,
	}
}

//ParseActiveQueryParam ...
func ParseActiveQueryParam(c *gin.Context) (bool, error) {
	a, err := strconv.ParseBool(c.DefaultQuery("active", "false"))
	if err != nil {
		return false, err
	}
	return a, nil
}

//ParseInactiveQueryParam ...
func ParseInactiveQueryParam(c *gin.Context) (bool, error) {
	a, err := strconv.ParseBool(c.DefaultQuery("inactive", "false"))
	if err != nil {
		return false, err
	}
	return a, nil
}

func (h *handlerV1) MakeProxy(c *gin.Context, proxyUrl, path string) (err error) {
	req := c.Request

	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		h.log.Error("error in parse addr: %v", logger.Error(err))
		c.String(http.StatusInternalServerError, "error")
		return
	}

	req.URL.Scheme = proxy.Scheme
	req.URL.Host = proxy.Host
	req.URL.Path = path
	transport := http.DefaultTransport

	resp, err := transport.RoundTrip(req)
	if HandleHTTPError(c, 500, "error in round trip:", err) {
		return
	}

	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Header(k, v)
		}
	}
	defer resp.Body.Close()

	c.Status(resp.StatusCode)
	_, _ = bufio.NewReader(resp.Body).WriteTo(c.Writer)
	return
}
func (h *handlerV1) UserInfo(c *gin.Context, required bool) (*models.LoginInfo, error) {
	var token = c.Request.Header.Get("Authorization")
	if token != "" {
		claims, err := security.ExtractClaims(token, h.cfg.LoginSecretAccessKey)
		if err != nil {
			HandleHTTPError(c, http.StatusUnauthorized, "please provide valid token", errors.New("unauthorized"))
			return &models.LoginInfo{}, err
		}
		if claims["user_type"].(string) == "staff" {
			return &models.LoginInfo{
				ID:       claims["id"].(string),
				Login:    claims["login"].(string),
				UserType: claims["user_type"].(string),
				Soato:    claims["soato"].(string),
			}, nil
		} else {
			return &models.LoginInfo{
				ID:       claims["id"].(string),
				Login:    claims["login"].(string),
				UserType: claims["user_type"].(string),
				FullName: claims["full_name"].(string),
			}, nil
		}
	} else if required {
		HandleHTTPError(c, http.StatusUnauthorized, "please provide token", errors.New("unauthorized"))
		return &models.LoginInfo{}, errors.New("not authorized")
	} else {
		return &models.LoginInfo{}, nil
	}
}

// This function is to parsing upcoming query params
// please, make sure you write proper message error
func ParseQueryParam(c *gin.Context, log logger.Logger, key string, defaultValue string) (int, error) {
	valueString := c.DefaultQuery(key, defaultValue)

	value, err := strconv.Atoi(valueString)
	if err != nil {
		message := "Error while parsing query"
		log.Error(message, logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		c.Abort()
		return 0, err
	}
	return value, err
}

// This function is to handle all HTTP client errors
// please, make sure you write proper message error
func HandleHTTPError(c *gin.Context, code int, message string, err error) bool {
	if err != nil && code == http.StatusBadRequest {
		log.Error(message+" --> Error: ", logger.Error(err))
		c.JSON(http.StatusBadRequest, models.FailureResponse{
			Success: false,
			Message: message,
			Error:   err,
		})
		return true
	} else if err != nil && code == http.StatusInternalServerError {
		log.Error(message+", --> Error	", logger.Error(err))
		c.JSON(http.StatusServiceUnavailable, models.FailureResponse{
			Success: false,
			Message: message,
			Error:   err,
		})
		return true
	} else if err != nil && code == http.StatusUnauthorized {
		log.Error(message+" --> Error: ", logger.Error(err))
		c.JSON(http.StatusUnauthorized, models.FailureResponse{
			Success: false,
			Message: message,
			Error:   err,
		})
		return true
	} else if err != nil && code == http.StatusConflict {
		log.Error(message+" --> Error: ", logger.Error(err))
		c.JSON(http.StatusConflict, models.FailureResponse{
			Success: false,
			Message: message,
			Error:   err,
		})
		return true
	}
	return false
}
func RandomSixDigits() (uint64, error) {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return n.Uint64(), nil
}

// This function is to handle all HTTP client validation errors
// please, make sure you write proper message error
func HandleHTTPValidationError(c *gin.Context, exist bool, code int, err error) bool {
	if !exist && code == http.StatusNotFound {
		log.Error(" --> Error: ", logger.Error(err))
		c.JSON(http.StatusNotFound, models.FailureResponse{
			Success: false,
			Message: err.Error(),
			Error:   err,
		})
		return true
	} else if exist && code == http.StatusConflict {
		log.Error(" --> Error: ", logger.Error(err))
		c.JSON(http.StatusConflict, models.FailureResponse{
			Success: false,
			Message: err.Error(),
			Error:   err,
		})
		return true
	} else if !exist && code == http.StatusInternalServerError {
		log.Error(", --> Error	", logger.Error(err))
		c.JSON(http.StatusServiceUnavailable, models.FailureResponse{
			Success: false,
			Message: err.Error(),
			Error:   err,
		})
		return true
	}
	return false
}

// func getUserInfo(c *gin.Context) (userInfo models.UserInfo) {
// 	tokenStr := c.GetHeader("Authorization")

// 	if tokenStr == "" {
// 		c.Status(http.StatusUnauthorized)
// 		return
// 	}

// 	claims, err := jwt.ExtractClaims(tokenStr, SigningKey)

// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	userInfo.UserID = claims["sub"].(string)
// 	userInfo.UserType = claims["user_type"].(string)
// 	userInfo.RoleID = claims["role"].(string)

// 	return
// }

func ParseBoolParam(c *gin.Context, log logger.Logger, key, defaultValue string) (bool, error) {
	valueString := c.DefaultQuery(key, defaultValue)

	value, err := strconv.ParseBool(valueString)
	if err != nil {
		message := "Error while parsing query"
		log.Error(message, logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
		c.Abort()
		return false, err
	}
	return value, err
}

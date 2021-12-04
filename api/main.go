// Note: please follow rules

package api

import (
	"net/http"

	_ "github.com/e-space-uz/backend/api/docs"
	v1 "github.com/e-space-uz/backend/api/handler/v1"
	"github.com/e-space-uz/backend/config"
	"github.com/e-space-uz/backend/pkg/logger"
	"github.com/e-space-uz/backend/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterOptions struct {
	Log     logger.Logger
	Cfg     config.Config
	Storage storage.StorageI
}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RouterOptions) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "*")

	router.Use(cors.New(corsConfig))

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Log:     opt.Log,
		Cfg:     opt.Cfg,
		Storage: opt.Storage,
	})
	routesV1 := router.Group("/v1")
	routesV1.Use()
	{
		routesV1.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    "ekadastr api_gateway",
			})
		})

		routesV1.POST("/file-upload", handlerV1.FileUpload)
		routesV1.POST("/image-upload", handlerV1.ImageUpload)

		//City endpoints
		routesV1.GET("/city/:city_id", handlerV1.GetCity)
		routesV1.GET("/city", handlerV1.GetAllCities)

		routesV1.GET("/region/:region_id", handlerV1.GetRegion)
		routesV1.GET("/region", handlerV1.GetAllRegions)
		routesV1.GET("/regions/:city_id", handlerV1.GetAllRegionsByCityID)
		//District Endpoints
		routesV1.GET("/district/:district_id", handlerV1.GetDistrict)
		routesV1.GET("/district", handlerV1.GetAllDistricts)

		routesV1.POST("/login", handlerV1.Login)
		routesV1.POST("/login-exists", handlerV1.LoginExist)
		routesV1.POST("/login-refresh", handlerV1.LoginRefresh)

		//Applicant endpoints
		routesV1.POST("/applicant", handlerV1.CreateApplicant)
		routesV1.GET("/applicant/:applicant_id", handlerV1.GetApplicant)
		routesV1.GET("/applicant", handlerV1.GetAllApplicants)
		routesV1.GET("/applicant-by-token", handlerV1.GetApplicantByToken)
		routesV1.PUT("/applicant/:applicant_id", handlerV1.UpdateApplicant)
		routesV1.GET("/applicant/one-id/:code", handlerV1.GetOneIDAccessToken)

		//Staff endpoints
		routesV1.POST("/staff", handlerV1.CreateStaff)
		routesV1.GET("/staff/:staff_id", handlerV1.GetStaff)
		routesV1.GET("/staff", handlerV1.GetAllStaffs)
		routesV1.GET("/staff-by-token", handlerV1.GetStaffByToken)
		routesV1.PUT("/staff/:staff_id", handlerV1.UpdateStaff)

		//Entity endpoints
		routesV1.POST("/entity", handlerV1.CreateEntity)
		routesV1.GET("/entity/:entity_id", handlerV1.GetEntity)
		routesV1.GET("/entity", handlerV1.GetAllEntities)
		routesV1.GET("/entity-properties", handlerV1.GetAllEntitiesWithProperties)

		//Entity Draft endpoints
		routesV1.POST("/entity-draft", handlerV1.CreateEntityDraft)
		routesV1.GET("/entity-draft/:entity_draft_id", handlerV1.GetEntityDraft)
		routesV1.GET("/entity-draft", handlerV1.GetAllEntityDraft)

		//Property endpoints
		routesV1.POST("/property", handlerV1.CreateProperty)
		routesV1.GET("/property/:property_id", handlerV1.GetProperty)
		routesV1.GET("/property", handlerV1.GetAllProperties)
		routesV1.PUT("/property/:property_id", handlerV1.UpdateProperty)
		routesV1.DELETE("/property/:property_id", handlerV1.DeleteProperty)

		// Group property endpoints
		routesV1.GET("/group-property/:group_property_id", handlerV1.GetGroupProperty)
		routesV1.GET("/group-property-status/:status_id", handlerV1.GetGroupPropertyByStatusID)
		routesV1.GET("/group-property", handlerV1.GetAllGroupProperties)
		routesV1.GET("/group-property-type", handlerV1.GetAllGroupPropertiesByType)
		routesV1.POST("/group-property", handlerV1.CreateGroupProperty)
		routesV1.PUT("/group-property/:group_property_id", handlerV1.UpdateGroupProperty)
	}

	// swagger
	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router

}

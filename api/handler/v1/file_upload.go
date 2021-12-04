package v1

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type FilesResponse struct {
	FilePath string `json:"file_path"`
	URL      string `json:"file_url"`
}
type ImageResponse struct {
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"`
}

type AllFilesResponse struct {
	Count int64           `json:"count"`
	Files []FilesResponse `json:"files"`
}

// func validation(ext string) bool {
// 	// for _, val := range allowedExtensions {
// 	// 	if val == ext {
// 	// 		return true
// 	// 	}
// 	// }
// 	return false
// }

// @Router /v1/file-upload [post]
// @Summary Create file
// @Description API for creating file
// @Tags file_upload
// @Param file formData file true "file"
// @Param region body ek_entity_service.EntityFilesSwag  true "region"
// @Accept multipart/form-data
// @Accept json
// @Produce json
// @Success 200 {object} FilesResponse

func (h *handlerV1) FileUpload(c *gin.Context) {
	var (
		fileURL     string
		entityFiles ek_entity_service.EntityFilesSwag
	)

	if err := c.ShouldBind(&entityFiles); HandleHTTPError(c, http.StatusBadRequest, "DiscussionLogicService.Action.Create.BindingAction", err) {
		return
	}

	fmt.Println("asas", entityFiles)
	minioClient, err := minio.New(h.cfg.MinioDomain, &minio.Options{
		Creds:  credentials.NewStaticV4(h.cfg.MinioAccessKeyID, h.cfg.MinioSecretAccesKey, ""),
		Secure: true,
	})
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}

	exists, _ := minioClient.BucketExists(context.Background(), h.cfg.BucketName)
	if !exists {
		err = minioClient.MakeBucket(context.Background(), h.cfg.BucketName, minio.MakeBucketOptions{Region: h.cfg.MinioDomain})
		if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
			return
		}
	}
	file, err := c.FormFile("file")
	if HandleHTTPError(c, http.StatusInternalServerError, "error while getting file", err) {
		return
	}
	contentType := file.Header["Content-Type"][0]
	if len(contentType) <= 1 {
		contentType = "multipart/form-data"
	}
	fileID := primitive.NewObjectID().Hex()
	fileURL = fileID + filepath.Ext(file.Filename)
	object, err := file.Open()
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}
	defer object.Close()
	_, err = minioClient.PutObject(context.Background(), h.cfg.BucketName, fileURL, object, file.Size,
		minio.PutObjectOptions{ContentType: contentType})
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}
	_, err = h.storage.EntityFilesService().Create(context.Background(),
		&entity_service.EntityFiles{
			Id:       fileID,
			Url:      fileURL,
			FileName: file.Filename,
			Comment:  entityFiles.Comment,
			User:     primitive.NewObjectID().Hex(),
		})
	if HandleHTTPError(c, "error while creating bucket in minio", err) {
		return
	}

	c.JSON(http.StatusOK, FilesResponse{
		FilePath: fileURL,
		URL:      h.cfg.MinioDomain + "/" + h.cfg.BucketName + "/" + file.Filename,
	})
}

// @Router /v1/region/{region_id}/upload [post]
// @Tags file_upload
// @Param file formData file true "file"
// @Param region_id path string  true "region_id"
// @Param comment body ek_setting_service.RegionFilesSwag  true "body"
// @Accept multipart/form-data
// @Success 200 {object} FilesResponse
func (h *handlerV1) UploadMainFile(c *gin.Context) {
	var (
		fileURL     string
		regionFiles ek_setting_service.RegionFilesSwag
		regionID    = c.Param("region_id")
	)

	objectID, err := primitive.ObjectIDFromHex(regionID)
	if HandleHTTPError(c, http.StatusBadRequest, "RegionFiles.Action.Create.ObjectID", err) {
		return
	}

	// resp, err := h.storage.RegionService().RegionExists(c, &setting_service.SSExistsRequest{
	// 	Id: objectID.Hex(),
	// })

	// if HandleHTTPError(c, "Region.Action.Exists.RegionID", err) {
	// 	return
	// }

	// if !resp.Exist {
	// 	HandleHTTPError(c, http.StatusBadRequest, "Region.Action.Exists.RegionID", errors.New("region is not exist"))
	// 	return
	// }

	if err := c.ShouldBind(&regionFiles); HandleHTTPError(c, http.StatusBadRequest, "RegionFiles.Action.Create.BindingAction", err) {
		return
	}

	fmt.Println(regionFiles.Comment, " ++", regionFiles.Type)

	minioClient, err := minio.New(h.cfg.MinioDomain, &minio.Options{
		Creds:  credentials.NewStaticV4(h.cfg.MinioAccessKeyID, h.cfg.MinioSecretAccesKey, ""),
		Secure: false,
	})
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}

	exists, _ := minioClient.BucketExists(context.Background(), h.cfg.BucketName)
	if !exists {
		err = minioClient.MakeBucket(context.Background(), h.cfg.BucketName, minio.MakeBucketOptions{Region: h.cfg.MinioDomain})
		if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
			return
		}
	}
	file, err := c.FormFile("file")
	if HandleHTTPError(c, http.StatusInternalServerError, "error while getting file", err) {
		return
	}
	contentType := file.Header["Content-Type"][0]
	if len(contentType) <= 1 {
		contentType = "multipart/form-data"
	}
	fileID := primitive.NewObjectID().Hex()
	fileURL = fileID + filepath.Ext(file.Filename)
	object, err := file.Open()
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}
	defer object.Close()
	_, err = minioClient.PutObject(context.Background(), h.cfg.BucketName, fileURL, object, file.Size,
		minio.PutObjectOptions{ContentType: contentType})
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}
	_, err = h.storage.RegionFilesService().Create(context.Background(),
		&setting_service.RegionFiles{
			Id:       fileID,
			RegionId: objectID.Hex(),
			Url:      h.cfg.MinioDomain + "/" + h.cfg.BucketName + "/" + fileURL,
			FileName: file.Filename,
			Comment:  regionFiles.Comment,
			User:     primitive.NewObjectID().Hex(),
		})
	if HandleHTTPError(c, "error while creating bucket in minio", err) {
		return
	}

	c.JSON(http.StatusOK, FilesResponse{
		FilePath: fileURL,
		URL:      h.cfg.MinioDomain + "/" + h.cfg.BucketName + "/" + file.Filename,
	})
}

// @Router /v1/image-upload [post]
// @Tags file_upload
// @Param image formData file true "image"
// @Accept multipart/form-data
// @Success 200 {object} FilesResponse
func (h *handlerV1) ImageUpload(c *gin.Context) {

	minioClient, err := minio.New(h.cfg.MinioDomain, &minio.Options{
		Creds:  credentials.NewStaticV4(h.cfg.MinioAccessKeyID, h.cfg.MinioSecretAccesKey, ""),
		Secure: true,
	})
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}

	exists, _ := minioClient.BucketExists(context.Background(), h.cfg.BucketName)
	if !exists {
		err = minioClient.MakeBucket(context.Background(), h.cfg.BucketName, minio.MakeBucketOptions{Region: h.cfg.MinioDomain})
		if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
			return
		}
	}
	file, err := c.FormFile("image")
	if HandleHTTPError(c, http.StatusInternalServerError, "error while getting file", err) {
		return
	}
	contentType := file.Header["Content-Type"][0]
	if len(contentType) <= 1 {
		contentType = "multipart/form-data"
	}
	fileName := file.Filename
	file.Filename = primitive.NewObjectID().Hex() + filepath.Ext(fileName)
	object, err := file.Open()
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}
	defer object.Close()
	_, err = minioClient.PutObject(context.Background(), h.cfg.BucketName, file.Filename, object, file.Size,
		minio.PutObjectOptions{ContentType: contentType})
	if HandleHTTPError(c, http.StatusInternalServerError, "error while creating bucket in minio", err) {
		return
	}

	c.JSON(http.StatusOK, FilesResponse{
		FilePath: file.Filename,
		URL:      h.cfg.MinioDomain + "/" + h.cfg.BucketName + "/" + file.Filename,
	})
}

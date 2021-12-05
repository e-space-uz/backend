package v1

import (
	"context"
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

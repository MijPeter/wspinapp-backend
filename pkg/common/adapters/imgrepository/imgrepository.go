package imgrepository

import (
	"context"
	"example/wspinapp-backend/pkg/common"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"log"
	"time"
)

type cloudinaryConfig struct {
	cloudName    string
	apiKey       string
	apiSecret    string
	uploadFolder string
}

type ImageRepository interface {
	Upload(input interface{}) (string, error)
	// RemoteUpload(url Url) (string, error) https://dev.to/hackmamba/robust-media-upload-with-golang-and-cloudinary-gin-gonic-version-54ii
}

// implementation of ImageRepository
type cloudinaryRepository struct {
	cloudinary       *cloudinary.Cloudinary
	cloudinaryConfig *cloudinaryConfig
}

// Upload uploads validated image file to image repository
func (imageStore *cloudinaryRepository) Upload(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // todo the fuck is that?
	defer cancel()

	//upload file
	uploadParam, err := imageStore.cloudinary.Upload.Upload(ctx, input, uploader.UploadParams{Folder: imageStore.cloudinaryConfig.uploadFolder})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}

// New creates a cloudinary config (from environment variables file) and connects to cloudinary.
func New() ImageRepository {
	cloudinaryConfig := getCloudinaryConfig()

	return &cloudinaryRepository{
		cloudinary:       connectCloudinary(cloudinaryConfig),
		cloudinaryConfig: cloudinaryConfig,
	}
}

func connectCloudinary(cfg *cloudinaryConfig) *cloudinary.Cloudinary {
	for {
		cld, err := cloudinary.NewFromParams(cfg.cloudName, cfg.apiKey, cfg.apiSecret)

		if err == nil {
			log.Println("Successfully connected to cloudinary.")
			return cld
		}
		log.Println("Failed to connect to cloudinary, retrying in a moment.")
		time.Sleep(1 * time.Second)
	}
}

func getCloudinaryConfig() *cloudinaryConfig {
	return &cloudinaryConfig{
		cloudName:    common.EnvCloudName(),
		apiKey:       common.EnvCloudAPIKey(),
		apiSecret:    common.EnvCloudAPISecret(),
		uploadFolder: common.EnvCloudUploadFolder(),
	}
}

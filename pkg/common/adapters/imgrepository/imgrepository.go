package imgrepository

import (
	"context"
	"example/wspinapp-backend/pkg/common"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
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

type Asset struct {
	ID  string
	Url string
}

type ImageRepository interface {
	Upload(input interface{}) (string, error)
	// RemoteUpload(url Url) (string, error) https://dev.to/hackmamba/robust-media-upload-with-golang-and-cloudinary-gin-gonic-version-54ii
	Assets() ([]Asset, error)
	DeleteAssets(assetsID []string) error
}

// implementation of ImageRepository
type cloudinaryRepository struct {
	cloudinary       *cloudinary.Cloudinary
	cloudinaryConfig *cloudinaryConfig
}

// Upload uploads validated image file to image repository
func (imageStore *cloudinaryRepository) Upload(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//upload file
	uploadParam, err := imageStore.cloudinary.Upload.Upload(ctx, input, uploader.UploadParams{Folder: imageStore.cloudinaryConfig.uploadFolder})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}

// TODO delete assets from wall folder only
func (imageStore *cloudinaryRepository) Assets() ([]Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var resources []Asset
	assets, err := imageStore.cloudinary.Admin.Assets(ctx, admin.AssetsParams{MaxResults: 500})
	if err != nil {
		return nil, err
	}

	for _, asset := range assets.Assets {
		resources = append(resources, Asset{ID: asset.PublicID, Url: asset.SecureURL})
	}

	return resources, nil
}

func (imageStore *cloudinaryRepository) DeleteAssets(assetsID []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(assetsID) > 100 {
		log.Println("Too many assets to delete, choosing only 100 to delete") // cloudinary lets us delete only up to 100 assets at once
		assetsID = assetsID[0:100]
	}

	_, err := imageStore.cloudinary.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{PublicIDs: assetsID})
	return err
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

package common

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnvironmentVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Couldn't load environment variables. Exiting. %s\n", err.Error())
	}
}

func EnvCloudName() string {
	return os.Getenv("CLOUDINARY_CLOUD_NAME")
}

func EnvCloudAPIKey() string {
	return os.Getenv("CLOUDINARY_API_KEY")
}

func EnvCloudAPISecret() string {
	return os.Getenv("CLOUDINARY_API_SECRET")
}

func EnvCloudUploadFolder() string {
	return os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
}

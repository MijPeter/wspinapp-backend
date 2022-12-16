package common

import (
	"log"
	"os"
	"strconv"
)

// TODO if this is to stay then probably use map or enum and not plain strings

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

func EnvDBUser() string {
	return os.Getenv("POSTGRES_USER")
}

func EnvDBPassword() string {
	return os.Getenv("POSTGRES_PASSWORD")
}

func EnvDBName() string {
	return os.Getenv("POSTGRES_DB")
}

func EnvDBHost() string {
	return os.Getenv("POSTGRES_HOST")
}

func EnvDBPort() int {
	ret, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))

	if err != nil {
		log.Fatalln("Couldn't load port number for db from environment variable.")
	}

	return ret
}

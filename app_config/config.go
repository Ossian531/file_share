package app_config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port int
	S3_endpoint string
	Bucket_name string
	Bucket_region string
	UsePathStyle bool
}

var AppConfig *Config

func LoadEnv() {

    godotenv.Load()

	bucket_name := os.Getenv("BUCKET_NAME")

	s3_endpoint := os.Getenv("S3_ENDPOINT")

	bucket_region := os.Getenv("BUCKET_REGION")

	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		port = 80
	}

	// UsePathStyle is required for MinIO, but should be false for AWS S3
	usePathStyle := os.Getenv("USE_PATH_STYLE") == "true"

	AppConfig = &Config{
		Port: port,
		S3_endpoint: s3_endpoint,
		Bucket_name: bucket_name,
		Bucket_region: bucket_region,
		UsePathStyle: usePathStyle,
	}
}


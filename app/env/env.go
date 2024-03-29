package env

import (
	"os"
)

func IsProd() bool {
	return os.Getenv("ENV") == "prod"
}

func IsDev() bool {
	return os.Getenv("ENV") == "dev"
}

func IsTest() bool {
	return os.Getenv("ENV") == "test"
}

func GetDatabaseUrl() string {
	return os.Getenv("DATABASE_URL")
}

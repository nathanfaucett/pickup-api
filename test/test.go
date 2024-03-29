package test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/aicacia/pickup/app"
	"github.com/aicacia/pickup/app/config"
	"github.com/aicacia/pickup/app/env"
	"github.com/aicacia/pickup/app/repository"
	"github.com/aicacia/pickup/app/router"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func ResponseBody[T any](t *testing.T, resp *http.Response) T {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %v", err)
	}
	var value T
	err = json.Unmarshal(body, &value)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %v", err)
	}
	return value
}

func setupTest() *fiber.App {
	godotenv.Load("../.env", "../.env.test")
	databaseUrl, err := url.Parse(env.GetDatabaseUrl())
	if err != nil {
		log.Fatalf("error parsing database url: %v\n", err)
	}
	databaseUrl.Path = fmt.Sprintf("/pickup-test-%s", uuid.New().String())
	os.Setenv("DATABASE_URL", databaseUrl.String())
	setupDatabase()

	err = repository.InitDB()
	if err != nil {
		log.Fatalf("error initializing database: %v\n", err)
	}
	err = config.InitConfig()
	if err != nil {
		log.Fatalf("error initializing config: %v\n", err)
	}
	app.Version.Version = "0.0.0-test"

	fiberApp := fiber.New(fiber.Config{
		Prefork:       false,
		Network:       "",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "",
		AppName:       "",
	})
	router.InstallRouter(fiberApp)

	return fiberApp
}
func setupDatabase() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current working directory: %v\n", err)
	}
	cmd := exec.Command("sqlx", "database", "drop", "-y")
	cmd.Dir = path.Join(cwd, "..")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Printf("error dropping database: %v\n", err)
	}
	cmd = exec.Command("sqlx", "database", "create")
	cmd.Dir = path.Join(cwd, "..")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("error creating database: %v\n", err)
	}
	cmd = exec.Command("sqlx", "migrate", "run")
	cmd.Dir = path.Join(cwd, "..")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("error migrating database: %v\n", err)
	}
}

func teardownTest() {
	err := config.CloseConfigListener()
	if err != nil {
		log.Fatalf("error closing config listener: %v\n", err)
	}
	err = repository.CloseDB()
	if err != nil {
		log.Fatalf("error closing database: %v\n", err)
	}
	teardownDatabase()
}

func teardownDatabase() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current working directory: %v\n", err)
	}
	cmd := exec.Command("sqlx", "database", "drop", "-y")
	cmd.Dir = path.Join(cwd, "..")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("error dropping database: %v\n", err)
	}
}

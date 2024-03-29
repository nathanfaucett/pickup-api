package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aicacia/pickup/app/model"
	"github.com/gofiber/fiber/v2"
)

var fiberApp *fiber.App

func TestMain(m *testing.M) {
	fiberApp = setupTest()
	exitVal := m.Run()
	teardownTest()
	os.Exit(exitVal)
}

func TestVersion(t *testing.T) {
	req := httptest.NewRequest("GET", "/version", nil)
	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status code is not OK: %d", resp.StatusCode)
	}
	version := ResponseBody[model.VersionST](t, resp)
	if version.Version != "0.0.0-test" {
		t.Fatalf("Version is not correct: %s", version.Version)
	}
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := fiberApp.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status code is not OK: %d", resp.StatusCode)
	}
	health := ResponseBody[model.HealthST](t, resp)
	if !health.IsHealthy() {
		t.Fatalf("Health is not correct: %v", health)
	}
}

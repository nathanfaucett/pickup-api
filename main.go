package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/aicacia/pickup/app"
	"github.com/aicacia/pickup/app/config"
	"github.com/aicacia/pickup/app/controller"
	"github.com/aicacia/pickup/app/repository"
	"github.com/aicacia/pickup/app/router"
	"github.com/aicacia/pickup/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

var (
	Version string = "0.1.0"
	Build   string = fmt.Sprint(time.Now().UnixMilli())
)

// @title Pickup API
// @description Pickup API
// @contact.name Nathan Faucett
// @contact.email nathanfaucett@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	defer func() {
		rec := recover()
		if rec != nil {
			log.Fatalf("application panic: %v\n", rec)
		}
	}()
	godotenv.Load(".env", ".env.dev", ".env.prod")
	err := repository.InitDB()
	if err != nil {
		log.Fatalf("error initializing database: %v\n", err)
	}
	err = config.InitConfig()
	if err != nil {
		log.Fatalf("error initializing config: %v\n", err)
	}
	err = controller.InitOAuth2()
	if err != nil {
		log.Fatalf("error initializing oauth2: %v\n", err)
	}

	app.Version.Version = Version
	app.Version.Build = Build

	docs.SwaggerInfo.Version = Version
	uri, err := url.Parse(config.Get().URI)
	if err != nil {
		log.Fatalf("error parsing URI: %v\n", err)
	}
	docs.SwaggerInfo.Host = uri.Host

	logWriter := os.Stdout
	log.SetOutput(logWriter)
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.LUTC)

	// https://docs.gofiber.io/api/fiber#config
	fiberApp := fiber.New(fiber.Config{
		Prefork:       false,
		Network:       "",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "",
		AppName:       "",
	})
	fiberApp.Use(fiberRecover.New())
	fiberApp.Use(logger.New(logger.Config{
		Output:     logWriter,
		TimeZone:   "UTC",
		TimeFormat: "2006/01/02 15:04:05",
		Format:     "${time} ${status} - ${ip} ${latency} ${method} ${path}\n",
	}))
	router.InstallRouter(fiberApp)
	if config.Get().Dashboard.Enabled {
		fiberApp.Use("/dashboard", monitor.New())
	}

	addr := fmt.Sprintf("%s:%d", config.Get().Host, config.Get().Port)
	log.Printf("Listening on %s\n", addr)

	log.Fatal(fiberApp.Listen(addr))
}

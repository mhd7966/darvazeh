package main

import (
	sentryfiber "github.com/aldy505/sentry-fiber"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/darvazeh/config"
	"github.com/mhd7966/darvazeh/connections"
	_ "github.com/mhd7966/darvazeh/docs"
	"github.com/mhd7966/darvazeh/log"
	"github.com/mhd7966/darvazeh/middleware"
	"github.com/mhd7966/darvazeh/routes"
)

//// @host localhost:3000 -> if set when you have domain you should set domain and then you want to test it localy -> change the host -> nazarim behtare

// @title Darvazeh API
// @version 1.0
// @description I have no specific description
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @termsOfService http://swagger.io/terms/
// @BasePath /v0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	config.SetConfig()
	log.LogInit()
	connections.ConnectMySQLDB()

	defer connections.MYSQL.Close()

	app := fiber.New(fiber.Config{
		Prefork: true,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			hub := sentryfiber.GetHubFromContext(c)
			if hub == nil {
				return fiber.DefaultErrorHandler(c, e)
			}
			hub.CaptureException(e)
			return nil
		},
	})

	middleware.InitApiSentry(app)
	app.Get("/docs/*", swagger.Handler)
	log.Log.Info("Swagger handler route created :)")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hi Stupid:)")
	})

	routes.MainRouter(*app)

	app.Listen(":3000")

}

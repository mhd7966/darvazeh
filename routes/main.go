package routes

import (
	"github.com/abr-ooo/go-pkgs"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func MainRouter(app fiber.App) {

	router := app.Group("/v0", gopkgs.Auth)

	DomainRouter(router)
	DomainsRouter(router)

	log.Info("All routes created :)")
}

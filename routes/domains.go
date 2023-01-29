package routes

import (
	"github.com/abr-ooo/darvazeh/controllers"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func DomainsRouter(app fiber.Router) {

	domains := app.Group("/domains")

	domains.Get("", controllers.GetAllDomains)
	domains.Post("", controllers.NewDomain)
	domains.Delete("/:domain_name", controllers.DeleteDomain)

	log.Info("Domains routes created :)")
}

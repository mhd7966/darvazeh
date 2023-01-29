package routes

import (
	"github.com/abr-ooo/darvazeh/controllers"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func DomainRouter(app fiber.Router) {

	domain := app.Group("/domain")

	domain.Get("/:domain_name", controllers.GetDomainRecords)
	domain.Post("", controllers.NewRecord)
	domain.Put("/:record_id", controllers.UpdateRecord)
	domain.Delete("/:record_id", controllers.DeleteRecord)
	domain.Get("/checkns/:domain_name", controllers.CheckNS)

	log.Info("Domain routes created :)")

	
}

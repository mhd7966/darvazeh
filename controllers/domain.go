package controllers

import (
	"errors"
	"strconv"

	gopkgs "github.com/abr-ooo/go-pkgs"
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/darvazeh/log"
	"github.com/mhd7966/darvazeh/models"
	"github.com/mhd7966/darvazeh/repositories"
	"github.com/mhd7966/darvazeh/services/pdns"
	"github.com/sirupsen/logrus"
)

// GetAllDomains godoc
// @Summary get all domains
// @Description return all domains of a userID
// @ID get_all_domains_by_userID
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /domains [get]
func GetAllDomains(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed"
		log.Log.WithFields(logrus.Fields{
			"user_id":  userID,
			"verify":   verifyUser,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("GetAllDomains. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("GetAllDomains. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// exist, err := repositories.Exist("domains", "account", userID)
	// if err != nil {
	// 	response.Message = "Check User Failed"
	// 	log.WithFields(logrus.Fields{
	// 		"user_id":  userID,
	// 		"response": response.Message,
	// 		"error":    err.Error(),
	// 	}).Error("GetAllDomains. Check exist user in domains DB failed!")
	// 	return c.Status(fiber.StatusBadRequest).JSON(response)
	// }

	// if !exist {
	// 	response.Message = "This user doesn't have any domain"
	// 	log.WithFields(logrus.Fields{
	// 		"user_id":  userID,
	// 		"response": response.Message,
	// 	}).Info("GetAllDomains. This user doesn't exist in DB for get domains!")
	// 	return c.Status(fiber.StatusBadRequest).JSON(response)
	// }

	domains, err := repositories.GetDomains(userID)

	if err != nil {
		response.Message = "Get Domains Failed"
		log.Log.WithFields(logrus.Fields{
			"user_id":  userID,
			"domains":  domains,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("GetAllDomains. Get domains from DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = domains
	log.Log.WithFields(logrus.Fields{
		"user_id":  userID,
		"domains":  domains,
		"response": response.Message,
	}).Info("GetAllDomains. Get all domains succeed :)")

	return c.Status(fiber.StatusOK).JSON(response)
}

// NewDomain godoc
// @Summary new domain
// @Description new domain
// @ID new_domain
// @Param domainModel body models.Domain true "Domain info ->[account = user_id]"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /domains [post]
func NewDomain(c *fiber.Ctx) error {

	var response models.Response
	if true {
		return errors.New("wwwwww")
	}
	response.Status = "error"

	domainModel := new(models.Domain)
	err := c.BodyParser(domainModel)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewDomain. Parse body to domainModel failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	domainModel.Account = strconv.Itoa(int(gopkgs.UID(c)))

	validDomain, err := Validation("DOMAIN", domainModel.URL)
	if err != nil {
		response.Message = "Check Domain Format Failed"
		log.Log.WithFields(logrus.Fields{
			"domain":   domainModel.URL,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewDomain. Validation function of domain have error!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if !*validDomain {
		response.Message = "Domain Format Doesn't Valid"
		log.Log.WithFields(logrus.Fields{
			"domain":   domainModel.URL,
			"response": response.Message,
		}).Info("NewDomain. This domain format doesn't valid!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	exist, err := repositories.Exist("domains", "name", domainModel.URL)
	if err != nil {
		response.Message = "Check Domain Failed"
		log.Log.WithFields(logrus.Fields{
			"domain":   domainModel.URL,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewDomain. Check exist domain in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if exist {
		response.Message = "Duplicate URL"
		log.Log.WithFields(logrus.Fields{
			"domain":   domainModel.URL,
			"response": response.Message,
		}).Info("NewDomain. This domain is duplicate. we have one of this!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	domainModel.Name = domainModel.URL + "."
	err = pdns.NewDomain(domainModel)
	if err != nil {
		response.Message = "Create Domain Failed"
		log.Log.WithFields(logrus.Fields{
			"domain":   domainModel.URL,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewDomain. Request to PDNS for create domain failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"domain":   domainModel.URL,
		"response": response.Message,
	}).Info("NewDomain. Create domain succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteDomain godoc
// @Summary delete domain
// @Description delete all domain info
// @ID delete_domain
// @Param domain_name path string true "domain_name"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /domains/{domain_name} [delete]
func DeleteDomain(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("DeleteDomain. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("DeleteDomain. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	domain := c.Params("domain_name")

	exist, err := repositories.Exist("domains", "name", domain)
	if err != nil {
		response.Message = "Check Domain Failed"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("DeleteDomain. Check exist doamin in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !exist {
		response.Message = "This domain doesn't exist"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"response": response.Message,
		}).Info("DeleteDomain. This doamin doesn't exist in DB for delete!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyDomain, err := VerifyUserDomain(userID, domain)
	if err != nil {
		response.Message = "Check Access User To domain Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id":       userID,
			"del_domain_id": domain,
			"verify":        verifyDomain,
			"error":         err.Error(),
		}).Error("DeleteDomain. Verify access user to domain failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyDomain {
		response.Message = "Access User To domain Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id":       userID,
			"del_domain_id": domain,
			"verify":        verifyDomain,
		}).Info("DeleteDomain. User doesn't have access to domain!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = pdns.RemoveDomain(domain)
	if err != nil {
		response.Message = "Delete Domains Failed"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("DeleteDomain. Request to PDNS for delete domain failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"domain":   domain,
		"response": response.Message,
	}).Info("DeleteDomain. Delete domain succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

func VerifyUser(user_id int) (bool, string, error) {

	exist, err := repositories.Exist("domains", "account", strconv.Itoa(user_id))
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"user_id": user_id,
			"error":   err.Error(),
		}).Debug("Verify User. check exist user failed!")
		return false, "", err
	}
	log.Log.WithFields(logrus.Fields{
		"user_id": user_id,
		"exist":   exist,
	}).Debug("VerifyUser. success!")
	return exist, strconv.Itoa(user_id), nil
}

func VerifyUserDomain(user_id string, domain_name string) (bool, error) {
	verify, err := repositories.VerifyUserDomain(user_id, domain_name)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"user_id":     user_id,
			"domain_name": domain_name,
			"error":       err.Error(),
		}).Debug("VerifyUserDomain. check access user to domain failed!")
		return false, err
	}

	log.Log.WithFields(logrus.Fields{
		"user_id":     user_id,
		"domain_name": domain_name,
		"verify":      verify,
	}).Debug("VerifyUserDomain. success!")
	return verify, nil
}

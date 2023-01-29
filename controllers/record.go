package controllers

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/abr-ooo/darvazeh/config"
	"github.com/abr-ooo/darvazeh/log"
	"github.com/abr-ooo/darvazeh/models"
	"github.com/abr-ooo/darvazeh/repositories"
	"github.com/abr-ooo/darvazeh/services/pdns"
	gopkgs "github.com/abr-ooo/go-pkgs"
	"github.com/gofiber/fiber/v2"
	"github.com/lixiangzhong/dnsutil"
	"github.com/sirupsen/logrus"
)

// GetDomainInfo godoc
// @Summary get domain info
// @Description return domain info
// @ID get_info_of_domain
// @Param domain_name path string true "domain_name"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /domain/{domain_name} [get]
func GetDomainRecords(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("GetDomainRecords. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("GetDomainRecords. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	domain := c.Params("domain_name")

	verifyDomain, err := repositories.Exist("domains", "name", domain)
	if err != nil {
		response.Message = "Check Domains Failed!"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"verify":   verifyDomain,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("GetDomainRecords. Check exist domain in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyDomain {
		response.Message = "This domain doesn't exist!"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"verify":   verifyDomain,
			"response": response.Message,
		}).Info("GetDomainRecords. This domain doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyUserDomain, err := VerifyUserDomain(userID, domain)
	if err != nil {
		response.Message = "Check Access User To domain Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id":       userID,
			"del_domain_id": domain,
			"verify":        verifyUserDomain,
			"error":         err.Error(),
		}).Error("GetDomainRecords. Verify access user to domain failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUserDomain {
		response.Message = "Access User To domain Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id":       userID,
			"del_domain_id": domain,
			"verify":        verifyUserDomain,
		}).Info("GetDomainRecords. User doesn't have access to domain!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	records, err := repositories.GetDomainInfo(domain)
	if err != nil {
		response.Message = "Get Domains Failed!"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"records":  records,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("GetDomainRecords. Get domain's records from DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = records
	log.Log.WithFields(logrus.Fields{
		"domain":   domain,
		"records":  records,
		"response": response.Message,
	}).Info("GetDomainRecords. Get records succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)

	// zone, err := pdns.GetDomainInfo(domain)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON("Get Domains Failed")
	// }

	// return c.Status(fiber.StatusOK).JSON(zone)
}

// NewRecord godoc
// @Summary new record
// @Description new record
// @ID new_record
// @Param recordBody body models.RecordBody true "Record info: *Just MX record must have priority*"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /domain [post]
func NewRecord(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	recordModel, domain, err := RecordProperty(c)
	if err != nil {
		response.Message = err.Error()
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewRecord. Set record property and checking content failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("NewRecord. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("NewRecord. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyDomain, err := repositories.Exist("domains", "name", domain)
	if err != nil {
		response.Message = "Check Domains Failed"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"verify":   verifyDomain,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewRecord. Check exist domain in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyDomain {
		response.Message = "This domain doesn't exist"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"verify":   verifyDomain,
			"response": response.Message,
		}).Info("NewRecord. This domain doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyUserDomain, err := VerifyUserDomain(userID, domain)
	if err != nil {
		response.Message = "Check Access User To domain Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id":       userID,
			"del_domain_id": domain,
			"verify":        verifyUserDomain,
			"error":         err.Error(),
		}).Error("NewRecord. Verify access user to domain failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUserDomain {
		response.Message = "Access User To domain Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id":       userID,
			"del_domain_id": domain,
			"verify":        verifyUserDomain,
		}).Info("NewRecord. User doesn't have access to domain!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	values, exist, err := repositories.PreviousRecords(recordModel.RRSet[0].Name, recordModel.RRSet[0].Type)
	if err != nil {
		response.Message = "Check Record Failed"
		log.Log.WithFields(logrus.Fields{
			"exist_previous_values": exist,
			"previous_values":       values,
			"response":              response.Message,
			"error":                 err.Error(),
		}).Error("NewRecord. Check exist and Get previous values of this name and type in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if exist {
		if recordModel.RRSet[0].Type == "CNAME" {
			response.Message = "Duplicate CNAME Record"
			log.Log.WithFields(logrus.Fields{
				"response": response.Message,
			}).Info("NewRecord. A domain can't have one more CNAME on a name!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		content := recordModel.RRSet[0].Records[0].Content
		SetPreviousValues(*values, recordModel)
		var record = new(models.Record)
		record.Content = content
		record.Disabled = false
		recordModel.RRSet[0].Records = append(recordModel.RRSet[0].Records, *record)

	}

	if recordModel.RRSet[0].Type == "TXT" {
		recordModel.RRSet[0].Records[len(recordModel.RRSet[0].Records)-1].Content = FormatValues(recordModel.RRSet[0].Records[len(recordModel.RRSet[0].Records)-1].Content, recordModel.RRSet[0].Type)

	} else {
		for key, record := range recordModel.RRSet[0].Records {
			recordModel.RRSet[0].Records[key].Content = FormatValues(record.Content, recordModel.RRSet[0].Type)
		}
	}

	err = pdns.ManageRecord(*recordModel)
	if err != nil {
		response.Message = err.Error()
		log.Log.WithFields(logrus.Fields{
			"record_model": recordModel,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("NewRecord. Request to PDNS for create new record failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"record_model": recordModel,
		"response":     response.Message,
	}).Info("NewRecord. Create new record succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

// UpdateRecord godoc
// @Summary update record
// @Description update_record
// @ID update_record
// @Param record_id path string true "record_id"
// @Param recordBody body models.RecordBody true "Record info"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /domain/{record_id} [put]
func UpdateRecord(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("UpdateRecord. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("UpdateRecord. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	recordID := c.Params("record_id")

	verifyRecord, err := repositories.Exist("records", "id", recordID)
	if err != nil {
		response.Message = "Check Record Failed!"
		log.Log.WithFields(logrus.Fields{
			"record_id": recordID,
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("UpdateRecord. Check exist record with record_id in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyRecord {
		response.Message = "This record does not exist"
		log.Log.WithFields(logrus.Fields{
			"record_id": recordID,
			"response":  response.Message,
		}).Info("UpdateRecord. This record doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	previousRecord, err := repositories.GetRecordContent(recordID)
	previousContent := previousRecord.Content
	if err != nil {
		response.Message = "Get Record Failed"
		log.Log.WithFields(logrus.Fields{
			"record_id":       recordID,
			"previous_record": previousRecord,
			"response":        response.Message,
			"error":           err.Error(),
		}).Error("UpdateRecord. Get previous record value from DB failed for update record!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyUserRecord, err := VerifyUserRecord(previousRecord.DomainID, userID, recordID)
	if err != nil {
		response.Message = "Check Access User To record Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id":   userID,
			"record_id": recordID,
			"verify":    verifyUserRecord,
			"error":     err.Error(),
		}).Error("UpdateRecord. Verify access user to record failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUserRecord {
		response.Message = "Access User To record Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id":   userID,
			"record_id": recordID,
			"verify":    verifyUserRecord,
		}).Info("UpdateRecord. User doesn't have access to record!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	recordModel, _, err := RecordProperty(c)
	if err != nil {
		response.Message = err.Error()
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("UpdateRecord. Set record property and checking content failed for update record!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	//ghatan record haye ghabli vojod dare (hadaghal 1 record ke hamon recordi hast ke mikhaym taghir bedim) -> be exist niazi nadarim
	values, _, err := repositories.PreviousRecords(recordModel.RRSet[0].Name, recordModel.RRSet[0].Type)
	if err != nil {
		response.Message = "Get Previous Values Failed"
		log.Log.WithFields(logrus.Fields{
			"record_model":    recordModel,
			"previous_values": values,
			"response":        response.Message,
			"error":           err.Error(),
		}).Error("UpdateRecord. Get previous values of this name and type in DB failed for update record!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	newContent := recordModel.RRSet[0].Records[0].Content
	SetPreviousValues(*values, recordModel)

	for key, record := range recordModel.RRSet[0].Records {
		if recordModel.RRSet[0].Records[key].Content == previousContent {
			var newRecord = new(models.Record)
			if recordModel.RRSet[0].Type == "TXT" {
				record.Content = FormatValues(newContent, recordModel.RRSet[0].Type)
			} else {
				record.Content = newContent
			}
			record.Disabled = false
			recordModel.RRSet[0].Records[key] = *newRecord
		}
		if recordModel.RRSet[0].Type == "TXT" {
			recordModel.RRSet[0].Records[key].Content = record.Content
		} else {
			recordModel.RRSet[0].Records[key].Content = FormatValues(record.Content, recordModel.RRSet[0].Type)
		}
	}

	err = pdns.ManageRecord(*recordModel)
	if err != nil {
		response.Message = "Update Record Failed"
		log.Log.WithFields(logrus.Fields{
			"record_model": recordModel,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("UpdateRecord. Request to PDNS for update record failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"record_model": recordModel,
		"response":     response.Message,
	}).Info("UpdateRecord. Update record succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteRecord godoc
// @Summary delete record
// @Description delete record info
// @ID delete_record
// @Param record_id path string true "record_id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /domain/{record_id} [delete]
func DeleteRecord(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("DeleteRecord. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("DeleteRecord. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	recordID := c.Params("record_id")
	verifyRecord, err := repositories.Exist("records", "id", recordID)
	if err != nil {
		response.Message = "Check Record Failed"
		log.Log.WithFields(logrus.Fields{
			"record_id": recordID,
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("DeleteRecord. Check exist record with record_id in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyRecord {
		response.Message = "This record does not exist"
		log.Log.WithFields(logrus.Fields{
			"record_id": recordID,
			"response":  response.Message,
		}).Info("DeleteRecord. This record doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	recordDB, err := repositories.GetRecordContent(recordID)
	if err != nil {
		response.Message = "Get Record Failed"
		log.Log.WithFields(logrus.Fields{
			"record_id":       recordID,
			"previous_record": recordDB,
			"response":        response.Message,
			"error":           err.Error(),
		}).Error("DeleteRecord. Get previous record value from DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyUserRecord, err := VerifyUserRecord(recordDB.DomainID, userID, recordID)
	if err != nil {
		response.Message = "Check Access User To record Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id":   userID,
			"record_id": recordID,
			"verify":    verifyUserRecord,
			"error":     err.Error(),
		}).Error("DeleteRecord. Verify access user to record failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUserRecord {
		response.Message = "Access User To record Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id":   userID,
			"record_id": recordID,
			"verify":    verifyUserRecord,
		}).Info("DeleteRecord. User doesn't have access to record!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	recordModel, err := RecordDBToModel(*recordDB)
	recordContent := recordModel.RRSet[0].Records[0].Content
	if err != nil {
		response.Message = "convert Record Failed"
		log.Log.WithFields(logrus.Fields{
			"record_id":    recordID,
			"record_db":    recordDB,
			"record_model": recordModel,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("DeleteRecord. Convert RecordDB model to RecordModel failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	//ghatan record haye ghabli vojod dare (hadaghal 1 record ke hamon recordi hast ke mikhaym taghir bedim) -> be exist niazi nadarim
	values, _, err := repositories.PreviousRecords(recordModel.RRSet[0].Name, recordModel.RRSet[0].Type)
	if err != nil {
		response.Message = "Get Previous Values Failed"
		log.Log.WithFields(logrus.Fields{
			"record_model":    recordModel,
			"previous_values": values,
			"response":        response.Message,
			"error":           err.Error(),
		}).Error("DeleteRecord. Get previous values of this name and type in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	SetPreviousValues(*values, recordModel)

	var index int
	for key, record := range recordModel.RRSet[0].Records {
		if recordModel.RRSet[0].Records[key].Content == recordContent {
			index = key
		}
		if recordModel.RRSet[0].Type != "TXT" {
			recordModel.RRSet[0].Records[key].Content = FormatValues(record.Content, recordModel.RRSet[0].Type)
		}
	}
	recordModel.RRSet[0].Records = append(recordModel.RRSet[0].Records[:index], recordModel.RRSet[0].Records[index+1:]...)

	err = pdns.ManageRecord(*recordModel)
	if err != nil {
		response.Message = "Delete Record Failed"
		log.Log.WithFields(logrus.Fields{
			"record_model": recordModel,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("DeleteRecord. Delete record with PDNS failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"record_model": recordModel,
		"response":     response.Message,
	}).Info("DeleteRecord. Delete record succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// CheckNS godoc
// @Summary check NS
// @Description check NS
// @ID check NS
// @Param domain_name path string true "domain_name"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /domain/checkns/{domain_name} [get]
func CheckNS(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"

	verifyUser, userID, err := VerifyUser(int(gopkgs.UID(c)))
	if err != nil {
		response.Message = "Check User Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
			"error":   err.Error(),
		}).Error("CheckNS. Verify user failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUser {
		response.Message = "User Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"verify":  verifyUser,
		}).Info("CheckNS. User doesn't verify!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	domain := c.Params("domain_name")

	verifyDomain, err := repositories.Exist("domains", "name", domain)
	if err != nil {
		response.Message = "Check Domains Failed"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"verify":   verifyDomain,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("CheckNS. Check exist domain in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyDomain {
		response.Message = "This domain doesn't exist"
		log.Log.WithFields(logrus.Fields{
			"domain":   domain,
			"verify":   verifyDomain,
			"response": response.Message,
		}).Info("CheckNS. This domain doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	verifyUserDomain, err := VerifyUserDomain(userID, domain)
	if err != nil {
		response.Message = "Check Access User To domain Failed!"
		log.Log.WithFields(logrus.Fields{
			"user_id":       userID,
			"del_domain_id": domain,
			"verify":        verifyUserDomain,
			"error":         err.Error(),
		}).Error("CheckNS. Verify access user to domain failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !verifyUserDomain {
		response.Message = "Access User To domain Doesn't Verify!"
		log.Log.WithFields(logrus.Fields{
			"user_id":       userID,
			"del_domain_id": domain,
			"verify":        verifyUserDomain,
		}).Info("CheckNS. User doesn't have access to domain!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	config := config.Cfg.NS
	NS := []string{strings.ToLower(config.NS1), strings.ToLower(config.NS2)}

	ns1, ns2, err := GetNS(domain)
	if err != nil {
		response.Message = "Get NSs from dig Failed, probebly there is no NS record for this domain!"
		response.Message = "Get NSs from dig Failed, probebly there is no NS record for this domain!"
		log.Log.Info("CheckNS. Get NSs form dig Failed with this error : \n", err)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	*ns1 = strings.ToLower(strings.TrimRight(*ns1, "."))
	*ns2 = strings.ToLower(strings.TrimRight(*ns2, "."))

	ns1Exist := Contains(NS, *ns1)
	ns2Exist := Contains(NS, *ns2)

	var nsResponse models.NSResponse
	nsResponse.CorrectNS1 = NS[0]
	nsResponse.CorrectNS2 = NS[1]
	nsResponse.CurrentNS1 = *ns1
	nsResponse.CurrentNS2 = *ns2
	response.Data = nsResponse

	if !ns1Exist && !ns2Exist {
		response.Message = "NS1, NS2 doesn't match"
		log.Log.WithFields(logrus.Fields{
			"ns1_exist": ns1Exist,
			"ns2_exist": ns2Exist,
			"response":  response.Message,
		}).Info("CheckNS. NS1, NS2 doesn't match!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !ns1Exist {
		response.Message = "NS1 doesn't match"
		log.Log.WithFields(logrus.Fields{
			"ns1_exist": ns1Exist,
			"ns1":       *ns1,
			"response":  response.Message,
		}).Info("CheckNS. NS1 doesn't match!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !ns2Exist {
		response.Message = "NS2 doesn't match"
		log.Log.WithFields(logrus.Fields{
			"ns2_exist": ns2Exist,
			"ns2":       *ns2,
			"response":  response.Message,
		}).Info("CheckNS. NS2 doesn't match!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = pdns.ManageRecord(*BuildNSRecord(domain, *ns1, *ns2))
	if err != nil {
		response.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"response": response.Message,
	}).Info("CheckNS. Check NS succeed :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

func GetNS(domain string) (*string, *string, error) {

	var dig dnsutil.Dig

	a, err := dig.NS(domain)
	if err != nil {

		log.Log.WithFields(logrus.Fields{
			"domain": domain,
			"error":  err.Error(),
		}).Debug("Aux. Get NSs with dig have error!")
		return nil, nil, err
	}
	if len(a) == 0 {
		return nil, nil, errors.New("there is no NS for this domain :(")
	}

	log.Log.WithFields(logrus.Fields{
		"NS1": a[0].Ns,
		"NS2": a[1].Ns,
	}).Debug("Aux.GetNS finish :))")
	return &a[0].Ns, &a[1].Ns, nil
}

func BuildNSRecord(domain, ns1, ns2 string) *models.RecordModel {
	recordModel := new(models.RecordModel)

	record := []models.Record{
		{Content: ns1 + ".", Disabled: false},
		{Content: ns2 + ".", Disabled: false}}

	rrset := []models.RRSet{{Name: domain + ".",
		Type:       "NS",
		TTL:        3600,
		ChangeType: "REPLACE",
		Records:    record},
	}
	recordModel.Name = domain + "."
	recordModel.RRSet = rrset

	log.Log.WithFields(logrus.Fields{
		"record_model": recordModel,
	}).Debug("Aux.BuildNSRecord finish :))")
	return recordModel

}

func VerifyUserRecord(domain_id int, userID string, recordID string) (bool, error) {
	verify, err := repositories.VerifyUserRecord(strconv.Itoa(domain_id), userID, recordID)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"user_id":   userID,
			"reocrd_id": recordID,
			"error":     err.Error(),
		}).Debug("Aux.VerifyUserRecord. check access user to record failed!")
		return false, err
	}
	log.Log.WithFields(logrus.Fields{
		"user_id":   userID,
		"reocrd_id": recordID,
		"verify":    verify,
	}).Debug("Aux.VerifyUserRecord. success!")
	return verify, nil
}

var validRecordType = []string{"A", "CNAME", "TXT", "MX", "NS"}

func RecordProperty(c *fiber.Ctx) (*models.RecordModel, string, error) {
	recordBody := new(models.RecordBody)
	err := c.BodyParser(recordBody)
	if err != nil {

		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("Aux. Parse body to RecordBody model failed!")
		return nil, "", errors.New("parse body failed")
	}

	if !Contains(validRecordType, recordBody.Type) {
		log.Log.WithFields(logrus.Fields{
			"record_type": recordBody.Type,
		}).Debug("Aux.Record type doesn't valid!")
		return nil, "", errors.New("record type doesn't valid")
	}

	valid, err := ValidationRecordValues(recordBody.Value, recordBody.Type)
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"record_type":  recordBody.Type,
			"record_value": recordBody.Value,
			"error":        err.Error(),
		}).Debug("Aux.Validation record value failed!")
		return nil, "", err
	}

	if !*valid {
		log.Log.WithFields(logrus.Fields{
			"record_type":  recordBody.Type,
			"record_value": recordBody.Value,
		}).Debug("Aux.Record value isn't valid!")
		return nil, "", errors.New("value doesn't valid")
	}

	recordModel := new(models.RecordModel)

	record := []models.Record{{Content: recordBody.Value,
		Disabled: false}}

	//set priority to mx record
	if recordBody.Type == "MX" {
		if recordBody.Priority != 7000 {
			record[0].Content = strconv.Itoa(recordBody.Priority) + " " + recordBody.Value
		} else {

			log.Log.Debug("Aux.MX record doesn't have priority!")
			return nil, "", errors.New("you must set priority for MX record")
		}
	}

	rrset := []models.RRSet{{Name: recordBody.Name + ".",
		Type:       recordBody.Type,
		TTL:        recordBody.TTL,
		ChangeType: "REPLACE",
		Records:    record},
	}

	checkSub := strings.Contains(recordBody.Name, recordBody.Domain)
	if checkSub {
		recordModel.Name = recordBody.Domain + "."
		recordModel.RRSet = rrset
	} else {
		log.Log.Debug("Aux.Name of record is out of domain")
		return nil, "", errors.New("name is out of domain")
	}

	log.Log.Debug("Aux.RecordProperty finish :))")
	return recordModel, recordBody.Domain, nil
}

func SetPreviousValues(values []string, recordModel *models.RecordModel) {

	var records []models.Record
	//add previous records in db
	for _, value := range values {

		var record = new(models.Record)
		record.Content = value
		record.Disabled = false
		records = append(records, *record)
	}

	recordModel.RRSet[0].Records = records

	log.Log.Debug("Aux.SetPreviousValues finish :))")
}

func ValidationRecordValues(recordValue string, recordType string) (*bool, error) {

	var valid *bool
	var err error

	switch recordType {
	case "TXT":
		valid, err = Validation("TXT", recordValue)
	case "MX", "CNAME":
		valid, err = Validation("DOMAIN", recordValue)

	case "A":
		valid, err = Validation("A", recordValue)

	default:
		v := true
		return &v, nil
	}
	log.Log.WithFields(logrus.Fields{
		"record_type":  recordType,
		"record_value": recordValue,
		"valid":        valid,
	}).Debug("Aux.ValidationRecordValues finish :))")
	return valid, err
}

func FormatValues(recordValue string, recordType string) string {
	switch recordType {
	case "TXT":
		return "\"" + recordValue + "\""
	case "MX", "CNAME", "NS":
		return recordValue + "."
	default:
		return recordValue
	}
}

func Validation(valueType string, value string) (*bool, error) {

	var r *regexp.Regexp
	var err error

	switch valueType {
	case "TXT":
		r, err = regexp.Compile(`([a-z0-9])+$`)
	case "DOMAIN":
		r, err = regexp.Compile(`^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`)
	case "A":
		r, err = regexp.Compile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	}

	match := r.MatchString(value)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"record_type":  valueType,
			"record_value": value,
			"error":        err.Error(),
		}).Debug("Aux.Validation have error!")
		return nil, err
	}

	log.Log.WithFields(logrus.Fields{
		"record_type":  valueType,
		"record_value": value,
		"match":        match,
	}).Debug("Aux.Validation finish :))")
	return &match, nil

}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func RecordDBToModel(recordDB models.RecordDB) (*models.RecordModel, error) {

	recordModel := new(models.RecordModel)

	record := []models.Record{{Content: recordDB.Content,
		Disabled: false}}

	//set priority to mx record
	if recordDB.Type == "MX" {
		if recordDB.Priority != 7000 {
			record[0].Content = strconv.Itoa(recordDB.Priority) + " " + recordDB.Content
		} else {
			log.Log.Debug("Aux.RecordDBToModel, MX record doesn't have priority")
			return nil, errors.New("record must have priority for MX record")
		}
	}

	rrset := []models.RRSet{{Name: recordDB.Name + ".",
		Type:       recordDB.Type,
		TTL:        recordDB.TTL,
		ChangeType: "REPLACE",
		Records:    record},
	}

	domain, err := repositories.GetDOmain(strconv.Itoa(recordDB.DomainID))
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"domain_id": recordDB.DomainID,
			"error":     err.Error(),
		}).Debug("Aux.RecordDBToModel, get domain from DB with domain_id have error!")
		return nil, err
	}

	recordModel.Name = *domain + "."
	recordModel.RRSet = rrset

	log.Log.WithFields(logrus.Fields{
		"reocrd_model": recordModel,
	}).Debug("Aux.RecordDBToModel finish :))")
	return recordModel, nil

}

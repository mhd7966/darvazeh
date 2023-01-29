package pdns

import (
	"bytes"
	"errors"
	"strings"

	"github.com/imroc/req"
	"github.com/mhd7966/darvazeh/config"
	"github.com/mhd7966/darvazeh/models"
	log "github.com/sirupsen/logrus"
)

func GetDomainInfo(domainName string) (*models.Zone, error) {

	config := config.Cfg.PowerDNS

	header := req.Header{
		"Accept":    "application/json",
		"X-API-Key": config.XAPIKey,
	}

	url := "http://" + config.Host + "/api/v1/servers/" + config.ServerID + "/zones/" + domainName + "."
	r, err := req.Get(url, header)
	if err != nil {
		log.WithFields(log.Fields{
			"url":    url,
			"result": r,
			"error":  err.Error(),
		}).Debug("PDNS. GetDomainInfo have error!")
		return nil, err
	}

	var resp models.Zone
	r.ToJSON(&resp)

	log.WithFields(log.Fields{
		"response": resp,
	}).Debug("PDNS. GetDomainInfo finish :))")
	return &resp, nil
}

func RemoveDomain(domainName string) error {

	config := config.Cfg.PowerDNS

	header := req.Header{
		"Accept":    "application/json",
		"X-API-Key": config.XAPIKey,
	}

	url := "http://" + config.Host + "/api/v1/servers/" + config.ServerID + "/zones/" + domainName + "."
	r, err := req.Delete(url, header)

	if err != nil {
		log.WithFields(log.Fields{
			"url":    url,
			"result": r,
			"error":  err.Error(),
		}).Debug("PDNS. RemoveDomain have error!")
		return err
	}

	if r.Response().StatusCode == 204 {
		log.WithFields(log.Fields{
			"response_code": 204,
		}).Debug("PDNS. GetDomainInfo finish :))")
		return nil
	}
	return err
}

func NewDomain(domain *models.Domain) error {

	config := config.Cfg.PowerDNS

	header := req.Header{
		"Accept":    "application/json",
		"X-API-Key": config.XAPIKey,
	}

	url := "http://" + config.Host + "/api/v1/servers/" + config.ServerID + "/zones"
	r, err := req.Post(url, header, req.BodyJSON(domain))

	if err != nil {
		log.WithFields(log.Fields{
			"url":    url,
			"result": r,
			"error":  err.Error(),
		}).Debug("PDNS. NewDomain have error!")
		return err
	}

	if r.Response().StatusCode == 201 {
		log.WithFields(log.Fields{
			"response_code": 201,
		}).Debug("PDNS. NewDomain finish :))")
		return nil
	}

	return err

}

func ManageRecord(record models.RecordModel) error {

	config := config.Cfg.PowerDNS

	header := req.Header{
		"Accept":    "application/json",
		"X-API-Key": config.XAPIKey,
	}
	url := "http://" + config.Host + "/api/v1/servers/" + config.ServerID + "/zones/" + record.Name
	r, err := req.Patch(url, header, req.BodyJSON(&record))

	if err != nil {
		log.WithFields(log.Fields{
			"url":    url,
			"result": r,
			"error":  err.Error(),
		}).Debug("PDNS. ManageRecord have error!")
		return err
	}

	if r.Response().StatusCode == 204 {
		log.WithFields(log.Fields{
			"response_code": 204,
		}).Debug("PDNS. ManageRecord finish :))")
		return nil
	} else {

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Response().Body)
		errorStr := buf.String()

		if strings.Contains(errorStr, "Conflicts with pre-existing RRset") {
			log.Debug("PDNS. ManageRecord name of CNAME record Shouldn't be the same as the names of other records!")
			return errors.New("name of CNAME record Shouldn't be the same as the names of other records")

		} else if strings.Contains(errorStr, "Duplicate record in RRset") {
			log.Debug("PDNS. ManageRecord duplicate record!")
			return errors.New("duplicate record")

		} else if strings.Contains(errorStr, "Not Found") {
			log.Debug("PDNS. ManageRecord domain not found!")
			return errors.New("domain not found")
		}
		log.Debug("PDNS. ManageRecord", errorStr)
		return errors.New(errorStr)
	}
}

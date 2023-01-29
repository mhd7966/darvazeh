package repositories

import (
	"strconv"

	"github.com/abr-ooo/darvazeh/connections"
	"github.com/abr-ooo/darvazeh/models"
	log "github.com/sirupsen/logrus"
)

//Find previous records that have same type
func PreviousRecords(recordName string, recordType string) (*[]string, bool, error) {

	if last := len(recordName) - 1; last >= 0 && recordName[last] == '.' {
		recordName = recordName[:last]
	}

	query := "SELECT * FROM records WHERE name=? and type=?"
	rows, err := connections.MYSQL.Query(query, recordName, recordType)

	if err != nil {
		log.WithFields(log.Fields{
			"record_name": recordName,
			"record_type": recordType,
			"error":       err.Error(),
		}).Debug("Execution *PreviousRecords* query in DB have error!")
		return nil, false, err
	}
	record := new(models.RecordDB)
	var values []string

	for rows.Next() {
		err := rows.Scan(&record.ID, &record.DomainID, &record.Name, &record.Type, &record.Content, &record.TTL, &record.Priority, &record.Disabled, &record.OrderName, &record.Auth)
		if err != nil {
			log.WithFields(log.Fields{
				"record": record,
				"error":  err.Error(),
			}).Debug("Scan result of *PreviousRecords* query have error!")
			return nil, false, err
		}

		if record.Type == "MX" {
			record.Content = strconv.Itoa(record.Priority) + " " + record.Content
		}
		values = append(values, record.Content)
	}

	if err := rows.Err(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Debug("Iterate on *PreviousRecords* query result have error!")
		return nil, false, err
	}

	defer rows.Close()
	exist := len(values) > 0
	log.WithFields(log.Fields{
		"values": values,
		"exist":  exist,
	}).Debug("Repo.GetRecordContent Finish :))")
	if exist {
		return &values, exist, nil
	} else {
		return nil, exist, nil
	}
}

func GetDomainInfo(domain_name string) (*[]models.RecordBody, error) {

	query := "SELECT r.id, r.domain_id, r.name, r.type, r.content, r.ttl, r.prio, r.disabled, r.ordername, r.auth FROM records r INNER JOIN domains d ON r.domain_id=d.id WHERE d.name=? and r.type!='NS' and r.type!='SOA'"
	rows, err := connections.MYSQL.Query(query, domain_name)

	if err != nil {
		log.WithFields(log.Fields{
			"domain_name": domain_name,
			"error":       err.Error(),
		}).Debug("Execution *GetDomainInfo* query in DB have error!")
		return nil, err
	}
	var record models.RecordDB
	var records []models.RecordBody
	for rows.Next() {
		err := rows.Scan(&record.ID, &record.DomainID, &record.Name, &record.Type, &record.Content, &record.TTL, &record.Priority, &record.Disabled, &record.OrderName, &record.Auth)
		if err != nil {
			log.WithFields(log.Fields{
				"record": record,
				"error":  err.Error(),
			}).Debug("Scan result of *GetDomainInfo* query have error!")
			return nil, err
		}

		var recordBody models.RecordBody
		recordBody.ID = record.ID
		recordBody.Domain = domain_name
		recordBody.Name = record.Name
		recordBody.Type = record.Type
		recordBody.TTL = record.TTL
		recordBody.Value = record.Content
		recordBody.Priority = record.Priority

		records = append(records, recordBody)
	}
	if err := rows.Err(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Debug("Iterate on *GetDomainInfo* query result have error!")
		return nil, err
	}

	defer rows.Close()
	log.WithFields(log.Fields{
		"records": record,
	}).Debug("Repo.GetDomainInfo Finish :))")
	return &records, nil

}

func GetRecordContent(recordID string) (*models.RecordDB, error) {

	query := "SELECT * FROM records WHERE id=?"
	result, err := connections.MYSQL.Query(query, recordID)

	if err != nil {
		log.WithFields(log.Fields{
			"record_id": recordID,
			"error":     err.Error(),
		}).Debug("Execution *GetRecordContent* query in DB have error!")
		return nil, err
	}
	record := new(models.RecordDB)

	for result.Next() {
		err := result.Scan(&record.ID, &record.DomainID, &record.Name, &record.Type, &record.Content, &record.TTL, &record.Priority, &record.Disabled, &record.OrderName, &record.Auth)
		if err != nil {
			log.WithFields(log.Fields{
				"record": record,
				"error":  err.Error(),
			}).Debug("Scan result of *GetRecordContent* query have error!")
			return nil, err
		}

		if record.Type == "MX" {
			record.Content = strconv.Itoa(record.Priority) + " " + record.Content
		}
	}

	if err := result.Err(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Debug("Iterate on *GetRecordContent* query result have error!")
		return nil, err
	}

	defer result.Close()
	log.WithFields(log.Fields{
		"record": record,
	}).Debug("Repo.GetRecordContent Finish :))")
	return record, nil
}

func VerifyUserRecord(domain_id string, userID string, recordID string) (bool, error) {
	query := "SELECT * FROM records WHERE domain_id = (SELECT id FROM domains WHERE account = ? and id = ?) and id=?"
	result, err := connections.MYSQL.Query(query, userID, domain_id, recordID)

	if err != nil {
		log.WithFields(log.Fields{
			"user_id":   userID,
			"record_id": recordID,
			"domain_id": domain_id,
			"error":     err.Error(),
		}).Debug("Execution *VerifyUserRecord* query in DB have error!")
		return false, err
	}
	verify := result.Next()

	log.WithFields(log.Fields{
		"verify":    verify,
		"user_id":   userID,
		"record_id": recordID,
		"domain_id": domain_id,
	}).Debug("Repo.VerifyUserRecord Finish :))")
	return verify, nil
}

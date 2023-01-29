package repositories

import (
	"github.com/mhd7966/darvazeh/connections"
	log "github.com/sirupsen/logrus"
)

func GetDomains(userID string) (*[]string, error) {
	query := "SELECT name FROM domains WHERE account=?"
	rows, err := connections.MYSQL.Query(query, userID)

	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Debug("Execution *GetDomains* query in DB have error")
		return nil, err
	}

	var domain string
	var domains []string
	for rows.Next() {
		err := rows.Scan(&domain)
		if err != nil {
			log.WithFields(log.Fields{
				"domain": domain,
				"error":  err.Error(),
			}).Debug("Scan result of *GetDomains* query have error")
			return nil, err
		}

		domains = append(domains, domain)
	}
	if err := rows.Err(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Debug("Iterate on *GetDomains* query result have error")
		return nil, err
	}

	defer rows.Close()

	log.WithFields(log.Fields{
		"domains": domains,
	}).Debug("Repo.GetDomains Finish :))")
	return &domains, nil
}

func Exist(table string, field string, value string) (bool, error) {

	query := "SELECT * FROM " + table + " WHERE " + field + "=?"
	result, err := connections.MYSQL.Query(query, value)

	if err != nil {
		log.WithFields(log.Fields{
			"table": table,
			"field": field,
			"value": value,
			"error": err.Error(),
		}).Debug("Execution *Exist* query in DB have error")
		return false, err
	}
	exist := result.Next()

	defer result.Close()
	log.WithFields(log.Fields{
		"table": table,
		"field": field,
		"value": value,
		"exist": exist,
	}).Debug("Repo.Exist Finish :))")
	return exist, nil
}

func GetDOmain(domainID string) (*string, error) {

	query := "SELECT name FROM domains WHERE id=?"
	result, err := connections.MYSQL.Query(query, domainID)

	if err != nil {
		log.WithFields(log.Fields{
			"domain_id": domainID,
			"error":     err.Error(),
		}).Debug("Execution *GetDomain* query in DB have error")
		return nil, err
	}

	var domain string

	for result.Next() {
		err := result.Scan(&domain)
		if err != nil {
			log.WithFields(log.Fields{
				"domain": domain,
				"error":  err.Error(),
			}).Debug("Scan result of *GetDomain* query have error")
			return nil, err
		}
	}

	if err := result.Err(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Debug("Iterate on *GetDomain* query result have error")
		return nil, err
	}

	defer result.Close()
	log.WithFields(log.Fields{
		"domain": domain,
	}).Debug("Repo.GetDomain Finish :))")
	return &domain, nil
}

func VerifyUserDomain(user_id string, domain_name string) (bool, error) {
	query := "SELECT * FROM domains WHERE account=? and name=?"
	result, err := connections.MYSQL.Query(query, user_id, domain_name)

	if err != nil {
		log.WithFields(log.Fields{
			"user_id":     user_id,
			"domain_name": domain_name,
			"error":       err.Error(),
		}).Debug("Execution *VerifyUserDomain* query in DB have error!")
		return false, err
	}
	verify := result.Next()

	log.WithFields(log.Fields{
		"verify":      verify,
		"user_id":     user_id,
		"domain_name": domain_name,
	}).Debug("Repo.Verify User Domain Finish :))")
	return verify, nil
}

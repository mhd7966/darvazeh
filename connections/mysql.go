package connections

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mhd7966/darvazeh/config"
	"github.com/mhd7966/darvazeh/log"
	"github.com/sirupsen/logrus"
)

var MYSQL *sql.DB

func ConnectMySQLDB() error {

	config := config.Cfg.MySQL

	var err error

	MYSQL, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.User, config.Pass, config.Host, config.Port, config.Name))

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Connect to DB Failed !!!!!")
		return err
	}

	err = MYSQL.Ping()
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Ping DB Have Error !!!!!")
		return err
	}

	return nil

}

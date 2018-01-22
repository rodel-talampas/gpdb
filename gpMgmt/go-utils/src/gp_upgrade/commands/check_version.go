package commands

import (
	"fmt"

	_ "github.com/lib/pq"

	"io"

	"gp_upgrade/db"
	"os"

	"regexp"

	"gp_upgrade/utils"

	"github.com/cppforlife/go-semi-semantic/version"
	"github.com/jmoiron/sqlx"
)

type CheckVersionCommand struct {
	Master_host string `long:"master-host" required:"yes" description:"Domain name or IP of host"`
	Master_port int    `long:"master-port" required:"no" default:"15432" description:"Port for master database"`
}

const (
	MINIMUM_VERSION = "4.3.9.0"
)

func (cmd CheckVersionCommand) Execute([]string) error {
	dbConn := db.NewDBConn(cmd.Master_host, cmd.Master_port, "template1")
	return cmd.execute(dbConn, os.Stdout)
}

func (cmd CheckVersionCommand) execute(dbConnector db.DBConnector, outputWriter io.Writer) error {

	err := dbConnector.Connect()
	if err != nil {
		return utils.DatabaseConnectionError{Parent: err}
	}
	defer dbConnector.Close()

	var connection *sqlx.DB
	connection = dbConnector.GetConn()
	var row string
	err = connection.QueryRow("SELECT version()").Scan(&row)
	if err != nil {
		return err
	}

	re := regexp.MustCompile("Greenplum Database (.*) build")

	version_string := re.FindStringSubmatch(row)[1]
	version_object := version.MustNewVersionFromString(version_string)

	if version_object.IsGt(version.MustNewVersionFromString(MINIMUM_VERSION)) {
		fmt.Fprintf(outputWriter, "gp_upgrade: Version Compatibility Check [OK]\n")
	} else {
		fmt.Fprintf(outputWriter, "gp_upgrade: Version Compatibility Check [Failed]\n")
	}
	return err
}

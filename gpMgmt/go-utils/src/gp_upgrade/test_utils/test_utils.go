package test_utils

import (
	"fmt"
	"gp_upgrade/config"
	"gp_upgrade/db"
	"io/ioutil"
	"os"

	"github.com/greenplum-db/gpbackup/testutils"

	"path"

	"errors"

	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	TempHomeDir = "/tmp/gp_upgrade_test_temp_home_dir"

	SAMPLE_JSON = `[{
    "address": "briarwood",
    "content": 2,
    "datadir": "/Users/pivotal/workspace/gpdb/gpAux/gpdemo/datadirs/dbfast_mirror3/demoDataDir2",
    "dbid": 7,
    "hostname": "briarwood",
    "mode": "s",
    "port": 25437,
    "preferred_role": "m",
    "role": "m",
    "san_mounts": null,
    "status": "u"
  }]`
)

func Check(msg string, e error) {
	if e != nil {
		panic(fmt.Sprintf("%s: %s\n", msg, e.Error()))
	}
}

func ResetTempHomeDir() string {
	config_dir := path.Join(TempHomeDir, ".gp_upgrade")
	if _, err := os.Stat(config_dir); !os.IsNotExist(err) {
		err = os.Chmod(config_dir, 0700)
		Check("cannot change mod", err)
	}
	err := os.RemoveAll(TempHomeDir)
	Check("cannot remove temp home", err)
	save := os.Getenv("HOME")
	err = os.MkdirAll(TempHomeDir, 0700)
	Check("cannot create home temp dir", err)
	err = os.Setenv("HOME", TempHomeDir)
	Check("cannot set home dir", err)
	return save
}

func WriteSampleConfig() {
	err := os.MkdirAll(config.GetConfigDir(), 0700)
	Check("cannot create sample dir", err)
	err = ioutil.WriteFile(config.GetConfigFilePath(), []byte(SAMPLE_JSON), 0600)
	Check("cannot write sample config", err)
}

func CreateMockDB() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	mockdb := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		Fail("Could not create mock database connection")
	}
	return mockdb, mock
}

func CreateMockDBConn(masterHost string, masterPort int) (*db.GPDBConnector, sqlmock.Sqlmock) {
	mockdb, mock := CreateMockDB()
	gpdbConnector := &db.GPDBConnector{
		Driver: testutils.TestDriver{DBExists: true, DB: mockdb, DBName: "testdb"},
		DBName: "testdb",
		Host:   masterHost,
		Port:   masterPort,
	}
	connection := gpdbConnector.GetConn()
	if connection != nil && connection.Stats().OpenConnections > 0 {
		Fail("connection before connect is called")
	}
	return gpdbConnector, mock
}

type FakeDbConnection struct {
	Driver db.DBDriver
}

func (fakedbconn FakeDbConnection) Connect() error {
	return errors.New("Invalid DB Connection")
}
func (fakedbconn FakeDbConnection) Close() {
}

func (fakdbconn FakeDbConnection) GetConn() *sqlx.DB {
	return nil
}

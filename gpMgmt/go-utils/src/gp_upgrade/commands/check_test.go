package commands_test

import (
	"database/sql"
	"os"

	"gp_upgrade/test_utils"

	"io/ioutil"

	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const (
	sqlite3_database_path = "/tmp/gp_upgrade_test_sqlite.db"
)

var _ = Describe("check", func() {

	AfterEach(func() {
		err := os.RemoveAll(sqlite3_database_path)
		test_utils.Check("Cannot remove sqllite database file", err)
	})
	Describe("happy: the database is running, master_host is provided, and connection is successful", func() {
		It("writes a file to ~/.gp_upgrade/cluster_config.json with correct json", func() {
			path := os.Getenv("GOPATH") + "/src/gp_upgrade/commands/fixtures/segment_config.sql"
			setupSqlite3Database(getFileContents(path))

			session := runCommand("check", "--master_host", "localhost", "--database_type", "sqlite3", "--database_config_file", sqlite3_database_path)

			Eventually(session).Should(Exit(0))
			content, err := ioutil.ReadFile(jsonFilePath())
			Expect(err).NotTo(HaveOccurred())
			expectedJson, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/gp_upgrade/commands/fixtures/segment_config.json")
			Expect(err).NotTo(HaveOccurred())
			Expect(expectedJson).To(Equal(content))
			var json_structure []map[string]interface{}
			err = json.Unmarshal(content, &json_structure)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("error cases", func() {
		Describe("the database cannot be opened", func() {
			It("returns error", func() {
				session := runCommand("check", "--master_host", "localhost", "--database_type", "foo", "--database_config_file", "bar")

				Eventually(session).Should(Exit(1))
				Expect(string(session.Err.Contents())).To(ContainSubstring(`sql: unknown driver "foo" (forgotten import?)`))
			})
		})
		Describe("the database query fails", func() {
			It("returns error", func() {
				session := runCommand("check", "--master_host", "localhost", "--database_type", "sqlite3", "--database_config_file", sqlite3_database_path)

				Eventually(session).Should(Exit(1))
				Expect(string(session.Err.Contents())).To(ContainSubstring(`no such table: gp_segment_configuration`))
			})
		})
	})
})

func jsonFilePath() string {
	return os.Getenv("HOME") + "/.gp_upgrade/cluster_config.json"
}

func setupSqlite3Database(inputSql string) {
	// clean any prior db
	err := ioutil.WriteFile(sqlite3_database_path, []byte(""), 0644)
	test_utils.Check("cannot delete sqlite config", err)

	db, err := sql.Open("sqlite3", sqlite3_database_path)
	test_utils.Check("cannot open sqlite config", err)
	defer db.Close()

	_, err = db.Exec(inputSql)
	test_utils.Check("cannot run sqlite config", err)

	err = os.RemoveAll(jsonFilePath())
	test_utils.Check("cannot remove json file", err)
}

func getFileContents(path string) string {
	segment_fixture_sql, err := ioutil.ReadFile(path)
	test_utils.Check("cannot open fixture:", err)
	return string(segment_fixture_sql)
}

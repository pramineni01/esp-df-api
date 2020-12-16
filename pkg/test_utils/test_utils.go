// Role of this package is to automate mundane tasks involved in writing and running tests.

package test_utils

import (
	"bytes"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"os/exec"
)

var global_db *sqlx.DB

// password user_name option used by mysql command line tool
var user_name string = fmt.Sprintf("--user=%s", os.Getenv("DB_USER"))
var password = fmt.Sprintf("--password=%s", os.Getenv("DB_PASSWORD"))

func setupDB() {
	dsn := fmt.Sprintf("%s:%s@/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sqlx.Open("mysql", dsn)
	global_db = db

	if err != nil {
		fmt.Println("err", err)
	}
	createTestDatabases()
}

func createTestDatabases() {
	global_db.MustExec("CREATE DATABASE IF NOT EXISTS seed")
	global_db.MustExec("CREATE DATABASE IF NOT EXISTS test")
	loadSchema("seed", takeSchemaDump())
	loadSchema("test", takeSchemaDump())
	loadSeedData()
}

func loadSchema(db string, schema *bytes.Buffer) {
	cmd := exec.Command("mysql", user_name, password, db)
	cmd.Stdin = schema
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func takeSchemaDump() *bytes.Buffer {
	cmd := exec.Command("mysqldump", "--no-data", user_name, password, os.Getenv("DB_NAME"))

	buff := new(bytes.Buffer)
	cmd.Stdout = buff
	err := cmd.Start()

	if err != nil {
		panic(err)
	}
	cmd.Wait()
	return buff

}

func loadSeedData() {
	f, err := os.Open("../test_utils/seed.sql")
	if err != nil {
		panic(err)

	}
	cmd := exec.Command("mysql", user_name, password, "seed")
	cmd.Stdin = f
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func LoadTables(tables []string) {
	for _, table := range tables {
		query := fmt.Sprintf("INSERT test.%[1]s  SELECT * FROM seed.%[1]s", table)
		global_db.MustExec(query)
	}

}

func UnLoadTables(tables []string) {
	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM test.%s", table)
		global_db.MustExec(query)
	}
}

// This function is the entry point of this package.
// The function starts by creating two new databases(seed, test), using the schema of the provided Databases.
// The role of the seed database is to hold all the test data(instead of it being distributed in ton of different test files) using a file currently name `seed.sql`
// After creating both the database and loading the data from `seed.sql` to seed databases, this function returns the connection to test database.
// Against which the test are going to run, and based on their requirements they can load or unload(after using) the tables from `seed db` to `test db`.
func GetTestDBConn() *sqlx.DB {
	setupDB()
	dsn := fmt.Sprintf("%s:%s@/test", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	db, err := sqlx.Open("mysql", dsn)

	if err != nil {
		panic(fmt.Errorf("Error %v occurred while connecting to test database", err))
	}
	return db
}

func StartServer() {
	cmd := exec.Command("go", "run", os.Getenv("SERVER_FILE_PATH"))
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

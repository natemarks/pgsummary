package pg

import (
	"database/sql"
	"fmt"
	"github.com/natemarks/pgsummary/util"
	"github.com/rs/zerolog"
	"sort"
)

const (
	DatabaseType = "postgres"
)

type ConnData struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

func (c ConnData) connString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.DbName)
}

type DatabaseUser struct {
	Name       string `json:"name"`
	Attributes string `json:"attributes"`
}

type TableReport struct {
	RowCount int               `json:"rowCount"`
	Columns  map[string]string `json:"columns"`
}

type DatabaseReport struct {
	Tables     map[string]TableReport `json:"tables"`
	Extensions []string               `json:"extensions"`
}

type InstanceReport struct {
	HostName  string                    `json:"hostName"`
	Port      int                       `json:"port"`
	Databases map[string]DatabaseReport `json:"databases"`
	Users     []DatabaseUser            `json:"users"`
}

func ValidateCredentials(
	cd ConnData,
	log *zerolog.Logger) {
	log.Info().Msgf("Validating %s credentials", DatabaseType)
	// open database
	db, err := sql.Open(DatabaseType, cd.connString())
	util.CheckError(err, log)

	// close database
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	err = db.Ping()
	util.CheckError(err, log)
}

func GetDatabaseNames(cd ConnData, log *zerolog.Logger) ([]string, error) {
	/// reserved databases that throw errors when we try to interrogate them
	invalidDatabases := []string{"rdsadmin", "template0", "template1"}
	// open database

	// open database
	db, err := sql.Open(DatabaseType, cd.connString())
	util.CheckError(err, log)

	// close database
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	rows, err := db.Query("SELECT datname FROM pg_database")
	util.CheckError(err, log)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// An album slice to hold data from returned rows.
	var allDBs []string

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var thisDB string
		if err := rows.Scan(&thisDB); err != nil {
			break
		}
		//  Don't append one of the invalid databases
		if util.Contains(invalidDatabases, thisDB) {
			continue
		}
		allDBs = append(allDBs, thisDB)
	}
	sort.Strings(allDBs)
	return allDBs, err
}

func GetInstanceUsers(
	cd ConnData,
	log *zerolog.Logger) ([]DatabaseUser, error) {
	var err error
	var users []DatabaseUser

	// open database
	db, err := sql.Open(DatabaseType, cd.connString())
	util.CheckError(err, log)

	// close database
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	sqlStatement := fmt.Sprintf("SELECT usename AS role_name,\n CASE\n  WHEN usesuper AND usecreatedb THEN\n    CAST('superuser, create database' AS pg_catalog.text)\n  WHEN usesuper THEN\n    CAST('superuser' AS pg_catalog.text)\n  WHEN usecreatedb THEN\n    CAST('create database' AS pg_catalog.text)\n  ELSE\n    CAST('' AS pg_catalog.text)\n END role_attributes\nFROM pg_catalog.pg_user\nORDER BY role_name desc;")

	rows, err := db.Query(sqlStatement)
	util.CheckError(err, log)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var dbu DatabaseUser
		if err := rows.Scan(&dbu.Name, &dbu.Attributes); err != nil {
			break
		}
		users = append(users, dbu)
	}
	return users, err
}

func GetTableRowCount(
	cd ConnData,
	table string,
	log *zerolog.Logger) (int, error) {
	var err error
	var exactCount int

	// open database
	db, err := sql.Open(DatabaseType, cd.connString())
	util.CheckError(err, log)

	// close database
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	sqlStatement := fmt.Sprintf("SELECT count(*) AS exact_count FROM %s", table)
	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&exactCount); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		return exactCount, err
	default:
		log.Fatal().Err(err)
	}
	return exactCount, err
}

func GetTableColumnsDetails(
	cd ConnData,
	table string,
	log *zerolog.Logger) (map[string]string, error) {
	var err error
	columns := make(map[string]string)

	// open database
	db, err := sql.Open(DatabaseType, cd.connString())
	util.CheckError(err, log)

	// close database
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	sqlStatement := fmt.Sprintf("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = '%s';", table)

	rows, err := db.Query(sqlStatement)
	util.CheckError(err, log)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			break
		}
		columns[columnName] = dataType
	}
	return columns, err
}

func GetTableReport(
	cd ConnData,
	table string,
	log *zerolog.Logger) (TableReport, error) {
	var err error
	rowCount, err := GetTableRowCount(cd, table, log)
	util.CheckError(err, log)
	columnDetails, err := GetTableColumnsDetails(cd, table, log)
	util.CheckError(err, log)

	tbr := TableReport{
		RowCount: rowCount,
		Columns:  columnDetails,
	}
	return tbr, err
}

func GetTableNames(cdata ConnData, log *zerolog.Logger) ([]string, error) {

	// open database
	db, err := sql.Open(DatabaseType, cdata.connString())
	util.CheckError(err, log)

	// close database
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	err = db.Ping()
	util.CheckError(err, log)
	// rows, err := db.Query("SELECT tablename FROM pg_catalog.pg_tables")
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name")
	util.CheckError(err, log)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// An album slice to hold data from returned rows.
	var allTables []string

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var thisTable string
		if err := rows.Scan(&thisTable); err != nil {
			break
		}
		allTables = append(allTables, thisTable)
	}
	sort.Strings(allTables)
	return allTables, err
}

func GetDatabaseExtensions(
	cd ConnData,
	log *zerolog.Logger) ([]string, error) {

	// open database
	db, err := sql.Open(DatabaseType, cd.connString())
	util.CheckError(err, log)

	// close database
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	err = db.Ping()
	util.CheckError(err, log)
	// rows, err := db.Query("SELECT tablename FROM pg_catalog.pg_tables")
	rows, err := db.Query("SELECT extname FROM pg_extension")
	util.CheckError(err, log)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// An album slice to hold data from returned rows.
	var extensions []string

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var thisExt string
		if err := rows.Scan(&thisExt); err != nil {
			break
		}
		extensions = append(extensions, thisExt)
	}
	return extensions, err
}

func GetDatabaseReport(
	cd ConnData,
	log *zerolog.Logger) (DatabaseReport, error) {
	var err error
	tableMap := make(map[string]TableReport)
	tableNames, err := GetTableNames(cd, log)
	for _, v := range tableNames {
		tbr, err := GetTableReport(cd, v, log)
		util.CheckError(err, log)
		tableMap[v] = tbr
	}

	extensions, err := GetDatabaseExtensions(cd, log)
	dbr := DatabaseReport{tableMap,
		extensions}

	return dbr, err
}

func GetInstanceReport(
	cd ConnData,
	log *zerolog.Logger) (InstanceReport, error) {
	var err error
	databaseMap := make(map[string]DatabaseReport)

	databaseNames, err := GetDatabaseNames(cd, log)
	util.CheckError(err, log)

	for _, v := range databaseNames {
		cd.DbName = v
		dbr, err := GetDatabaseReport(cd, log)
		util.CheckError(err, log)
		databaseMap[cd.DbName] = dbr
		//dbReportList = append(dbReportList, dbr)
	}
	users, err := GetInstanceUsers(cd, log)
	util.CheckError(err, log)
	instanceReport := InstanceReport{
		HostName:  cd.Host,
		Port:      cd.Port,
		Databases: databaseMap,
		Users:     users,
	}
	return instanceReport, err
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// database connection parameters
const (
	DB_TYPE  = "mysql"
	USERNAME = "root"
	PASSWORD = "capa1993"
	DB_NAME  = "event_database"
	DB_PORT  = "3306"
)

func createTables() {
	// try to open database connection
	db, err = sql.Open(DB_TYPE, USERNAME+":"+PASSWORD+"@tcp(127.0.0.1:"+DB_PORT+")/"+DB_NAME)
	if err != nil {
		fmt.Println("a")
		fmt.Println(err.Error())
		fmt.Println("b")
		panic(err)
	}

	// make sure connection is available
	err = db.Ping() //check that you can establish a network connection and log in
	if err != nil {
		panic(err)
	} else {
		createUserTable(db)
		createAPIListTable(db)
		createResponseTimesTable(db)
	}
	defer db.Close()
}
func createUserTable(db *sql.DB) {
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS user (id int NOT NULL, PRIMARY KEY (id));")
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	stmt.Close()
}
func createResponseTimesTable(db *sql.DB) {
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS responsetimes (id int AUTO_INCREMENT, responsetime int NOT NULL, PRIMARY KEY (id));")
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	stmt.Close()
}
func createAPIListTable(db *sql.DB) {
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS apilist (apiid int NOT NULL,  PRIMARY KEY (apiid));")
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	//assume we have 10 different api keys pre-defined
	var apikeys []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	// check if there is any row that exists in apilist table

	var apiCount = checkCountInApiList()
	checkErr(err)
	//if apilist table is empty, that means this is the first time to create these tables
	//So fill the table and also create all respective tables for each apikey to store events
	if apiCount < 1 {
		for i := 0; i < len(apikeys); i++ {
			stmt, err = db.Prepare("INSERT INTO apilist(apiid) VALUES(?)")
			checkErr(err)
			_, err = stmt.Exec(apikeys[i])
			//APİ_KEYS' definition storage is important, abort in case of a failure
			if err != nil {
				panic(err)
			}
			stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS " + "events_for_api_" + apikeys[i] +
				" (eventid int NOT NULL AUTO_INCREMENT ,  userID int NOT NULL,  ts TIMESTAMP," +
				" PRIMARY KEY (eventid), FOREIGN KEY(userID) REFERENCES user(id) ON DELETE CASCADE );")

			checkErr(err)
			_, err = stmt.Exec()
			checkErr(err)
		}
	}
	stmt.Close()
}
func checkCountInApiList() (count int) {
	rows, err = db.Query("SELECT COUNT(*) as count FROM apilist WHERE apiid = 1") //assume apiid 1 must exist
	defer rows.Close()
	//return number of rows in rows list taken in parameter
	for rows.Next() {
		err = rows.Scan(&count)
		checkErr(err)
	}
	return count
}

//just prints error info
func checkErr(err error) {
	if err != nil {
		log.Println("Error: ", err)
	}
}

//checks if the specified API_KEY exists in the database
func doesApiKeyExist(API_KEY string) bool {
	rows, err = db.Query("SELECT apiid FROM apilist WHERE apiid=?", API_KEY)
	defer rows.Close()
	checkErr(err)
	// if the user id in the processed event definition does not exist in database
	// create a new user in the users table
	//defer rows.Close()
	return rows.Next()
}
func doesUserExist(USER_ID string) bool {
	rows, err = db.Query("SELECT id FROM user WHERE id=?", USER_ID)
	defer rows.Close()
	checkErr(err)
	return rows.Next()
}
func addUser(USER_ID string) {
	stmt, err = db.Prepare("insert into user (id) values(?);")
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec(USER_ID)
	checkErr(err)
}
func addEvent(API_KEY string, USER_ID string, UNIX_TIMESTAMP string) {
	var tableName = "events_for_api_" + API_KEY
	stmt, err = db.Prepare("insert into " + tableName + " (userID,ts) values(?,?);")
	defer stmt.Close()
	i, err := strconv.ParseInt(UNIX_TIMESTAMP, 10, 64)
	tm := time.Unix(i, 0)
	_, err = stmt.Exec(USER_ID, tm)
	checkErr(err)
}
func addResponseTime(randSleepMS string) {
	stmt, err = db.Prepare("insert into responsetimes (responsetime) values(?);")
	defer stmt.Close()
	checkErr(err)
	_, err = stmt.Exec(randSleepMS)
	checkErr(err)
	/*if the response time is added successfully to database, make the isBucketChanged 1, so that
	when there is a request for get method, response times from database is retrived again
	before printing results*/
	/*if err == nil { //uncomment this part if needed
		isBucketChanged = 1
	}*/
}

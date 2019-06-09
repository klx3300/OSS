package dbconn

import (
	"candies"
	"database/sql"
	"evbus"
	"logger"

	_ "github.com/go-sql-driver/mysql" // well i have to use it
)

// Db is the database connection
var Db *sql.DB

var log = logger.Logger{
	Name:     "DBConnector:",
	LogLevel: 0,
}

// Initialization initializes the connection to database
func Initialization(ebus *evbus.EventBus) {
	readDBAccountFromDisk(ebus)
	log.Debugln("Start connecting to database...")
	dbconn, dberr := sql.Open("mysql", dbaccount.toDSN())
	if dberr != nil {
		log.FailLn("Failed to connect to database:", dberr.Error())
		ebus.PublishEvent(candies.NewEvent("DBConnectFail:" + dberr.Error()).Type(candies.EventTypeUnrecoverableFail).Fin())
	}
	Db = dbconn
	log.Debugln("Starting database PING...")
	dberr = Db.Ping()
	if dberr != nil {
		log.FailLn("Failed to connect to database:", dberr.Error())
		ebus.PublishEvent(candies.NewEvent("DBConnectFail:" + dberr.Error()).Type(candies.EventTypeUnrecoverableFail).Fin())
	}
	log.Succln("Successfully connected to database!")
	log.Succln("Initialization OK")
}

// Finalization closes the database connection.
func Finalization(ebus *evbus.EventBus) {
	log.Debugln("Closing database connection...")
	Db.Close()
}

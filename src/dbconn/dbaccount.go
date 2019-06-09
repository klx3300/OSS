package dbconn

import (
	"candies"
	"encoding/json"
	"evbus"
	"staticres"
	"strconv"
)

// DBAccount means the database access account information
type DBAccount struct {
	User   string
	Pass   string
	Addr   string
	Port   int
	Dbname string
}

var dbaccount DBAccount

func (dbac DBAccount) toDSN() string {
	return dbac.User + ":" + dbac.Pass + "@tcp(" + dbac.Addr + ":" + strconv.Itoa(dbac.Port) + ")/" + dbac.Dbname
}

func readDBAccountFromDisk(ebus *evbus.EventBus) {
	log.Debugln("Start reading database information from dbconn.json")
	dbafres, dbaferr := staticres.GetResource("DBAccount", "dbconn.json")
	if dbaferr != nil {
		log.FailLn("Unable to read database access information from disk:", dbaferr.Error())
		ebus.PublishEvent(candies.NewEvent("DatabaseInformationAccessFail:" + dbaferr.Error()).Type(candies.EventTypeUnrecoverableFail).Fin())
		return
	}
	jsonerr := json.Unmarshal(dbafres.Content, &dbaccount)
	if jsonerr != nil {
		log.FailLn("Database access information illegal:", jsonerr.Error())
		ebus.PublishEvent(candies.NewEvent("DatabaseInformationAccessFail:" + jsonerr.Error()).Type(candies.EventTypeUnrecoverableFail).Fin())
		return
	}
	log.Succln("Database access information OK.")
}

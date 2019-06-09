package candies

import (
	"database/sql"
	"evbus"
)

// FastSQLPrep helps you ignore the errors
func FastSQLPrep(pstmt string, db *sql.DB, ebus *evbus.EventBus) *sql.Stmt {
	s, e := db.Prepare(pstmt)
	if e != nil {
		FastGG("Programmed Wrong SQL: "+e.Error(), ebus)
		return nil
	}
	return s
}

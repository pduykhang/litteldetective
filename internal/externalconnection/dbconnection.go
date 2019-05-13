package externalconnection

import "github.com/globalsign/mgo"

type (
	DBConnection struct {
		Session *mgo.Session
	}
)

func (db *DBConnection) Close() error {
	if db.Session != nil {
		db.Session.Close()
	}
	return nil
}

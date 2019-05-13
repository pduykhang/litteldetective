package dataconnection

import (
	"github.com/globalsign/mgo"

	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
)

func GetConnection() (*mgo.Session, error) {
	l := flog.New()
	l.SetLocal("dataconnection")

	mgoss, err := mgo.Dial("mongodb:27017")
	if err != nil {
		l.Errorf("have error when connect to database with")
		return nil, err
	}
	mgoss.SetMode(mgo.Monotonic, true)
	l.Errorf("connection to mongodb at %v is successfully", "mongodb:27017")
	return mgoss, nil
}

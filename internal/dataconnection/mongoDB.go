package dataconnection

import (
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
	"github.com/globalsign/mgo"
	"time"
)

type (
	MongoConf struct {
		Addrs    []string      `envconfig:"MONGODB_ADDRS" default:"127.0.0.1:27017"`
		Database string        `envconfig:"MONGODB_DATABASE" default:"imdbfilms"`
		Username string        `envconfig:"MONGODB_USERNAME"`
		Password string        `envconfig:"MONGODB_PASSWORD"`
		Timeout  time.Duration `envconfig:"MONGODB_TIMEOUT" default:"10s"`
	}
)

func GetConnection() (*mgo.Session, error) {
	l := flog.New()
	l.SetLocal("dataconnection")
	var cof MongoConf

	mgoss, err := mgo.Dial("mongodb:27017")
	if err != nil {
		l.Errorf("have error when connect to database with")
		return nil, err
	}
	mgoss.SetMode(mgo.Monotonic, true)
	l.Errorf("connection to mongodb at %v is successfully", cof.Addrs)
	return mgoss, nil
}

package mongodb

import (
	"github.com/globalsign/mgo"

	"github.com/PhamDuyKhang/littledetective/internal/pkg/envconf"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
)

type (
	MongoConf struct {
		MongoAddr string `json:"mongo_addr",envconfig:"MONGOADDR"`
	}
)

var l flog.Logger = flog.New()

func GetConnection() (*mgo.Session, error) {
	l.SetLocal("externalconnection")
	var conf MongoConf
	err := envconf.Load(&conf)
	if err != nil {
		l.Errorf("%v", err)
		return nil, err
	}
	mongoSession, err := mgo.Dial(conf.MongoAddr)
	if err != nil {
		l.Errorf("have error when connect to database with")
		return nil, err
	}
	mongoSession.SetMode(mgo.Monotonic, true)
	l.Infof("connection to mongodb at %v is successfully", conf.MongoAddr)
	return mongoSession, nil
}

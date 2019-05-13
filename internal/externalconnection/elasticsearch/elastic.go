package elasticsearch

import (
	"fmt"
	"time"

	"github.com/olivere/elastic"
	"github.com/pkg/errors"

	"github.com/PhamDuyKhang/littledetective/internal/pkg/envconf"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
)

type (
	ElasticSearchConfig struct {
		ELasticAddr string `json:"mongo_addr",envconfig:"ELASTICADDR"`
	}
)

const MAX_RETRY = 5

var l flog.Logger = flog.New()

func GetElasticConnection() (*elastic.Client, error) {
	l.SetLocal("dataconnection")
	counting := 0
	var client *elastic.Client
	l.Infof("connecting to elastic search")
	var conf ElasticSearchConfig
	err := envconf.Load(&conf)
	if err != nil {
		return nil, err
	}
	for counting <= MAX_RETRY {
		cl, err := elastic.NewClient(elastic.SetURL(conf.ELasticAddr))
		if err != nil {
			l.Errorf("connection to elastic cluster is fail \n retry connect")
			counting++
			time.Sleep(2 * time.Second)
		} else {
			client = cl
			break
		}
	}
	if counting > MAX_RETRY {
		return nil, errors.New(fmt.Sprintf("can't connection to elastic server at %s", conf.ELasticAddr))
	}
	l.Infof("connecting  to elastic cluster  at %s is successfully ", conf.ELasticAddr)
	return client, nil
}

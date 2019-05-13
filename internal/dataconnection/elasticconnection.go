package dataconnection

import (
	"github.com/olivere/elastic"

	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
)

func GetElasticConnection() (*elastic.Client, error) {
	l := flog.New()
	l.SetLocal("dataConnection")

	client, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		l.Errorf("connection to elastic cluster is fail")
		return nil, err
	}
	l.Infof("connection to elastic cluster successfully ")
	return client, nil
}

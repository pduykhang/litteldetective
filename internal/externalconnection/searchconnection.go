package externalconnection

import "github.com/olivere/elastic"

type SearchConnection struct {
	Client *elastic.Client
}

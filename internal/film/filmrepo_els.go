package film

import (
	"context"
	"fmt"
	"reflect"

	"github.com/olivere/elastic"

	"github.com/PhamDuyKhang/littledetective/internal/types"
)

type (
	ElasticFilmRepository struct {
		client *elastic.Client
	}
)

func NewElasticFilmRepository(c *elastic.Client) *ElasticFilmRepository {
	return &ElasticFilmRepository{
		client: c,
	}
}
func (el *ElasticFilmRepository) Search(text string) (types.SearchResult, error) {
	var resultType types.Film
	termQuery := elastic.NewQueryStringQuery(text)
	searchResult, err := el.client.Search().
		Index("imdb").
		Query(termQuery).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return types.SearchResult{}, err
	}
	if searchResult.TotalHits() == 0 {
		return types.SearchResult{
			TextSearch:   text,
			TookTime:     searchResult.TookInMillis,
			NumberResult: searchResult.TotalHits(),
			Result:       nil,
		}, nil
	}
	listFilm := []types.Film{}
	for _, item := range searchResult.Each(reflect.TypeOf(resultType)) {
		t := item.(types.Film)
		listFilm = append(listFilm, t)
	}
	return types.SearchResult{
		TextSearch:   text,
		TookTime:     searchResult.TookInMillis,
		NumberResult: searchResult.TotalHits(),
		Result:       listFilm,
	}, nil
}
func (el *ElasticFilmRepository) InsertData(film types.Film) (string, error) {
	put, err := el.client.Index().Index("imdb").Type("film").Id(film.ID).BodyJson(&film).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return put.Index, nil
}

package film

import (
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
	"github.com/PhamDuyKhang/littledetective/internal/types"
)

type (
	Filmer interface {
		Search(text string) (types.SearchResult, error)
		InsertData(film types.Film) (string, error)
	}
	FilmService struct {
		l flog.Longer
		f Filmer
	}
)

func NewFilmService(f Filmer, longer flog.Longer) *FilmService {
	longer.SetLocal("film")
	return &FilmService{
		l: longer,
		f: f,
	}
}
func (fs *FilmService) FulTextSearch(searchText string) (types.SearchResult, error) {
	searchRespond, err := fs.f.Search(searchText)
	if err != nil {
		return types.SearchResult{}, err
	}
	if err != nil {
		return types.SearchResult{}, err
	}
	return searchRespond, nil
}
func (fs *FilmService) InsertDataToElastic(film types.Film) (string, error) {
	return fs.f.InsertData(film)
}

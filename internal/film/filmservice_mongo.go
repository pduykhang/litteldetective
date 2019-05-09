package film

import "github.com/PhamDuyKhang/littledetective/internal/types"

type (
	MongoFilmer interface {
		InsertFilm(f types.Film) (types.Film, error)
		GetFilmByID(ID string) (types.Film, error)
		GetAllFilm() ([]types.Film, error)
	}
	FilmMongoService struct {
		repo MongoFilmer
	}
)

func NewFilmMongoService(r MongoFilmer) *FilmMongoService {
	return &FilmMongoService{
		repo: r,
	}
}
func (s FilmMongoService) GetFilmByID(id string) (types.Film, error) {
	return s.repo.GetFilmByID(id)
}
func (s FilmMongoService) GetAllFilm() ([]types.Film, error) {
	return s.repo.GetAllFilm()
}
func (s FilmMongoService) AddFilm(f types.Film) (types.Film, error) {
	return s.repo.InsertFilm(f)
}

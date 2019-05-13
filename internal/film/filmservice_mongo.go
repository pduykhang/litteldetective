package film

import (
	"github.com/PhamDuyKhang/littledetective/internal/pkg/crawler"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
	"github.com/PhamDuyKhang/littledetective/internal/types"
	"sync"
)

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
	if f.ID != "" {
		f.ID = crawler.NewUUID()
	}
	return s.repo.InsertFilm(f)
}
func (s FilmMongoService) GenerateFilm() error {
	l := flog.New()
	l.SetLocal("film")
	filmIn := make(chan types.Film, 20)
	filmOut := make(chan types.Film, 20)
	go func() {
		var wg sync.WaitGroup
		for i := 1; i <= 20; i++ {
			wg.Add(1)
			go crawler.ExtractDetail(i, &wg, filmIn, filmOut, l)
		}
		wg.Wait()
	}()
	go func() {
		for {
			select {
			case film, ok := <-filmOut:
				if !ok {
					l.Infof("all data is saved")
					return
				}
				if film.ID != "" {
					film.ID = crawler.NewUUID()
				}
				_, err := s.repo.InsertFilm(film)
				if err != nil {
					l.Errorf("can't save data to database %v", err)
				}
				l.Infof("%s is saved", film.Title)
			}
		}
	}()
	l.Infof("extraction is stared")
	err := crawler.MakeURLTopRate(filmIn, l)
	if err != nil {
		return err
	}
	l.Infof("extraction is done close channel")
	defer close(filmIn)
	defer close(filmOut)
	return nil
}

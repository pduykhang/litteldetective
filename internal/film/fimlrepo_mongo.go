package film

import (
	imdb "github.com/PhamDuyKhang/littledetective/internal/pkg/crawler"
	"github.com/PhamDuyKhang/littledetective/internal/types"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type (
	FilmMongoRepository struct {
		ss *mgo.Session
	}
)

func NewFilmMongoRepository(s *mgo.Session) *FilmMongoRepository {
	return &FilmMongoRepository{
		ss: s,
	}
}

func (frp FilmMongoRepository) InsertFilm(f types.Film) (types.Film, error) {
	session := frp.ss.Clone()
	defer session.Clone()
	f.ID = imdb.NewUUID()
	err := session.DB("imdbfilms").C("movie").Insert(f)
	if err != nil {
		return f, err
	}
	return f, nil
}
func (frp FilmMongoRepository) GetFilmByID(ID string) (types.Film, error) {
	session := frp.ss.Clone()
	defer session.Clone()
	var film types.Film
	if err := session.DB("imdbfilms").C("movie").Find(bson.M{"_id": ID}).One(&film); err != nil {
		return types.Film{}, err
	}
	return film, nil
}
func (frp FilmMongoRepository) GetAllFilm() ([]types.Film, error) {
	session := frp.ss.Clone()
	defer session.Clone()
	var films []types.Film
	if err := session.DB("imdbfilms").C("movie").Find(bson.M{}).All(&films); err != nil {
		return nil, err
	}
	return films, nil
}

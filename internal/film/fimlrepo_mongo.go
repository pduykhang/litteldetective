package film

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/PhamDuyKhang/littledetective/internal/types"
)

const (
	DATABASE_NAME   = "imdbfilms"
	COLLECTION_NAME = "moive"
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
	defer session.Close()
	err := session.DB(DATABASE_NAME).C(COLLECTION_NAME).Insert(f)
	if err != nil {
		return f, err
	}
	return f, nil
}
func (frp FilmMongoRepository) GetFilmByID(ID string) (types.Film, error) {
	session := frp.ss.Clone()
	defer session.Close()
	var film types.Film
	if err := session.DB(DATABASE_NAME).C(COLLECTION_NAME).Find(bson.M{"_id": ID}).One(&film); err != nil {
		return types.Film{}, err
	}
	return film, nil
}
func (frp FilmMongoRepository) GetAllFilm() ([]types.Film, error) {
	session := frp.ss.Clone()
	defer session.Close()
	var films []types.Film
	if err := session.DB(DATABASE_NAME).C(COLLECTION_NAME).Find(bson.M{}).All(&films); err != nil {
		return nil, err
	}
	return films, nil
}

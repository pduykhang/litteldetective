package handler

import (
	imdb "github.com/PhamDuyKhang/littledetective/internal/pkg/crawler"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/marshal"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/respond"
	"github.com/PhamDuyKhang/littledetective/internal/types"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	FilmService interface {
		GetFilmByID(id string) (types.Film, error)
		GetAllFilm() ([]types.Film, error)
		AddFilm(f types.Film) (types.Film, error)
	}
	FilmHandler struct {
		s  FilmService
		el SearchService
		l  flog.Longer
	}
	CrawlURL struct {
		URL string `json:"url"`
	}
)

func NewFilmHandler(s FilmService, logger flog.Longer, fr SearchService) *FilmHandler {
	logger.SetLocal("handler")
	return &FilmHandler{
		s:  s,
		el: fr,
		l:  logger,
	}
}
func (h FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	listFilm, err := h.s.GetAllFilm()
	if err != nil {
		h.l.Errorf("error when insert data %v", err)
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"status": "500", "message": "error when insert data"})
		return
	}
	respond.JSON(w, http.StatusAccepted, listFilm)
	return
}
func (h FilmHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	film := types.Film{}
	err := marshal.ParseRequest(r, &film)
	if err != nil {
		h.l.Errorf("can't parse data form http request err: %v", err)
		respond.JSON(w, http.StatusBadRequest, map[string]string{"status": "400", "message": "can't get content form your request"})
		return
	}
	newFilm, err := h.s.AddFilm(film)
	if err != nil {
		h.l.Errorf("have error when insert data into database: %v", err)
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"status": "500", "message": "have error when insert data into database"})
		return
	}
	_, err = h.el.InsertDataToElastic(film)
	if err != nil {
		h.l.Errorf("have error when insert data into elastic search server: %v", err)
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"status": "500", "message": "your data is not synchronized"})
	}
	respond.JSON(w, http.StatusAccepted, newFilm.ID)
	return

}
func (h FilmHandler) GetFilmWithID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	film, err := h.s.GetFilmByID(id)
	if err != nil {
		h.l.Errorf("error when insert data %v", err)
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"status": "500", "message": "error when insert data"})
		return
	}
	respond.JSON(w, http.StatusAccepted, film)
	return
}
func (h FilmHandler) GenerateData(w http.ResponseWriter, r *http.Request) {
	imdb.Crawler()
	respond.JSON(w, http.StatusAccepted, map[string]string{"status": "200", "message": "generation is successfully "})
	return
}
func (h FilmHandler) Sync(w http.ResponseWriter, r *http.Request) {
	listFilm, err := h.s.GetAllFilm()
	if err != nil {
		h.l.Errorf("error when get list data %v", err)
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"status": "500", "message": "error when get data"})
		return
	}
	for _, film := range listFilm {
		id, err := h.el.InsertDataToElastic(film)
		if err != nil {

		}
		h.l.Infof("insert %s with index := %s", film.Title, id)
	}
	return
}

package handler

import (
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/marshal"
	"github.com/PhamDuyKhang/littledetective/internal/pkg/respond"
	"github.com/PhamDuyKhang/littledetective/internal/types"
	"net/http"
)

type (
	SearchService interface {
		FulTextSearch(searchText string) (types.SearchResult, error)
		InsertDataToElastic(film types.Film) (string, error)
	}
	SearchHandler struct {
		s      SearchService
		logger flog.Longer
	}
)

func NewSearchHandeler(s SearchService, l flog.Longer) *SearchHandler {
	l.SetLocal("handler")
	return &SearchHandler{
		s:      s,
		logger: l,
	}
}
func (h SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	requestData := types.SearchRequest{}
	err := marshal.ParseRequest(r, &requestData)
	if err != nil {
		h.logger.Errorf("can't parse data form http request err: %v", err)
		respond.JSON(w, http.StatusBadRequest, map[string]string{"status": "400", "message": "can't get content form your request"})
		return
	}
	result, err := h.s.FulTextSearch(requestData.Keyword)
	if err != nil {
		h.logger.Errorf("search fail with err: %v", err)
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"status": "500", "message": "have error when search "})
		return
	}
	respond.JSON(w, http.StatusAccepted, result)
	return
}

func (h SearchHandler) Test(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusAccepted, map[string]string{"status": "200", "result": "Hello"})
	return
}
func (h SearchHandler) AddFilmToElastic(w http.ResponseWriter, r *http.Request) {
	filmData := types.Film{}
	err := marshal.ParseRequest(r, &filmData)
	if err != nil {
		h.logger.Errorf("can't parse data form http request err: %v", err)
		respond.JSON(w, http.StatusBadRequest, map[string]string{"status": "400", "message": "can't get content form your request"})
		return
	}
	index, err := h.s.InsertDataToElastic(filmData)
	if err != nil {
		h.logger.Errorf("error when insert data to elastic server err: %v", err)
		respond.JSON(w, http.StatusInternalServerError, map[string]string{"status": "500", "message": "data wasn't inserted "})
		return
	}
	respond.JSON(w, http.StatusAccepted, map[string]string{"status": "200", "message": index})
	return
}

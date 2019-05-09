package imdb

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/globalsign/mgo"

	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
)

type (
	Film struct {
		ID          string              `json:"id" bson:"_id"`
		Rank        int                 `json:"rank" bson:"rank"`
		URL         string              `json:"url" bson:"url"`
		Title       string              `json:"title" bson:"title"`
		Rate        string              `json:"rate" bson:"rate"`
		ReleaseDate int                 `json:"release_date" bson:"release_date"`
		Description string              `json:"description" bson:"description"`
		Credit      map[string][]string `json:"credit" bson:"credit"`
	}
)

func Crawler() {
	l := flog.New()
	l.SetLocal("imdb")
	filmIn := make(chan Film, 20)
	filmOut := make(chan Film, 20)
	go func() {
		var wg sync.WaitGroup
		for i := 1; i <= 20; i++ {
			wg.Add(1)
			go ExtractDetail(i, wg, filmIn, filmOut, l)
		}
		wg.Wait()
	}()
	go func() {
		SaveData(filmOut, l)
	}()
	l.Infof("extraction is stared")
	MakeURLTopRate(filmIn, l)
	l.Infof("extraction is done close channel")
	close(filmIn)
	return

}
func MakeURLTopRate(filmURLChan chan Film, l *flog.Flog) {
	topURL := "https://www.imdb.com/chart/top?ref_=nv_mv_250"
	doc, err := GetDocFormURL(topURL)
	if err != nil {
		fmt.Println(err)
	}
	doc.Find("table.chart").Find("tbody.lister-list").Find("tr").Each(func(rank int, films *goquery.Selection) {
		f := Film{}
		f.Rank = rank + 1
		filmURL, ok := films.Find("td.titleColumn").Find("a").Attr("href")
		if ok {
			f.URL = NormalizeURL(filmURL)
		}
		f.Title = films.Find("td.titleColumn").Find("a").Text()
		rawRelease := films.Find("td.titleColumn").Find("span.secondaryInfo").Text()
		rawRelease = strings.TrimLeft(rawRelease, "(")
		rawRelease = strings.TrimRight(rawRelease, ")")
		releaseDate, err := strconv.Atoi(rawRelease)
		if err != nil {
			fmt.Println(err)
			f.ReleaseDate = 0
		} else {
			f.ReleaseDate = releaseDate
		}
		f.Rate = films.Find("td.imdbRating").Find("strong").Text()
		f.ID = NewUUID()
		filmURLChan <- f
	})

}

func ExtractDetail(idx int, wg sync.WaitGroup, fimIn chan Film, filmOut chan Film, l *flog.Flog) {
	defer wg.Done()
	for {
		select {
		case film, ok := <-fimIn:
			if !ok {
				l.Infof("crawler #%d done", idx)
				return
			}
			l.Infof("crawler #%d is crawling", idx)
			newFilm := film
			l.Infof("url %s being extracted", newFilm.URL)
			filmDoc, err := GetDocFormURL(newFilm.URL)
			if err != nil {
				l.Errorf("%s extraction is fail", newFilm.Title)
			}
			film.Description = strings.TrimSpace(filmDoc.Find("div.plot_summary").Find("div.summary_text").Text())
			m := make(map[string][]string)
			filmDoc.Find("div.plot_summary").Find("div.credit_summary_item").Each(func(i int, item *goquery.Selection) {
				key := strings.TrimRight(item.Find("h4.inline").Text(), ":")
				value := []string{}
				item.Find("a").Each(func(j int, values *goquery.Selection) {
					if !strings.HasPrefix(values.Text(), "See full cast") && !strings.HasSuffix(values.Text(), "credits ") && !strings.HasPrefix(values.Text(), "1 more credit") {
						value = append(value, values.Text())
					}
				})
				m[key] = value
			})
			film.Credit = m
			l.Infof("%s is pushed to film channel to save ", newFilm.Title)
			filmOut <- film
		}
	}
}
func NormalizeURL(url string) (fullUrl string) {
	rootURL := "https://www.imdb.com"
	fullUrl = rootURL + url
	return
}
func GetDocFormURL(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36")
	Client := http.Client{}
	body, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(body.Body)
	body.Body.Close()
	if err != nil {
		return nil, err
	}
	return doc, nil
}
func SaveData(in chan Film, l *flog.Flog) {
	session, err := mgo.Dial("mongodb:27017")
	if err != nil {
		l.Errorf("dialing is fail %v", err)
		return
	}
	session.SetMode(mgo.Monotonic, true)
	s := session.Clone()
	for {
		select {
		case listFilms, ok := <-in:
			if !ok {
				l.Infof("all data is saved")
				return
			}
			err = s.DB("imdbfilms").C("movie").Insert(listFilms)
			l.Infof("%s is saved", listFilms.Title)
			if err != nil {
				l.Errorf("can't save data to database %v", err)
			}
		}
	}
}

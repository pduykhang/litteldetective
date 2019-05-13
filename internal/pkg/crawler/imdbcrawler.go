package crawler

import (
	"fmt"
	"github.com/PhamDuyKhang/littledetective/internal/types"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/globalsign/mgo"

	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
)

func MakeURLTopRate(filmURLChan chan types.Film, l *flog.Flog) error {
	topURL := "https://www.imdb.com/chart/top?ref_=nv_mv_250"
	doc, err := GetDocFormURL(topURL)
	if err != nil {
		return err
	}
	doc.Find("table.chart").Find("tbody.lister-list").Find("tr").Each(func(rank int, films *goquery.Selection) {
		f := types.Film{}
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
	return nil
}

func ExtractDetail(idx int, wg *sync.WaitGroup, fimIn chan types.Film, filmOut chan types.Film, l *flog.Flog) {
	defer wg.Done()
	for {
		select {
		case film, ok := <-fimIn:
			if !ok {
				l.Infof("crawler #%d done", idx)
				return
			}
			l.Infof("crawler #%d is crawling", idx)
			l.Infof("url %s being extracted", film.URL)
			filmDoc, err := GetDocFormURL(film.URL)
			if err != nil {
				l.Errorf("%s extraction is fail", film.Title)
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
			l.Infof("%s is pushed to film channel to save ", film.Title)
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
func SaveData(in chan types.Film, l *flog.Flog, session *mgo.Session) {
	defer session.Close()
	for {
		select {
		case listFilms, ok := <-in:
			if !ok {
				l.Infof("all data is saved")
				return
			}
			err := session.DB("imdbfilms").C("movie").Insert(listFilms)
			if err != nil {
				l.Errorf("can't save data to database %v", err)
			}
			l.Infof("%s is saved", listFilms.Title)
		}
	}
}

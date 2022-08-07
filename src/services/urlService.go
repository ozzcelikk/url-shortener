package services

import (
	"encoding/json"
	"github.com/ozzcelikk/url-shortener/src/helpers"
	"net/http"
	"strings"
)

type urlPostModel struct {
	Url string `json:"url"`
}

type createResponseModel struct {
	Success     bool
	UserMessage string
	ShortUrl    string
}

type listResponseModel struct {
	Success     bool
	UserMessage string
	Stats       []UrlData
}

func find(url string) *UrlData {
	for _, n := range *store.Data {
		if n.Url == url {
			return &n
		}
	}

	return nil
}
func findByShorUrl(url string) (urlData *UrlData, i int) {
	for i, n := range *store.Data {
		if n.ShortUrl == url {
			return &n, i
		}
	}

	return nil, -1
}

func generateShortUrl() string {
	isUrlAvailable := false

	shortUrl := ""

	for isUrlAvailable == false {
		shortUrl = helpers.GetRandomShortUrl()

		c, _ := findByShorUrl(shortUrl)

		if c == nil {
			isUrlAvailable = true
		}
	}

	return shortUrl
}
func createUrlData(postModel urlPostModel) createResponseModel {
	existUrlData := find(postModel.Url)

	if existUrlData != nil {
		return createResponseModel{
			Success:     false,
			UserMessage: "Url already exist",
		}
	}

	checkUrl := helpers.IsValidUrl(postModel.Url)

	if checkUrl == false {
		return createResponseModel{
			Success:     false,
			UserMessage: "Invalid Url",
		}
	}

	urlData := UrlData{
		Url:        postModel.Url,
		ShortUrl:   generateShortUrl(),
		VisitCount: 0,
	}

	*store.Data = append(*store.Data, urlData)

	return createResponseModel{
		Success:     true,
		ShortUrl:    strings.Join([]string{"http://localhost:8080", urlData.ShortUrl}, "/"),
		UserMessage: "Success",
	}
}
func getList() listResponseModel {
	if len(*store.Data) != 0 {
		var res = new(listResponseModel)
		res.Success = true
		res.UserMessage = "Success"

		statsArr := new([]UrlData)

		for _, n := range *store.Data {
			s := UrlData{
				Url:        n.Url,
				ShortUrl:   n.ShortUrl,
				VisitCount: n.VisitCount,
			}

			*statsArr = append(*statsArr, s)
		}

		return listResponseModel{
			Success:     true,
			UserMessage: "Success",
			Stats:       *statsArr,
		}
	}

	return listResponseModel{
		Success:     true,
		UserMessage: "Success",
		Stats:       nil,
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var postModel urlPostModel
	err := json.NewDecoder(r.Body).Decode(&postModel)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	res := createUrlData(postModel)
	saveStore()

	json.NewEncoder(w).Encode(res)
}
func GetHandler(w http.ResponseWriter, r *http.Request) {
	url := strings.Replace(r.URL.Path, "/", "", 1)

	urlData, i := findByShorUrl(url)

	(*store.Data)[i].VisitCount += 1
	saveStore()

	http.Redirect(w, r, urlData.Url, http.StatusSeeOther)
}
func ListHandler() listResponseModel {
	return getList()
}

func Init() {
	loadStore()
}

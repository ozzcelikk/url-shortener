package services

import (
	"encoding/json"
	"io/ioutil"
)

type UrlData struct {
	Url        string
	ShortUrl   string
	VisitCount int
}

type Store struct {
	Data *[]UrlData `json:"Items"`
}

var store Store

var filePath = "./store.json"

func loadStore() bool {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(bytes, &store)

	return true
}
func saveStore() bool {
	content, err := json.Marshal(store)

	err = ioutil.WriteFile(filePath, content, 0644)

	if err != nil {
		panic(err)
	}

	return true
}

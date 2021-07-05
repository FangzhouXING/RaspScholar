package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sp "github.com/fangzhouxing/raspscholar/scholarparser"
)

type Paper struct {
	Title         string
	Year          int
	CitationCount int
}

type Author struct {
	Name   string
	School string
	Focus  []string
	Papers []Paper
}

type ParserConfig struct {
	ScholarCode string
	UseScraper  bool
	ScraperKey  string
}

func (c ParserConfig) GetScholarCode() string {
	return c.ScholarCode
}

func (c ParserConfig) GetScraperInfo() (bool, string) {
	return c.UseScraper, c.ScraperKey
}

func ReadConfig(filename string) (ParserConfig, error) {
	configfile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Unable to read config.")
	}

	var config ParserConfig
	err = json.Unmarshal(configfile, &config)
	if err != nil {
		panic("config parsing failed.")
	}

	fmt.Println(config)
	return config, err
}

const DEFAULT_CONFIG = "CONFIG.json"

func ReadDefaultConfig() (ParserConfig, error) {
	return ReadConfig(DEFAULT_CONFIG)
}

func fetchCachedScholar(filename string) sp.Scholar {
	var scholar sp.Scholar
	bytes, _ := ioutil.ReadFile(filename)
	json.Unmarshal(bytes, scholar)
	return scholar
}
func writeScholarToCache(scholar sp.Scholar, filename string) {
	bytes, _ := json.Marshal(scholar)
	ioutil.WriteFile(filename, bytes, 0644)

}

func FetchScholar(config ParserConfig) sp.Scholar {
	useCache := true

	var scholar sp.Scholar
	if !useCache {
		scholar, _ = sp.GetScholar(config)
		writeScholarToCache(scholar, "parsedScholar.json")
	} else {
		scholar = fetchCachedScholar("parsedScholar.json")
	}
	return scholar
}

func main() {

	config, _ := ReadConfig("CONFIG_real.json")

	fmt.Println(config)
	FetchScholar(config)

	http.Handle("/", http.FileServer(http.Dir("./Website")))

	http.ListenAndServe(":10000", nil)

}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	sp "github.com/fangzhouxing/raspscholar/scholarparser"
)

var gScholar sp.Scholar

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
	err := json.Unmarshal(bytes, &scholar)
	if err != nil {
		fmt.Println("fetchCacheScholar error. ", err)
	}
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

func handleScholarBasicInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In BasicInfo.", r.URL)
	switch r.Method {
	case "GET":
		fmt.Printf("GET URL: %v\n", r.URL)
		w.Header().Set("Content-Type", "application/json")
		bytes, _ := json.Marshal(gScholar)
		w.Write(bytes)
	default:
		fmt.Printf("Unexpected method: %v.\n", r.Method)
	}
}

func handleRecentPapers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In RecentPaper.", r.URL)
}

type slashFix struct {
	mux http.Handler
}

func (h *slashFix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.Replace(r.URL.Path, "//", "/", -1)
	h.mux.ServeHTTP(w, r)
}

func main() {

	config, _ := ReadConfig("CONFIG_real.json")

	fmt.Println(config)
	gScholar = FetchScholar(config)
	fmt.Println("Scholar name: ", gScholar.Name)

	httpMux := http.NewServeMux()

	fs := http.FileServer(http.Dir("Website"))
	httpMux.Handle("/", fs)
	httpMux.HandleFunc("/bascinfo/", handleScholarBasicInfo)
	httpMux.HandleFunc("/recentpapers/", handleRecentPapers)

	//http.Handle("/static/", http.FileServer(http.Dir("Website/static")))

	http.ListenAndServe(":10000", &slashFix{httpMux})

}

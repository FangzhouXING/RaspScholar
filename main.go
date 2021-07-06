package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	sp "github.com/fangzhouxing/raspscholar/scholarparser"
)

const MAX_GLOBAL_COPY = 5
const POLL_FREQUENCY = time.Hour //A bad name
const POLL_VARITION = time.Hour / 4

type ScholarVersion struct {
	scolarlist []sp.Scholar
	mu         sync.Mutex

	lastUpdate time.Time
	maxLen     int
}

var gScholarVersiion *ScholarVersion = NewScolarVersion(MAX_GLOBAL_COPY)

func NewScolarVersion(maxLen int) *ScholarVersion {
	return &ScholarVersion{
		scolarlist: []sp.Scholar{},
		mu:         sync.Mutex{},
		lastUpdate: time.Now(),
		maxLen:     maxLen,
	}
}

func (sv *ScholarVersion) AddVersion(newScholar sp.Scholar) {
	sv.mu.Lock()
	if len(sv.scolarlist) >= sv.maxLen {
		sv.scolarlist = sv.scolarlist[sv.maxLen-len(sv.scolarlist)+1:]
	}
	sv.scolarlist = append(sv.scolarlist, newScholar)
	sv.lastUpdate = time.Now()

	sv.mu.Unlock()
}

func (sv *ScholarVersion) LatestVersion() sp.Scholar {
	sv.mu.Lock()
	defer sv.mu.Unlock()
	if len(sv.scolarlist) < 1 {
		return sp.Scholar{}
	}
	temp := sv.scolarlist[len(sv.scolarlist)-1]
	return temp
}
func (sv *ScholarVersion) Size() int {
	return len(sv.scolarlist)
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

type BasicInfoResponse struct {
	Name          string
	RankAndSchool string
	Focus         []string
}

type RecentPaperResponse struct {
	Papers []sp.Paper
}

func getBasicInfo(scholar sp.Scholar) BasicInfoResponse {
	return BasicInfoResponse{Name: scholar.Name, RankAndSchool: scholar.RankAndSchool, Focus: scholar.Focus}
}

func getRecentPaper(scholar sp.Scholar) RecentPaperResponse {
	return RecentPaperResponse{Papers: scholar.Papers}
}

func handleScholarBasicInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In BasicInfo.", r.URL)
	switch r.Method {
	case "GET":
		fmt.Printf("GET URL: %v\n", r.URL)
		w.Header().Set("Content-Type", "application/json")

		latestScholar := gScholarVersiion.LatestVersion()
		bytes, _ := json.Marshal(getBasicInfo(latestScholar))
		w.Write(bytes)
	default:
		fmt.Printf("Unexpected method: %v.\n", r.Method)
	}
}

func handleRecentPapers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In RecentPaper.", r.URL)
	switch r.Method {
	case "GET":
		fmt.Printf("GET URL: %v\n", r.URL)
		w.Header().Set("Content-Type", "application/json")

		latestScholar := gScholarVersiion.LatestVersion()
		bytes, _ := json.Marshal(getRecentPaper(latestScholar))
		w.Write(bytes)
	default:
		fmt.Printf("Unexpected method: %v.\n", r.Method)
	}
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

func fetchScholar(config ParserConfig) sp.Scholar {
	useCache := false

	var scholar sp.Scholar
	if !useCache {
		scholar, _ = sp.GetScholar(config)
		writeScholarToCache(scholar, "parsedScholar.json")
	} else {
		scholar = fetchCachedScholar("parsedScholar.json")
	}
	return scholar
}

func ScholarUpdateRoutine(sv *ScholarVersion, config ParserConfig) {
	sv.AddVersion(fetchScholar(config))

	var fetchAvailable int = 0
	var fetchMutex sync.Mutex
	fetchCountDown := func() {
		for {
			nextPoll := rand.Int63n(2*int64(POLL_VARITION)) + int64(POLL_FREQUENCY) - int64(POLL_VARITION)
			fmt.Println("Next poll at: ", time.Duration(nextPoll))
			timer1 := time.NewTimer(time.Duration(nextPoll))

			<-timer1.C
			fmt.Println("Adding one poll.")

			fetchMutex.Lock()
			fetchAvailable += 1
			fetchMutex.Unlock()
		}
	}

	go fetchCountDown()

	for {
		if fetchAvailable > 0 {
			fetchMutex.Lock()
			fetchAvailable--
			fetchMutex.Unlock()

			fmt.Println("Fetching new Scolar Version.")
			newScolar := fetchScholar(config)
			sv.AddVersion(newScolar)
			fmt.Printf("SV now contain %v copies\n", sv.Size())
		}
		time.Sleep(10 * time.Second)
	}

}

func main() {

	rand.Seed(time.Now().UnixNano())
	config, _ := ReadConfig("CONFIG_real.json")

	fmt.Println(config)

	go ScholarUpdateRoutine(gScholarVersiion, config)

	httpMux := http.NewServeMux()

	fs := http.FileServer(http.Dir("Website"))
	httpMux.Handle("/", fs)
	httpMux.HandleFunc("/bascinfo/", handleScholarBasicInfo)
	httpMux.HandleFunc("/recentpapers/", handleRecentPapers)

	//http.Handle("/static/", http.FileServer(http.Dir("Website/static")))

	http.ListenAndServe(":10000", httpMux)

}

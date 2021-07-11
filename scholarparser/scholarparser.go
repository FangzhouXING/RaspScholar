package scholarparser

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ScholarParserConfig interface {
	GetScholarCode() string
	GetScraperInfo() (bool, string)
}

type Paper struct {
	Title    string
	Citation int
	PubYear  int
}

type Scholar struct {
	Name          string
	RankAndSchool string
	Focus         []string
	HIndex        int
	Papers        []Paper
}

func parseCitation(text string) int {
	if text == "" || text == "*" {
		return 0
	}
	if strings.Contains(text, "*") {
		text = strings.ReplaceAll(text, "*", "")
	}
	res, err := strconv.Atoi(text)
	if err != nil {
		errstr := fmt.Sprintf("Unable to parse %v to int", text)
		panic(errstr)
	}
	return res
}

func parseYear(text string) int {
	if text == "" {
		return 0
	}
	year, err := strconv.Atoi(text)
	if err != nil {
		errstr := fmt.Sprintf("Unable to parse year text %v to int", text)
		panic(errstr)
	}
	return year
}

func parseHIndexInt(text string) int {
	if text == "" {
		return 0
	}
	hindex, err := strconv.Atoi(text)
	if err != nil {
		errstr := fmt.Sprintf("Unable to parse hindex %v to int", text)
		panic(errstr)
	}
	return hindex
}

func ParseScholarPage(page io.Reader) (Scholar, error) {
	doc, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		fmt.Printf("ParseScholarPage Error: %v", err)
		return Scholar{}, err
	}

	papers := make([]Paper, 0)

	parsePaper := func(i int, s *goquery.Selection) {
		//fmt.Printf("\nPaper No.: %d\n", i)
		//fmt.Printf("Title: %v\n", s.Find(".gsc_a_at").Text())
		//fmt.Printf("Cite: %d\n", parseCitation(s.Find(".gsc_a_c").Text()))
		//fmt.Printf("Year: %d\n", parseYear(s.Find(".gsc_a_y").Text()))
		title := s.Find(".gsc_a_at").Text()
		citation := parseCitation(s.Find(".gsc_a_c").Text())
		year := parseYear(s.Find(".gsc_a_y").Text())
		papers = append(papers, Paper{Title: title, Citation: citation, PubYear: year})
	}
	doc.Find(".gsc_a_tr").Each(parsePaper)

	fmt.Println("Paper Len:", len(papers))

	var name string
	var rankAndSchool string
	var focus []string = []string{}
	parseScholar := func(i int, s *goquery.Selection) {
		nameSelector := "#gsc_prf_in"
		focusSelector := "#gsc_prf_int"
		name = s.Find(nameSelector).Text()
		rankAndSchool = s.Find(nameSelector).Next().Text()
		//fmt.Printf("Name: %v\n", name)
		//fmt.Printf("title: %v\n", rankAndSchool)

		getFocus := func(i int, fs *goquery.Selection) {
			focus = append(focus, fs.Text())
		}
		s.Find(focusSelector).Children().Each(getFocus)
		//fmt.Printf("Focus size: %d\n", len(focus))
	}

	doc.Find("#gsc_prf").Each(parseScholar)

	var hindex int = 0
	parseHIndex := func(i int, s *goquery.Selection) {
		rowText := s.First().Text()
		fmt.Printf("rowText: %s\n", rowText)
		if rowText == "h-index" {
			numberPart := s.First().Next().Text()
			fmt.Printf("numberPart: %s\n", numberPart)
			hindex = parseHIndexInt(numberPart)
		}

	}
	doc.Find(".gsc_rsb_sc1").Each(parseHIndex)

	scholar := Scholar{Name: name, RankAndSchool: rankAndSchool, Focus: focus, Papers: papers, HIndex: hindex}
	//os.WriteFile("samplepage2.html")
	return scholar, nil
}

func buildGoogleScholarUrl(name string, start int, pagesize int) string {
	urlTemplate := "https://scholar.google.com/citations?user=%s&hl=en&cstart=%d&pagesize=%d&sortby=pubdate"
	return fmt.Sprintf(urlTemplate, name, start, pagesize)
}

func buildScraperApi(name string, start int, pagesize int, apikey string) string {
	scholarUrl := buildGoogleScholarUrl(name, start, pagesize)
	return fmt.Sprintf("http://api.scraperapi.com?api_key=%s&url=%s", apikey, scholarUrl)
}

func GetScholar(config ScholarParserConfig) (Scholar, error) {
	useScraper, scraperKey := config.GetScraperInfo()
	var reqUrl string
	if useScraper {
		reqUrl = buildScraperApi(config.GetScholarCode(), 0, 100, scraperKey)
	} else {
		reqUrl = buildGoogleScholarUrl(config.GetScholarCode(), 0, 100)
	}
	fmt.Println("reqstr:", reqUrl)

	req, e := http.NewRequest("GET", reqUrl, nil)
	if e != nil {
		panic(e)
	}

	//setHeader(req)
	fmt.Println("Created Request.")

	res, e := new(http.Client).Do(req)
	if e != nil {
		panic(e)
	}
	defer res.Body.Close()

	scholar, err := ParseScholarPage(res.Body)
	if err != nil {
		fmt.Printf("Error parsing scholar page.")
		return Scholar{}, nil
	}
	return scholar, nil
}

package scholarparser_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/fangzhouxing/raspscholar/scholarparser"
)

func TestParseScholarParsing(t *testing.T) {
	reader, err := os.Open("samplepage.html")
	if err != nil {
		t.Error("Cannot open samplepage")
	}
	scholar, err := scholarparser.ParseScholarPage(reader)
	if err != nil {
		t.Error("Parse failed.")
	}
	fmt.Println(scholar)
}

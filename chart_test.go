package timelapsechart

import (
	"testing"
)

const checkMark = "\u2713"
const xMark = "\u2717"

const basicCsv = `
Country,Date,Age,Amount,Id
"United States",2012-02-01,50,112.1,01234
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,17,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-05-07,NA,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United States",2012-02-01,32,321.31,54320
Spain,2012-02-01,66,555.42,00241
`

func TestNewChart(t *testing.T) {

	t.Log("a chart")

	chart := New("", Config{})

	if chart.Config.Width == 800 {
		t.Log(" should be 800 px wide by default.", checkMark)
	} else {
		t.Fatal(" should be 800 px wide by default.", xMark)
	}
}

package timelapsechart

import (
	"bytes"
	"strings"
	"testing"
)

func TestBarChart(t *testing.T) {

	csvStr := `
continent,country,date,amount
North America,Canada,2021-01-01,25.0
North America,United States,2021-01-01,40.0
North America,Mexico,2021-01-01,15.0
`

	t.Log("a developer")

	data := Data{}
	err := data.ReadCSV(strings.NewReader(csvStr), []string{"country", "date", "amount"})
	if err != nil {
		t.Fatal(" should be able to read the csv.", xMark, err)
	}

	chart := New("", "bar", Config{})
	chart.AddData(&data)

	buffer := bytes.NewBuffer([]byte{})

	barChart := BarChart{}

	err = barChart.render(chart, buffer)

	if err == nil {
		t.Log(" should be able to render a bar chart.", checkMark)
	} else {
		t.Fatal(" should be able to render a bar chart.", xMark, err)
	}
}

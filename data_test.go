package timelapsechart

import (
	"os"
	"strings"
	"testing"
)

func TestData(t *testing.T) {

	t.Log("a chart builder")

	data := Data{}
	data.AddCategory("Canada")

	if len(data.Categories) == 1 {
		t.Log(" should be able to add a category.", checkMark)
	} else {
		t.Fatal(" should be able to add a category.", xMark)
	}

}

func TestDataColDefaults(t *testing.T) {

	csvStr := `
country,date,total
Canada,2021-01-01,25.0
United States,2021-01-01,40.0
Mexico,2021-01-01,15.0
`
	t.Log("the category column")

	data := Data{}
	err := data.ReadCSV(strings.NewReader(csvStr), []string{"country", "date", "total"})
	if err != nil {
		t.Fatal(" should be able to read the csv.", xMark, err)
	}

	colNum, err := data.getCategoryColNum()
	if colNum == 0 {
		t.Log(" should be 0 by default.", checkMark)
	} else {
		t.Fatal(" should be 0 by default.", xMark, colNum)
	}
}

func TestDataConfig(t *testing.T) {

	t.Log("a developer")

	csvFile, err := os.Open("example/input.csv")
	if err != nil {
		t.Fatal(" should be able to open the input test file.", xMark, err)
	}
	defer csvFile.Close()

	data := Data{}
	err = data.ReadCSV(csvFile, []string{"country", "date", "total"})
	if err != nil {
		t.Fatal(" should be able to read the csv.", xMark, err)
	}

	colNum, err := data.getCategoryColNum()
	if err != nil {
		t.Fatal(" should be able to get the col num.", xMark, err)
	}
	if colNum == 0 {
		t.Log(" should be able to specify the category column name for an input csv.", checkMark)
	} else {
		t.Fatal(" should be able to specify the category column name for an input csv.", xMark, colNum)
	}

	csvFile.Close()

	csvStr := `
continent,country,date,amount
North America,Canada,2021-01-01,25.0
North America,United States,2021-01-01,40.0
North America,Mexico,2021-01-01,15.0
`

	err = data.ReadCSV(strings.NewReader(csvStr), []string{"country", "date", "amount"})
	if err != nil {
		t.Fatal(" should be able to read the csv.", xMark, err)
	}

	colNum, err = data.getCategoryColNum()
	if err != nil {
		t.Fatal(" should be able to get the col num.", xMark, err)
	}
	if colNum == 1 {
		t.Log(" regardless of the order of columns in the csv.", checkMark)
	} else {
		t.Fatal(" regardless of the order of columns in the csv", xMark, colNum)
	}

	csvFile.Close()

}

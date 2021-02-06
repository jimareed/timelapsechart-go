package timelapsechart

import (
	"context"
	"errors"
	"io"

	"github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/imports"
)

// Data
type Data struct {
	Categories     []string `json:"categories"`
	TimeRange      []string `json:"timeRange"`
	MaxValue       float64  `json:"maxValue"`
	CategoryColumn string   `json:"categoryColumn"`
	TimeColumn     string   `json:"timeColumn"`
	ValueColumn    string   `json:"valueColumn"`
	Value2Column   string   `json:"value2Column"`
	seriesData     *dataframe.DataFrame
	series2Data    *dataframe.DataFrame
	inputData      *dataframe.DataFrame
}

var ctx = context.TODO()

func (data *Data) getCategoryColNum() (int, error) {
	if data.CategoryColumn == "" {
		return 0, nil
	}
	return data.inputData.NameToColumn(data.CategoryColumn)
}

func (data *Data) getTimeColNum() (int, error) {
	if data.TimeColumn == "" {
		return 1, nil
	}
	return data.inputData.NameToColumn(data.TimeColumn)
}

func (data *Data) getValueColNum() (int, error) {
	if data.ValueColumn == "" {
		return 2, nil
	}
	return data.inputData.NameToColumn(data.ValueColumn)
}

func (data *Data) getValue2ColNum() (int, error) {
	if data.Value2Column == "" {
		return -1, nil
	}
	return data.inputData.NameToColumn(data.Value2Column)
}

func (data *Data) columnsDefined() bool {
	return data.CategoryColumn != "" && data.TimeColumn != "" && data.ValueColumn != ""
}

func (data *Data) ReadCSV(r io.ReadSeeker, cols []string) error {

	if len(cols) > 0 {
		data.CategoryColumn = cols[0]
	}
	if len(cols) > 1 {
		data.TimeColumn = cols[1]
	}
	if len(cols) > 2 {
		data.ValueColumn = cols[2]
	}
	if len(cols) > 3 {
		data.Value2Column = cols[3]
	}

	err := errors.New("")

	if data.columnsDefined() {
		data.inputData, err = imports.LoadFromCSV(ctx, r, imports.CSVLoadOptions{
			DictateDataType: map[string]interface{}{
				data.CategoryColumn: "",
				data.TimeColumn:     "",
				data.ValueColumn:    float64(0),
			}})
	} else {
		data.inputData, err = imports.LoadFromCSV(ctx, r, imports.CSVLoadOptions{InferDataTypes: true, NilValue: &[]string{"NA"}[0]})
	}

	if err != nil {
		return err
	}

	categoryColNum, err := data.getCategoryColNum()
	if err != nil {
		return err
	}
	timeColNum, err := data.getTimeColNum()
	if err != nil {
		return err
	}
	valueColNum, err := data.getValueColNum()
	if err != nil {
		return err
	}
	value2ColNum, err := data.getValue2ColNum()
	if err != nil {
		return err
	}

	if categoryColNum >= len(data.inputData.Series) {
		return errors.New("invalid category column")
	}
	if timeColNum >= len(data.inputData.Series) {
		return errors.New("invalid time column")
	}
	if valueColNum >= len(data.inputData.Series) {
		return errors.New("invalid value column")
	}
	if value2ColNum >= len(data.inputData.Series) {
		return errors.New("invalid value column")
	}

	categorySeries := data.inputData.Series[categoryColNum]
	timeSeries := data.inputData.Series[timeColNum]
	valuesSeries := data.inputData.Series[valueColNum]
	values2Series := data.inputData.Series[valueColNum]

	if value2ColNum != -1 {
		values2Series = data.inputData.Series[value2ColNum]
	}

	for i := 0; i < categorySeries.NRows(); i++ {
		data.AddCategory(categorySeries.Value(i).(string))
	}

	for i := 0; i < timeSeries.NRows(); i++ {
		catIndex := data.getCategoryIndex(categorySeries.Value(i).(string))

		if catIndex >= 0 {
			t := timeSeries.Value(i).(string)
			data.addTimeRangeValue(t)
		}
	}

	seriesInit := dataframe.SeriesInit{}
	seriesInit.Size = len(data.TimeRange)
	seriesInit.Capacity = len(data.TimeRange)

	series2Init := dataframe.SeriesInit{}
	series2Init.Size = len(data.TimeRange)
	series2Init.Capacity = len(data.TimeRange)

	s1 := dataframe.NewSeriesFloat64(data.Categories[0], &seriesInit)
	data.seriesData = dataframe.NewDataFrame(s1)

	s2 := dataframe.NewSeriesFloat64(data.Categories[0], &series2Init)
	data.series2Data = dataframe.NewDataFrame(s2)

	for i := 1; i < len(data.Categories); i++ {
		s := dataframe.NewSeriesFloat64(data.Categories[i], &seriesInit)
		data.seriesData.AddSeries(s, nil)
		ss := dataframe.NewSeriesFloat64(data.Categories[i], &series2Init)
		data.series2Data.AddSeries(ss, nil)
	}

	for i := 0; i < valuesSeries.NRows(); i++ {
		catIndex := data.getCategoryIndex(categorySeries.Value(i).(string))
		if catIndex >= 0 {
			timeRangeIndex := data.getTimeRangeIndex(timeSeries.Value(i).(string))
			if timeRangeIndex >= 0 {
				data.seriesData.Series[catIndex].Update(timeRangeIndex, valuesSeries.Value(i))
				if value2ColNum != -1 {
					data.series2Data.Series[catIndex].Update(timeRangeIndex, values2Series.Value(i))
				}
			}
		}
	}

	for i := 0; i < len(data.Categories); i++ {
		for j := 0; j < len(data.TimeRange); j++ {
			value := data.GetValue(i, j)
			if value > data.MaxValue {
				data.MaxValue = value
			}
		}
	}

	if data.MaxValue == 0 {
		data.MaxValue = 100.0
	}

	return nil
}

func (data *Data) AddCategory(newCategory string) error {

	found := false

	for _, category := range data.Categories {
		if newCategory == category {
			found = true
		}
	}

	if !found {
		data.Categories = append(data.Categories, newCategory)
	}

	return nil
}

func (data *Data) GetMaxValue() float64 {

	return data.MaxValue
}

func (data *Data) getDefaultValue(category int, index int) float64 {

	for i := index; i >= 0; i-- {
		if data.seriesData.Series[category].Value(index) != nil {
			return data.seriesData.Series[category].Value(index).(float64)
		}
	}

	return 0.0
}

func (data *Data) GetValue(category int, index int) float64 {

	if data.seriesData.Series[category].Value(index) == nil {
		return data.getDefaultValue(category, index)
	}

	return data.seriesData.Series[category].Value(index).(float64)
}

func (data *Data) GetValue2(category int, index int) float64 {

	return data.series2Data.Series[category].Value(index).(float64)
}

func (data *Data) addTimeRangeValue(newTime string) error {

	found := false

	for _, time := range data.TimeRange {
		if newTime == time {
			found = true
		}
	}

	if !found {
		if len(data.TimeRange) == 0 || (data.TimeRange[len(data.TimeRange)-1] < newTime) {
			data.TimeRange = append(data.TimeRange, newTime)
		} else {
			insertAtIndex := -1
			for i, value := range data.TimeRange {
				if value < newTime {
					insertAtIndex = i
				}
			}
			temp := append([]string{}, data.TimeRange[insertAtIndex+1:]...)
			data.TimeRange = append(data.TimeRange[0:insertAtIndex+1], newTime)
			data.TimeRange = append(data.TimeRange, temp...)
		}
	}
	return nil
}

func update(array []string, newItem string) []string {

	if len(array) == 0 || (array[len(array)-1] < newItem) {
		array = append(array, newItem)
	} else {
		insertAtIndex := -1
		for i, value := range array {
			if value < newItem {
				insertAtIndex = i
			}
		}
		temp := append([]string{}, array[insertAtIndex+1:]...)
		array = append(array[0:insertAtIndex+1], newItem)
		array = append(array, temp...)
	}
	return array
}

func (data *Data) getCategoryIndex(name string) int {

	for i, category := range data.Categories {
		if name == category {
			return i
		}
	}

	return -1
}

func (data *Data) getTimeRangeIndex(t string) int {

	for i, value := range data.TimeRange {
		if t == value {
			return i
		}
	}

	return -1
}

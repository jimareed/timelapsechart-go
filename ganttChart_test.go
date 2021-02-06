package timelapsechart

import (
	"bytes"
	"strings"
	"testing"
)

func TestGanttChart(t *testing.T) {

	csvStr := `
task,date,start,duration
Task 1,2021-01-04,0.0,30.0
Task 2,2021-01-04,20.0,40.0
Task 3,2021-01-04,50.0,50.0
`

	t.Log("a developer")

	data := Data{}
	err := data.ReadCSV(strings.NewReader(csvStr), []string{"task", "date", "duration", "start"})
	if err != nil {
		t.Fatal(" should be able to read the csv.", xMark, err)
	}

	config := Config{}
	config.Width = 400
	config.ChartX = 200

	chart := New("", "gantt", config)
	chart.AddData(&data)

	buffer := bytes.NewBuffer([]byte{})

	ganttChart := GanttChart{}

	err = ganttChart.render(chart, buffer)
	if err == nil {
		t.Log(" should be able to render a gantt chart.", checkMark)
	} else {
		t.Fatal(" should be able to render a gantt chart.", xMark, err)
	}

	width := ganttChart.RectWidth(chart, 50.0)
	if width == 50.0 {
		t.Log(" which should have the correct width.", checkMark)
	} else {
		t.Fatal(" which should have the correct width.", xMark, width)
	}

}

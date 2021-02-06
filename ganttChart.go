package timelapsechart

import (
	"bytes"
	"fmt"
)

// GanttChart
type GanttChart struct {
	maxValue float64
}

func (ganttChart *GanttChart) getMaxValue(chart *Chart) float64 {
	maxValue := 0.0

	for i, _ := range chart.Data.Categories {
		for j := 0; j < len(chart.Data.TimeRange); j++ {
			value := chart.Data.GetValue(i, j)
			value2 := chart.Data.GetValue2(i, j)
			if (value + value2) > maxValue {
				maxValue = value + value2
			}
		}
	}

	return maxValue
}

func (ganttChart *GanttChart) render(chart *Chart, buffer *bytes.Buffer) error {

	ganttChart.maxValue = ganttChart.getMaxValue(chart)

	for i, category := range chart.Data.Categories {

		fmt.Fprintf(buffer, `	<text x="%d" y="%d" fill="black" text-anchor="end" font-size="%dpx">%s</text>`,
			chart.Config.ChartX-10, chart.Config.ChartY+25+i*40, chart.Config.LabelSize, category)
		fmt.Fprintf(buffer, "\n")

		lastValue := 0.0
		begin := ""
		for j := 0; j < len(chart.Data.TimeRange); j++ {
			if j == 0 { // add rect based on first time interval
				value2 := chart.Data.GetValue2(i, len(chart.Data.TimeRange)-1)
				fmt.Fprintf(buffer, `	<rect x="%.02f" y="%d" fill="%s" width="1" height="40">`,
					float64(chart.Config.ChartX)+ganttChart.RectWidth(chart, value2), chart.Config.ChartY+i*40, chart.GetColor(i))
				fmt.Fprintf(buffer, "\n")
			}
			value := chart.Data.GetValue(i, j)
			fmt.Fprintf(buffer, `		<animate id="r%dstep%d" attributeName="width" from="%.02f" to="%.02f" %s dur="1.0s" fill="freeze" />`, i+1, j, ganttChart.RectWidth(chart, lastValue), ganttChart.RectWidth(chart, value), begin)
			fmt.Fprintf(buffer, "\n")
			lastValue = value
			begin = fmt.Sprintf(`begin="r%dstep%d.end"`, i+1, j)
		}

		fmt.Fprintf(buffer, "\t</rect>\n")
	}

	return nil
}

func (ganttChart *GanttChart) RectWidth(chart *Chart, value float64) float64 {
	maxWidth := float64(chart.Config.Width - (chart.Config.ChartX + 100))

	return maxWidth * (value / ganttChart.maxValue)
}

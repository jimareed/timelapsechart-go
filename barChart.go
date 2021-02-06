package timelapsechart

import (
	"bytes"
	"fmt"
)

// BarChart
type BarChart struct {
	reserved int
}

func (barChart *BarChart) render(chart *Chart, buffer *bytes.Buffer) error {

	for i, category := range chart.Data.Categories {
		fmt.Fprintf(buffer, `	<text x="%d" y="%d" fill="black" text-anchor="end" font-size="%dpx">%s</text>`,
			chart.Config.ChartX-10, chart.Config.ChartY+25+i*40, chart.Config.LabelSize, category)
		fmt.Fprintf(buffer, "\n")
		fmt.Fprintf(buffer, `	<rect x="%d" y="%d" fill="%s" width="1" height="40">`,
			chart.Config.ChartX, chart.Config.ChartY+i*40, chart.GetColor(i))
		fmt.Fprintf(buffer, "\n")

		lastValue := 0.0
		begin := ""
		for j := 0; j < len(chart.Data.TimeRange); j++ {
			value := chart.Data.GetValue(i, j)
			fmt.Fprintf(buffer, `		<animate id="r%dstep%d" attributeName="width" from="%.02f" to="%.02f" %s dur="1.0s" fill="freeze" />`, i+1, j, chart.RectWidth(lastValue), chart.RectWidth(value), begin)
			fmt.Fprintf(buffer, "\n")
			lastValue = value
			begin = fmt.Sprintf(`begin="r%dstep%d.end"`, i+1, j)
		}

		fmt.Fprintf(buffer, "\t</rect>\n")
	}

	return nil
}

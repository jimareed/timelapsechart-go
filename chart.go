package timelapsechart

import (
	"bytes"
	"fmt"
)

// Chart
type Chart struct {
	Config Config `json:"config"`
	Title  string `json:"title"`
	Data   *Data  `json:"data"`
	Type   string `json:"type"`
}

// Config
type Config struct {
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	LabelSize int      `json:"labelSize"`
	ChartX    int      `json:"chartX"`
	ChartY    int      `json:"chartY"`
	Palette   []string `json:"palette"`
}

func New(title string, chartType string, config Config) *Chart {
	chart := Chart{}

	chart.Title = title
	chart.Type = chartType

	if config.Width == 0 {
		config.Width = 800
	}
	if config.Height == 0 {
		config.Height = 600
	}
	if config.LabelSize == 0 {
		config.LabelSize = 20
	}
	if config.ChartX == 0 {
		config.ChartX = 200
	}
	if config.ChartY == 0 {
		config.ChartY = 50
	}
	if len(config.Palette) == 0 {
		config.Palette = append(config.Palette, "#f00")
		config.Palette = append(config.Palette, "#f70")
		config.Palette = append(config.Palette, "#ec0")
	}

	chart.Config = config

	return &chart
}

func (chart *Chart) GetColor(index int) string {

	index = index % len(chart.Config.Palette)

	return chart.Config.Palette[index]
}

func (chart *Chart) AddData(data *Data) error {

	chart.Data = data

	return nil
}

func (chart *Chart) Table() {

	fmt.Println("country,date,total")
	for i, category := range chart.Data.Categories {
		for j := 0; j < len(chart.Data.TimeRange); j++ {
			//			t, _ := time.Parse("2006-01-02", chart.Data.TimeRange[j])
			//			if int(t.Weekday()) == 1 {
			fmt.Printf("%s,%s,%.01f\n", category, chart.Data.TimeRange[j], chart.Data.GetValue(i, j))
			//			}
		}
	}
}

func (chart *Chart) Render(buffer *bytes.Buffer) error {

	header := `<svg version="1.1" xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">
 <g>
`
	body := ` 
<text x="400" y="35" fill="black" text-anchor="middle" font-size="20px">%s</text>
`
	footer := `
 </g>
</svg>
`

	_, err := fmt.Fprintf(buffer, header, chart.Config.Width, chart.Config.Height)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(buffer, body, chart.Title)
	if err != nil {
		return err
	}

	if chart.Type == "gantt" {
		err = chart.renderGanttChart(buffer)
	} else {
		err = chart.renderBarChart(buffer)
	}
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(buffer, footer)
	if err != nil {
		return err
	}

	return err
}

func (chart *Chart) RectWidth(value float64) float64 {
	maxWidth := float64(chart.Config.Width - (chart.Config.ChartX + 100))

	return maxWidth * (value / chart.Data.GetMaxValue())
}

func (chart *Chart) renderBarChart(buffer *bytes.Buffer) error {

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

func (chart *Chart) renderGanttChart(buffer *bytes.Buffer) error {

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

# timelapsechart-go

Go library to generate a time-lapse bar chart in SVG format.

<p  align="center">
    <img src="./example/output.svg" alt="timelapsechart-go output"/>
</p>

## usage

```golang
package main

import (
    "bytes"
    "fmt"
    "log"

    "github.com/jimareed/timelapsechart-go"
)

func main() {

    csv, err := os.Open("input.csv")
    if err != nil {
        log.Fatalln(err)
    }
    defer csv.Close()

    data := timelapsechart.Data{}
    err = data.ReadCSV(csv, []string{"country", "date", "total"})
    if err != nil {
        log.Fatalln("Couldn't open the csv file", err)
    }

    chart := timelapsechart.New("Covid Cases", timelapsechart.Config{})
    chart.AddData(&data)

    buffer := bytes.NewBuffer([]byte{})
    chart.Render(buffer)
    fmt.Println(buffer)
}
```

## Sources
- [dataframe-go](https://github.com/rocketlaunchr/dataframe-go)
- [covid-19-data](https://github.com/owid/covid-19-data)


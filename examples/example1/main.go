package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

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

	chart := timelapsechart.New("Covid Cases", "bar", timelapsechart.Config{})
	chart.AddData(&data)

	buffer := bytes.NewBuffer([]byte{})
	chart.Render(buffer)
	fmt.Println(buffer)
}

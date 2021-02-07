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
	err = data.ReadCSV(csv, []string{"task", "date", "duration", "start"})
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	chart := timelapsechart.New("Gantt Chart", "gantt", timelapsechart.Config{})
	chart.AddData(&data)
	chart.AddLabels([]string{"Q1", "Q2", "Q3", "Q4"})

	buffer := bytes.NewBuffer([]byte{})
	chart.Render(buffer)
	fmt.Println(buffer)
}

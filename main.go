package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

var (
	idx = map[string]int{
		"start":        1,
		"end":          2,
		"title":        3,
		"description":  4,
		"speakers":     5,
		"organization": 6,
	}
)

type event []string
type schedule [][]string

func main() {
	fmt.Println("\nReading CSV")
	csvFile, err := os.Open("schedule.csv")
	if err != nil {
		panic(err)
	}

	csvr := csv.NewReader(csvFile)
	csvr.Comma = ','
	csvr.Comment = '#'

	records, err := csvr.ReadAll()
	if err != nil {
		fmt.Printf("error with csv: %v\n", err)
	}

	f, err := os.Create("schedule.txt")
	if err != nil {
		panic(fmt.Sprintf("Can't create schedule.txt in the current folder: %e", err))
	}
	defer f.Close()

	table := tablewriter.NewWriter(f)
	table.SetHeader([]string{"Day", "Start", "End", "Title", "Speaker(s)", "Organization(s)"})
	table.SetRowLine(true)

	s := genSchedule(records)

	table.AppendBulk(s)
	table.Render()
}

func genSchedule(recs [][]string) schedule {
	var s schedule
	var d string
	for _, r := range recs {
		e := genEvent(r)
		if e[0] != d && d != "" {
			s = append(s, event{"----", "----", "----", "----", "----", "----"})
		}
		d = e[0]
		s = append(s, e)
	}

	return s
}

func genEvent(r []string) event {
	d, ts, te, _ := getTimeInfo(r[idx["start"]], r[idx["end"]])
	e := event{fmt.Sprintf("%s/%s", d[9:], d[5:7]), ts, te, r[idx["title"]], r[idx["speakers"]], r[idx["organization"]]}
	return e
}

func getTimeInfo(start, end string) (d, ts, te, dur string) {
	// Format: 11/5/15 11:25
	tFmt := "01/2/06 15:04"
	var s, e time.Time
	var err error
	if s, err = time.Parse(tFmt, start); err != nil {
		fmt.Printf("Could not parse start time: %v", err)
	}
	if e, err = time.Parse(tFmt, end); err != nil {
		fmt.Printf("Could not parse end time: %v", err)
	}

	dFmt := "2006-01-02"
	d = s.Format(dFmt)
	tFmt = "15:04"
	ts = s.Format(tFmt)
	te = e.Format(tFmt)
	diff := e.Sub(s)
	dur = fmt.Sprintf("00:%2.f", diff.Minutes())
	return d, ts, te, dur
}

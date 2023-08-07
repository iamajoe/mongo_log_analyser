package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func help() {
	log.Println(
		"Usage:\n./mongo_log_analyser <stats> [options...]\n\nCheck documentation for more information",
	)
}

func templateStatistics(logs map[string][]mongoLogAnalyse) (string, error) {
	fullCount := 0
	var slowest mongoLogAnalyse

	sep := "\n=====================================\n"
	tmpl := sep + "Collection count result:\n\n"

	// go per collection count
	for c, v := range logs {
		// calculate the total
		fullCount += len(v)

		// get collection statistics
		tmpl += fmt.Sprintf("- %s : %d \n", c, len(v))

		// find the lowest
		for _, l := range v {
			if slowest.duration < l.duration {
				slowest = l
			}
		}
	}

	tmpl += sep + "\n"
	tmpl += fmt.Sprintf("slowest: \n- %dms\n- %s\n", slowest.duration, slowest.raw)
	tmpl += fmt.Sprintf("total: %d\n", fullCount)
	tmpl += sep

	return tmpl, nil
}

func main() {
	statsFs := flag.NewFlagSet("stats", flag.ExitOnError)
	statsInput := statsFs.String("i", "", "input with mongo logs")
	statsOutput := statsFs.String("o", "", "output with results")
	statsFilter := statsFs.String("f", "[]", "filters out an array of patterns")
	statsMinDurationMs := statsFs.Int("t", 1000, "minimum request duration time in ms")
	statsHelp := statsFs.Bool("h", false, "help manual")

	if len(os.Args) < 2 {
		help()
		return
	}

	switch os.Args[1] {

	case "stats":
		if err := statsFs.Parse(os.Args[2:]); err != nil {
			statsFs.PrintDefaults()
			log.Fatal(err)
		}

		if *statsHelp {
			statsFs.PrintDefaults()
			return
		}

		// parse the filter
		filter := []string{}
		if statsFilter != nil && len(*statsFilter) > 0 {
			err := json.Unmarshal([]byte(*statsFilter), &filter)
			if err != nil {
				log.Fatal(err)
			}
		}

		res, err := stats(
			*statsInput,
			*statsMinDurationMs,
			filter,
		)
		if err != nil {
			log.Fatal(err)
		}

		output := *statsOutput
		if len(output) > 0 {
			// TODO: save to file but what?
		} else {
			tmpl, err := templateStatistics(res)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(tmpl)
		}
		break
	default:
		help()
	}
}

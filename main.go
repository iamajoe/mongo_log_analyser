package main

import (
	"flag"
	"log"
	"os"
)

func help() {
	log.Println(
		"Usage:\n./mongo_log_analyser <stats> [options...]\n\nCheck documentation for more information",
	)
}

func main() {
	statsFs := flag.NewFlagSet("stats", flag.ExitOnError)
	statsInputRaw := statsFs.String("i", "", "input with mongo logs")
	statsOutputRaw := statsFs.String("o", "tmp_run_result.txt", "output with results")
	statsHelpRaw := statsFs.Bool("h", false, "help manual")

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

		if *statsHelpRaw {
			statsFs.PrintDefaults()
			return
		}

		if err := stats(
			*statsInputRaw,
			*statsOutputRaw,
		); err != nil {
			log.Fatal(err)
		}
		break
	default:
		help()
	}
}

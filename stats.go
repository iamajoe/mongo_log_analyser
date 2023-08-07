package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type mongoLogCommand struct {
	GetMore    uint64                 `json:"geMore"`
	Find       string                 `json:"find"`
	Collection string                 `json:"collection"`
	BatchSize  uint64                 `json:"batchSize"`
	Projection map[string]int         `json:"projection"`
	Lsid       map[string]interface{} `json:"lsid"`
	Db         string                 `json:"$db"`
	Filter     map[string]interface{} `json:"filter"`
}

// TODO: we should analyse errors too

// REF: https://www.mongodb.com/docs/manual/reference/log-messages/
type mongoLog struct {
	T struct {
		Date string `json:"$date"`
	} `json:"t"`
	S    string `json:"s"`
	C    string `json:"c"`
	Id   uint64 `json:"id"`
	Ctx  string `json:"ctx"`
	Msg  string `json:"msg"`
	Attr *struct {
		T                  string           `json:"type"`
		Ns                 string           `json:"ns"`
		Command            *mongoLogCommand `json:"command"`
		OriginatingCommand *mongoLogCommand `json:"originatingCommand"`
		PlanSummary        string           `json:"planSummary"`
		Cursorid           uint64           `json:"cursorid"`
		KeysExamined       uint64           `json:"keysExamined"`
		DocsExamined       uint64           `json:"docsExamined"`
		NumYields          uint64           `json:"numYields"`
		Nreturned          uint64           `json:"nreturned"`
		Reslen             uint64           `json:"reslen"`
		Locks              *struct{}        `json:"locks"`
		Storage            *struct{}        `json:"storage"`
		Protocol           string           `json:"protocol"`
		DurationMillis     uint64           `json:"durationMillis"`
		Message            string           `json:"message"`
	} `json:"attr"`
	Tags []string `json:"tags"`
	Size uint64   `json:"size"`
}

type mongoLogAnalyse struct {
	collection string
	duration   uint64
	raw        []byte
	log        mongoLog
	command    *mongoLogCommand
}

func analyseRequest(a mongoLogAnalyse) (mongoLogAnalyse, error) {
	l := a.log

	if l.Attr == nil {
		return a, errors.New(fmt.Sprintf("no Attribute found, id: %d , date: %s", l.Id, l.T.Date))
	}

	if l.Attr.OriginatingCommand != nil {
		a.command = l.Attr.OriginatingCommand
	} else if l.Attr.Command != nil {
		a.command = l.Attr.Command
	}

	if a.command == nil {
		return a, errors.New(fmt.Sprintf("no command found, id: %d , date: %s", l.Id, l.T.Date))
	}

	a.duration = l.Attr.DurationMillis

	// check the collection on the original command first
	if len(a.command.Collection) > 0 {
		a.collection = a.command.Collection
	} else if len(a.command.Find) > 0 {
		a.collection = a.command.Find
	}

	return a, nil
}

func isFiltered(l mongoLogAnalyse, ignorePatterns []string) bool {
	for _, p := range ignorePatterns {
		p = strings.ReplaceAll(p, " ", "")
		if len(p) == 0 {
			continue
		}

		// now check the url
		r, err := regexp.Compile(p)
		if err != nil {
			continue
		}

		i := r.FindIndex(l.raw)
		if len(i) > 0 {
			return true
		}
	}

	return false
}

func analyse(raw []byte, ignorePatterns []string) (mongoLogAnalyse, error) {
	a := mongoLogAnalyse{
		raw: raw,
	}

	// parse the log into the right structure
	if err := json.Unmarshal(raw, &a.log); err != nil {
		return a, err
	}

	// filtered out so, dont continue analyse
	if isFiltered(a, ignorePatterns) {
		return a, nil
	}

	switch a.log.C {
	case "COMMAND":
		return analyseRequest(a)
	}

	return a, nil
}

// stats fetches from the input the statistics
func stats(inputPath string, minDurationMs int, ignorePattern []string) (map[string][]mongoLogAnalyse, error) {
	statistics := make(map[string][]mongoLogAnalyse)

	if len(inputPath) == 0 {
		return statistics, errors.New("input path is required")
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return statistics, err
	}
	defer file.Close()

	// for statistics, run a scanner line by line on the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v := scanner.Text()
		if err != nil {
			return statistics, err
		}

		if len(v) == 0 {
			continue
		}

		l, err := analyse([]byte(v), ignorePattern)
		if err != nil || len(l.collection) == 0 {
			continue
		}

		arr, ok := statistics[l.collection]
		if !ok {
			arr = []mongoLogAnalyse{}
		}

		statistics[l.collection] = append(arr, l)
	}

	if err := scanner.Err(); err != nil {
		return statistics, err
	}

	return statistics, nil
}

package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type mongoLog struct {
	t struct {
		date string `json:"$date"`
	}
	s    string
	c    string
	id   uint
	ctx  string
	msg  string
	attr struct {
		t       string `json:"type"`
		ns      string
		command struct {
			getMore    uint64
			find       string
			collection string
			batchSize  uint64
			projection map[string]int
			lsid       map[string]interface{}
			db         string `json:"$db"`
			filter     map[string]interface{}
		}
		originatingCommand struct{}
		planSummary        string
		cursorid           uint64
		keysExamined       uint64
		docsExamined       uint64
		numYields          uint64
		nreturned          uint64
		reslen             uint64
		locks              struct{}
		storage            struct{}
		protocol           string
		durationMillis     uint64
	}
}

// stats fetches from the input the statistics
func stats(inputPath string, outputPath string) error {
	if len(inputPath) == 0 {
		return errors.New("input path is required")
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// for statistics, run a scanner line by line on the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_ = strings.ToLower(scanner.Text())
		if err != nil {
			return err
		}

		// TODO: unmarshall v, this is a mongo log
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// TODO: need to inform

	return nil
}

package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func csvParse(jacocoPath string, blackList BlackList) (map[string]Target, error) {
	jacoco, err := os.Open(jacocoPath)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(jacoco)
	targets := map[string]Target{}
	// skip first line
	r.Read()
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if check(blackList, record) == false {
			continue
		}

		target, err := convertTarget(record)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := targets[target.String()]; !ok {
			targets[target.String()] = target
		}
	}
	return targets, nil
}

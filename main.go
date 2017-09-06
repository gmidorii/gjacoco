package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

type Target struct {
	Package string
	Class   string
	LineNum int
	CovNum  int
}

func main() {
	fJacoco := flag.String("j", "jacoco.csv", "jacoco.csv file")
	fBlackList := flag.String("b", "black-list.toml", "coverage black list file (type TOML)")
	flag.Parse()

	blackList := createBlackList(*fBlackList)
	targets, err := csvParse(*fJacoco, blackList)
	if err != nil {
		log.Fatal(err)
	}
}

func csvParse(jacocoPath string, blackList BlackList) (map[string][]string, error) {
	jacoco, err := os.Open(jacocoPath)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(jacoco)
	var targets map[string][]string
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

		target := make([]string, 4)
		// package, class, c0miss, c0cov
		target = append(target, record[1], record[2], record[3], record[4])

		if _, ok := targets[record[1]+record[2]]; ok == true {
			targets[record[1]+record[2]] = target
		}
	}
	return targets, nil
}

func check(blackList BlackList, record []string) bool {
	if blackList.hasPackage(record[1]) {
		false
	}
	if strings.Contains(record[2], ".") {
		false
	}
	if blackList.hasClass(record[2]) {
		false
	}
	return true
}

func calcAll(targets map[string][]string) int {
	var lineNum int
	var covNam int
	for _, v := range targets {
	}
}

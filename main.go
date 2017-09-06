package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Target struct {
	Package string
	Class   string
	LineNum int
	CovNum  int
}

func (t Target) String() string {
	return fmt.Sprintf("%s.%s", t.Package, t.Class)
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

	fmt.Println(calcAllCov(targets))
}

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

func check(blackList BlackList, record []string) bool {
	if blackList.hasPackage(record[1]) {
		return false
	}
	if strings.Contains(record[2], ".") {
		return false
	}
	if blackList.hasClass(record[2]) {
		return false
	}
	return true
}

func convertTarget(record []string) (Target, error) {
	missNum, err := strconv.Atoi(record[3])
	if err != nil {
		return Target{}, err
	}

	covNum, err := strconv.Atoi(record[4])
	if err != nil {
		return Target{}, err
	}

	return Target{
		Package: record[1],
		Class:   record[2],
		LineNum: missNum + covNum,
		CovNum:  covNum,
	}, nil
}

func calcAllCov(targets map[string]Target) string {
	var sumLineNum int
	var sumCovNum int
	for _, v := range targets {
		sumLineNum += v.LineNum
		sumCovNum += v.CovNum
	}

	cov := float64(sumCovNum) / float64(sumLineNum) * 100
	return fmt.Sprint(cov)[0:5]
}

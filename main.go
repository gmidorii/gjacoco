package main

import (
	"flag"
	"fmt"
	"log"
)

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
	if err := outputTarget("./template/output.html", targets); err != nil {
		log.Println(err)
	}
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

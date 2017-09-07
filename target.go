package main

import (
	"fmt"
	"html/template"
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

type Output struct {
	CovAll  string
	Targets []Target
}

func (t Target) String() string {
	return fmt.Sprintf("%s.%s", t.Package, t.Class)
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

func outputTarget(fpath string, targets map[string]Target) error {
	t := template.Must(template.ParseFiles("./template/template.html"))

	fout, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer fout.Close()

	covAll := calcAllCov(targets)
	targetSlice := make([]Target, 0)
	for _, v := range targets {
		targetSlice = append(targetSlice, v)
	}

	output := Output{
		CovAll:  covAll,
		Targets: targetSlice,
	}

	return t.Execute(fout, &output)
}

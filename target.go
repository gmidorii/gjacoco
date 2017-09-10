package main

import (
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Target struct {
	Package  string
	Class    string
	Coverage float64
	LineNum  int
	CovNum   int
}

type Output struct {
	Statistic Statistic
	Targets   []Target
}

type Statistic struct {
	CovAll   string
	CovRatio map[int]int
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
	cov := float64(covNum) / float64(covNum+missNum) * 100

	return Target{
		Package:  record[1],
		Class:    record[2],
		Coverage: cov,
		LineNum:  missNum + covNum,
		CovNum:   covNum,
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

	sort.Slice(targetSlice, func(i, j int) bool {
		return targetSlice[i].Coverage < targetSlice[j].Coverage
	})

	output := Output{
		Statistic: Statistic{
			CovAll:   covAll,
			CovRatio: makeCovRatio(targetSlice),
		},
		Targets: targetSlice,
	}

	return t.Execute(fout, &output)
}

func makeCovRatio(targets []Target) map[int]int {
	var covRationMap = make(map[int]int, 5)
	targetLen := len(targets)
	count := 0

	// 0
	index := sort.Search(targetLen, func(i int) bool {
		return targets[i].Coverage > 0
	})
	covRationMap[0] = index
	count += index

	// 0 < cov < 30
	index = sort.Search(targetLen, func(i int) bool {
		return targets[i].Coverage >= 30
	})
	covRationMap[30] = index - count
	count += covRationMap[30]

	// 30 <= cov < 60
	index = sort.Search(targetLen, func(i int) bool {
		return targets[i].Coverage > 60
	})
	covRationMap[60] = index - count
	count += covRationMap[60]

	// 60 <= cov < 90
	index = sort.Search(targetLen, func(i int) bool {
		return targets[i].Coverage > 90
	})
	covRationMap[90] = index - count
	count += covRationMap[90]

	// 90 <= cov <= 100
	covRationMap[100] = targetLen - count

	return covRationMap
}

package main

import (
	"reflect"
	"testing"
)

var ratioTests = []struct {
	in  []Target
	out map[int]int
}{
	{[]Target{
		Target{Coverage: 0},
	},
		map[int]int{0: 1, 30: 0, 60: 0, 90: 0, 100: 0},
	},
	{[]Target{
		Target{Coverage: 29},
	},
		map[int]int{0: 0, 30: 1, 60: 0, 90: 0, 100: 0},
	},
	{[]Target{
		Target{Coverage: 59},
	},
		map[int]int{0: 0, 30: 0, 60: 1, 90: 0, 100: 0},
	},
	{[]Target{
		Target{Coverage: 89},
	},
		map[int]int{0: 0, 30: 0, 60: 0, 90: 1, 100: 0},
	},
	{[]Target{
		Target{Coverage: 100},
	},
		map[int]int{0: 0, 30: 0, 60: 0, 90: 0, 100: 1},
	},
}

func TestMakeCovRatio(t *testing.T) {
	for _, v := range ratioTests {
		r := makeCovRatio(v.in)
		if eq := reflect.DeepEqual(v.out, r); !eq {
			t.Error("not equal")
			t.Logf("expected:%v, actual:%v", v.out, r)
		}
	}
}

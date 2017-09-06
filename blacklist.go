package main

import (
	"strings"

	"github.com/BurntSushi/toml"
)

// BlackList is black list setting is not measuring coverage
type BlackList struct {
	Packages []string
	Classes  []string
}

func (b *BlackList) hasPackage(pack string) bool {
	for _, v := range b.Packages {
		if strings.Contains(pack, v) {
			return true
		}
	}
	return false
}

func (b *BlackList) hasClass(class string) bool {
	for _, v := range b.Classes {
		if class == v {
			return true
		}
	}
	return false
}

func createBlackList(fpath string) BlackList {
	var blackList BlackList
	_, err := toml.DecodeFile(fpath, &blackList)
	if err != nil {
		return BlackList{}
	}
	return blackList
}

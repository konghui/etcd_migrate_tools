package glob

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type GlobalParttern struct {
	BlackList []string `json: "blacklist"`
	WhiteList []string `json: whitelist`
}

func NewGlobalParttern(file string) (*GlobalParttern, error) {
	var globalParttern GlobalParttern
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &globalParttern)
	if err != nil {
		return nil, err
	}

	log.Printf("BlackList=%s\nWhiteList=%s\n", globalParttern.BlackList, globalParttern.WhiteList)
	return &globalParttern, err
}

func (this *GlobalParttern) IsInBlackList(subject string) bool {
	for _, blackItem := range this.BlackList {
		if Glob(blackItem, subject) == true {
			return true
		}
	}
	return false
}

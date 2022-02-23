package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

const (
	startURL = "https://api.hh.ru/vacancies"
	dataFile = "data.json"
)

var (
	dataFileLock sync.RWMutex
)

func PrepareRequestURL(receivedURL string) string {
	parts := strings.Split(receivedURL, "?")
	return startURL + "?" + parts[1]
}

func Check(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
func Pretty(data interface{}) ([]byte, error) {
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	return file, nil
}
func ChangeFile(data interface{}, fileName string) error {
	switch fileName {
	case dataFile:
		dataFileLock.RLock()
		defer dataFileLock.RUnlock()
	default:
		return fmt.Errorf("unknown file name")
	}
	prettyData, err := Pretty(data)
	Check(err)
	err = ioutil.WriteFile(fileName, prettyData, 0644)
	if err != nil {
		return err
	}
	return nil
}
func ReadFile(fileName string) ([]byte, error) {
	switch fileName {
	case dataFile:
		dataFileLock.RLock()
		defer dataFileLock.RUnlock()
	default:
		return nil, fmt.Errorf("unknown file name")
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

package util

import (
	"encoding/json"
	"io/ioutil"
)

func ParseJsonFile(path string, object interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, object)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonFile(path string, object interface{}) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

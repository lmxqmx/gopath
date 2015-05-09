package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

//load
func LoadJSON(j interface{}, t string, r map[string]string) error {
	t = replace(t, r)

	//decode string(JSON format)
	dec := json.NewDecoder(strings.NewReader(t))
	err := dec.Decode(j)
	if err != nil {
		return err
	}

	return nil
}

//load config file
func LoadJSONFile(j interface{}, path string, r map[string]string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return LoadJSON(j, string(b), r)
}

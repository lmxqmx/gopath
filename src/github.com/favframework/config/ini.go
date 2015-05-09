package config

import (
	"log"
	"io/ioutil"
)

//load
func LoadINI(j interface{}, t string, r map[string]string) error {
	t = replace(t, r)

	//FIXME: decode INI
	log.Fatal("FIXME")
	return nil
}

//load config file
func LoadINIFile(j interface{}, path string, r map[string]string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return LoadINI(j, string(b), r)
}

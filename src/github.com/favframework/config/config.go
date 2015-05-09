package config

import (
	"log"
	"os"
	"strings"
	"runtime"
)

//Replace by map
func replace(t string, replaces map[string]string) string {
	if len(replaces) <= 0 {
		return t
	}

	for k, v := range replaces {
		//t = regexp.MustCompile(k).ReplaceAllLiteralString(t, v)
		t = strings.Replace(t, k, v, -1)
	}

	return t
}

//Get working directory.
func WorkingDir() string {
	w, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//replace backslash(\) to slash(/) on windows platform
	if runtime.GOOS == "windows" {
		w = strings.Replace(w, "\\", "/", -1)
	}

	return w
}

package util

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// var DEFAULT_MODE os.FileMode = os.ModePerm

func Fexists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FTouch(path string) error {
	return FTouch2(path, os.ModePerm)
}
func FTouch2(path string, fm os.FileMode) error {
	f, err := os.Open(path)
	if err != nil {
		p := filepath.Dir(path)
		if !Fexists(p) {
			err := os.MkdirAll(p, fm)
			if err != nil {
				return err
			}
		}
		f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fm)
		if f != nil {
			defer f.Close()
		}
		return err
	}
	defer f.Close()
	fi, _ := f.Stat()
	if fi.IsDir() {
		return errors.New("can't touch path")
	}
	return nil
}
func FWrite(path, data string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	return err
}
func FAppend(path, data string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	return err
}
func FCopy(src string, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}
func ReadLine(r *bufio.Reader, limit int, end bool) ([]byte, error) {
	var isPrefix bool = true
	var bys []byte
	var tmp []byte
	var err error
	for isPrefix {
		tmp, isPrefix, err = r.ReadLine()
		if err != nil {
			return nil, err
		}
		bys = append(bys, tmp...)
	}
	if end {
		bys = append(bys, '\n')
	}
	return bys, nil
}

func Timestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}
func Time(timestamp int64) time.Time {
	return time.Unix(0, timestamp*1e6)
}
func AryExist(ary interface{}, obj interface{}) bool {
	switch reflect.TypeOf(ary).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(ary)
		for i := 0; i < s.Len(); i++ {
			if obj == s.Index(i).Interface() {
				return true
			}
		}
		return false
	default:
		return false
	}
}

var C_SH string = "/bin/bash"

func Exec(args ...string) (string, error) {
	bys, err := exec.Command(C_SH, "-c", strings.Join(args, " ")).Output()
	return string(bys), err
}

func IsType(v interface{}, t string) bool {
	t = strings.Trim(t, " \t")
	if v == nil || len(t) < 1 {
		return false
	}
	return reflect.Indirect(reflect.ValueOf(v)).Type().Name() == t
}

func Append(ary []interface{}, args ...interface{}) []interface{} {
	for _, arg := range args {
		ary = append(ary, arg)
	}
	return ary
}

func List(root string, reg string) []string {
	return ListFunc(root, reg, func(t string) string {
		return t
	})
}
func ListFunc(root string, reg string, f func(t string) string) []string {
	pathes := []string{}
	regx := regexp.MustCompile(reg)
	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if regx.MatchString(path) {
			pathes = append(pathes, f(path))
		}
		return nil
	})
	return pathes
}

// Copyright 2014 The fav Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"log"
	"os"
)

//env: 0:development 1:testing 2:staging 3:production
var env int8

//output to screen
var logStderr *log.Logger = log.New(os.Stderr, "", log.Lshortfile)

//write to file
var logFile *log.Logger

func Init(env int8){
	//var Log *log.Logger
	//Log = log.New(os.Stderr, "", flag)

	flag := log.Lshortfile
	logStderr = log.New(os.Stderr, "", flag)

/*
	switch env {
	case 0:
		log.SetOutput(os.Stderr)
		log.SetPrefix("")
		log.SetFlags(log.Lshortfile)
	case 1:
		log.SetOutput(os.Stderr)
		log.SetPrefix("")
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	case 2, 3:
		log.SetOutput(os.Stderr)
		log.SetPrefix("")
		log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	}
*/
}

func Print(v ...interface{}) {
	logStderr.Print(v...)
}

func Printf(format string, v ...interface{}) {
	logStderr.Printf(format, v...)
}

func Println(v ...interface{}) {
	logStderr.Println(v...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	logStderr.Fatal(v...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	logStderr.Fatalf(format, v...)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	logStderr.Fatalln(v...)
}

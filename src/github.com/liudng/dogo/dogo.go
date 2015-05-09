// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/favframework/console"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//Dogo object
type Dogo struct {
	//source files
	SourceDir []string

	//file extends
	SourceExt string

	//Working Dir
	WorkingDir string

	//build command
	BuildCmd string

	//run command
	RunCmd string

	//file list
	Files map[string]time.Time

	//Cmd object
	cmd *exec.Cmd

	//file modified
	isModified bool

	//build error
	buildErr string

	//build retry
	retries int64
}

//start new monitor
func (d *Dogo) NewMonitor() {
	//fmt.Printf("%#v\n", d.SourceDir)

	if d.WorkingDir == "" {
		//log.Fatalf("[dogo] dogo.json (BuildCmd) error. \n")
		d.WorkingDir = WorkingDir
	}
	if len(d.SourceDir) == 0 {
		//log.Fatalf("[dogo] dogo.json (SourceDir) error. \n")
		d.SourceDir = append(d.SourceDir, WorkingDir)
	}
	if d.SourceExt == "" {
		//log.Fatalf("[dogo] dogo.json (SourceExt) error. \n")
		d.SourceExt = ".go|.c|.cpp|.h"
	}
	if d.BuildCmd == "" {
		//log.Fatalf("[dogo] dogo.json (BuildCmd) error. \n")
		d.BuildCmd = "go build ."
	}
	if d.RunCmd == "" {
		//log.Fatalf("[dogo] dogo.json (RunCmd) error. \n")
		d.RunCmd = filepath.Base(WorkingDir)
		if runtime.GOOS == "windows" {
			d.RunCmd += ".exe"
		}
	}

	console.Chdir(d.WorkingDir)

	fmt.Printf("[dogo] Working Directory:\n")
	fmt.Printf("       %s\n", d.WorkingDir)

	fmt.Printf("[dogo] Monitoring Directories:\n")
	for _, dir := range d.SourceDir {
		fmt.Printf("       %s\n", dir)
	}

	fmt.Printf("[dogo] File extends:\n")
	fmt.Printf("       %s\n", d.SourceExt)

	fmt.Printf("[dogo] Build command:\n")
	fmt.Printf("       %s\n", d.BuildCmd)

	fmt.Printf("[dogo] Run command:\n")
	fmt.Printf("       %s\n", d.RunCmd)

	d.Files = make(map[string]time.Time)

	d.InitFiles()

	//FIXME: add console support.

	//FIXME: moniting directories: add file, delete file.

	//FIXME: Multi commands.
}

func (d *Dogo) InitFiles() {
	extends := strings.Split(d.SourceExt, "|")

	//scan source directories
	for _, dir := range d.SourceDir {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				d.FmtPrintf("%s\n", err)
				return err
			}

			for _, ext := range extends {
				if filepath.Ext(path) == ext {
					//fmt.Println(path)
					d.Files[path] = f.ModTime()
					break
				}
			}

			return nil
		})
	}
}

func (d *Dogo) Monitor() {
	for {
		d.Compare()

		if d.isModified == true {
			d.BuildAndRun()
		}

		time.Sleep(time.Duration(1 * time.Second))
	}
}

//compare source file's modify time
func (d *Dogo) Compare() {
	changed := false

	for p, t := range d.Files {
		info, err := os.Stat(p)
		if err != nil {
			d.FmtPrintf("%s\n", err)
			continue
		}

		//new modtime
		nt := info.ModTime()

		if nt.Sub(t) > 0 {
			d.Files[p] = nt
			changed = true
			d.FmtPrintf("[dogo] Changed files: %s\n", filepath.Base(p))
		}
	}

	if changed == true {
		d.isModified = true
	} else {
		d.isModified = false
	}
}

func (d *Dogo) BuildAndRun() {
	if d.cmd != nil {
		d.FmtPrintf("[dogo] Terminate the process %d: ", d.cmd.Process.Pid)
		if err := d.cmd.Process.Kill(); err != nil {
			d.FmtPrintf("\n%s\n", err)
		} else {
			d.FmtPrintf("success.\n")
		}
	}

	if err := d.Build(); err != nil {
		d.FmtPrintf("[dogo] Build failed: %s\n\n", err)
	} else {
		//run program
		d.FmtPrintf("[dogo] Start the process: %s\n\n", d.RunCmd)
		go d.Run()
	}
}

//build
func (d *Dogo) Build() error {
	d.FmtPrintf("[dogo] Start build: ")
	args := console.ParseText(d.BuildCmd)
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		fullOut := string(out)
		if d.buildErr == "" || d.buildErr != fullOut {
			d.FmtPrintf("\n%s", fullOut)
			d.retries = 0
			d.buildErr = fullOut
		} else {
			//d.FmtPrintf(".")
			d.retries++
		}
		return err
	} else {
		d.retries = 0
		d.buildErr = ""
		d.FmtPrintf("success.\n")
		return nil
	}
}

//run it
func (d *Dogo) Run() {
	args := console.ParseText(d.RunCmd)
	d.cmd = exec.Command(args[0], args[1:]...)
	d.cmd.Stdin = os.Stdin
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr
	err := d.cmd.Run()
	if err != nil {
		//fmt.Printf("%s\n", err)
	} else {
		d.cmd = nil
		//fmt.Printf("exit status 0.\n")
	}
}

func (d *Dogo) LogPrintf(format string, v ...interface{}) {
	if d.retries == 0 {
		log.Printf(format, v...)
	}
}

func (d *Dogo) FmtPrintf(format string, v ...interface{}) {
	if d.retries == 0 {
		fmt.Printf(format, v...)
	}
}

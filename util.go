package main

import (
	"errors"
	"log"
	"os/exec"
	"bytes"
	"time"
	"runtime"
	"strings"
)

func newErr(input string, err error) error {
	err = errors.New(input + " " + err.Error())
	log.Print(err)
	return err
}

func handle(input string, err error) error {
	pc, fn, line, _ := runtime.Caller(1)
	if input[len(input)-1:] != "." { // Add a period.
		input += "."
	}
	input += " " + err.Error()
	p := strings.Split(fn, "/")
	fn = p[len(p)-1]
	log.Printf("[error] in %s[%s:%d] %s", runtime.FuncForPC(pc).Name(), fn, line, input)
	return errors.New(input)
}

// Outputs a system command to log with all output on error.
func commandVerboseOnErr(input string) (string, string, error) {
	stdout, stderr, err := commandWithOutput(input)
	if err != nil {
		log.Print("Command: " + input)
		if stdout != "" {
			log.Print(stdout)
		}
		if stderr != "" {
			log.Print(stderr)
		}
		err = newErr("Error in running command.", err)
		log.Print(err)
	}
	return stdout, stderr, err
}

// Outputs a system command to log with stdout, stderr, and err output.
func commandVerbose(input string) (string, string, error) {
	log.Print("Command: " + input)
	stdout, stderr, err := commandWithOutput(input)
	if stdout != "" {
		log.Print(stdout)
	}
	if stderr != "" {
		log.Print(stderr)
	}
	if err != nil {
		err = newErr("Error in running command.", err)
		log.Print(err)
	} else {
		log.Print("Command ran with no errors.")
	}
	return stdout, stderr, err
}

// Executes a shell command and returns the stdout, stderr, and err
func commandWithOutput(input string) (string, string, error) {
	cmd := exec.Command("sh", "-cx", input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outResp := stdout.String()
	errResp := stderr.String()
	return outResp, errResp, err
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
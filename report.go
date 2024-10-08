package main

import (
	"fmt"
	"os"
)

type Report struct {
	stdout       string
	stderr       string
	duration     float32
	processState *os.ProcessState
}

func (r *Report) Print() {
	fmt.Println("Run report")
	fmt.Print(r.stdout)
	fmt.Println("---")
	fmt.Print(r.stderr)
	fmt.Println("---")
	if r.processState.ExitCode() >= 0 {
		fmt.Println("Exit Code:", r.processState.ExitCode())
	} else {
		fmt.Println("Process was killed by signal")
	}
	fmt.Println(r.processState.ExitCode())
	fmt.Println("Duration", r.duration)
}

func readFile(path string) string {
	body, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func NewReport(procState *os.ProcessState, files *ProcFiles, duration int64) *Report {
	return &Report{
		stdout:       readFile(files.stdout.Name()),
		stderr:       readFile(files.stderr.Name()),
		duration:     float32(duration) / 1000,
		processState: procState,
	}
}

package main

import (
	"os"
	"sync"
	"syscall"
	"time"
)

type ProcFiles struct {
	stdin  *os.File
	stdout *os.File
	stderr *os.File
}

func (files *ProcFiles) getProcAttr() []*os.File {
	return []*os.File{files.stdin, files.stdout, files.stderr}
}

func (files *ProcFiles) cleanup() {
	err := os.Remove(files.stdout.Name())
	if err != nil {
		panic(err)
	}
	err = os.Remove(files.stderr.Name())
	if err != nil {
		panic(err)
	}
}

func createTemp(dir string, filename string) *os.File {
	file, err := os.CreateTemp(dir, filename)
	if err != nil {
		panic(err)
	}
	return file
}

func NewProcFiles(dir string) *ProcFiles {
	return &ProcFiles{
		stdin:  os.Stdin,
		stdout: createTemp(dir, "stdout"),
		stderr: createTemp(dir, "stderr"),
	}
}

func run(args *Args, files *ProcFiles) (*os.ProcessState, int64) {
	waitCommandGroup := new(sync.WaitGroup)

	chProcState := make(chan *os.ProcessState)
	defer close(chProcState)

	chDuration := make(chan int64)
	defer close(chDuration)

	//chPid := make(chan int)
	//defer close(chPid)

	chPid := make(chan int)
	defer close(chPid)

	chKillerQuit := make(chan bool)
	defer close(chKillerQuit)

	waitCommandGroup.Add(1)
	go runCommand(args.command, files, waitCommandGroup, chProcState, chDuration, chPid, chKillerQuit)
	go runKiller(args.timeout, chPid, chKillerQuit)
	waitCommandGroup.Wait()
	//proc := <-chProc
	procState :=
		<-chProcState
	duration := <-chDuration

	return procState, duration
}

func runCommand(command string, files *ProcFiles, waitGroup *sync.WaitGroup, chProcState chan<- *os.ProcessState, chDuration chan<- int64, chPid chan<- int, chKillerQuit chan<- bool) {
	//defer waitGroup.Done()
	var procAttr os.ProcAttr
	//procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	procAttr.Files = files.getProcAttr()
	procAttr.Sys = &syscall.SysProcAttr{Setpgid: true}
	cmdArgs := []string{
		"/usr/bin/env bash",
		"-c",
		command,
	}
	startTimeStamp := time.Now()
	proc, err := os.StartProcess("/bin/bash", cmdArgs, &procAttr)
	if err != nil {
		panic(err)
	}
	//chPid <- proc.Pid
	chPid <- proc.Pid

	pState, err := proc.Wait()
	if err != nil {
		panic(err)
	}
	waitGroup.Done()

	//return pState, time.Now().UnixMilli() - startTimeStamp.UnixMilli()
	if pState.Exited() {
		chKillerQuit <- true
	}
	//chProc <- proc
	chProcState <- pState
	chDuration <- time.Now().UnixMilli() - startTimeStamp.UnixMilli()
}

func runKiller(timeout int, chPid <-chan int, chQuit <-chan bool) {
	pid := <-chPid
	chRingClock := make(chan bool)
	defer close(chRingClock)

	go func(timeout int, ch chan<- bool) {
		time.Sleep(time.Duration(timeout) * time.Second)
		ch <- true
	}(timeout, chRingClock)

	select {
	case <-chRingClock:
		err := syscall.Kill(-pid, syscall.SIGKILL)
		if err != nil {
			panic(err)
		}
		return
	case <-chQuit:
		return
	}
}

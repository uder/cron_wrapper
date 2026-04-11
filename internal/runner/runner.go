package runner

import (
	"os"
	"sync"
	"time"

	"golang.org/x/sys/unix"

	"cron_wrapper/internal/command"
)

type runChannels struct {
	chPid        chan int
	chKillerQuit chan bool
}

func (rc *runChannels) cleanup() {
	close(rc.chPid)
	close(rc.chKillerQuit)
}

func newRunChannels() *runChannels {
	return &runChannels{
		chPid:        make(chan int),
		chKillerQuit: make(chan bool),
	}
}

func Run(cmd *command.Command) {
	wgCommand := new(sync.WaitGroup)

	wgCommand.Add(1)
	// TODO: make command fields tread safe. add mutex?

	chs := newRunChannels()
	defer chs.cleanup()
	go runCommand(cmd, wgCommand, chs)

	// TODO: replace killer thread using context package?
	go runKiller(cmd.Timeout(), chs)
	wgCommand.Wait()
}

func runCommand(cmd *command.Command, waitGroup *sync.WaitGroup, chs *runChannels) {
	proc, err := os.StartProcess("/bin/bash", *cmd.CommandToExecute(), cmd.ProcAttrs())
	if err != nil {
		panic(err)
	}
	chs.chPid <- proc.Pid

	pState, err := proc.Wait()
	if err != nil {
		panic(err)
	}

	cmd.SetProcState(pState)
	waitGroup.Done()

	if pState.Exited() {
		chs.chKillerQuit <- true
	}
}

func runKiller(timeout int, chs *runChannels) {
	pid := <-chs.chPid
	chRingClock := make(chan bool)
	defer close(chRingClock)

	go func(timeout int, ch chan<- bool) {
		time.Sleep(time.Duration(timeout) * time.Second)
		ch <- true
	}(timeout, chRingClock)

	select {
	case <-chRingClock:
		err := unix.Kill(-pid, unix.SIGKILL)
		if err != nil {
			panic(err)
		}
		return
	case <-chs.chKillerQuit:
		return
	}
}

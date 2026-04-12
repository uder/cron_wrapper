package runner

import (
	"context"
	"os"
	"time"

	"cron_wrapper/internal/command"
)

type quit struct{}

func Run(cmd *command.Command) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cmd.Timeout())*time.Second)
	defer cancel()

	chQuit := make(chan quit)

	proc, err := os.StartProcess("/bin/bash", *cmd.CommandToExecute(), cmd.ProcAttrs())

	//TODO: add debug logs
	if err != nil {
		panic(err)
	}
	cmd.SetPid(proc.Pid)

	go func() {
		//TODO: add debug logs
		pState, err := proc.Wait()
		if err != nil {
			panic(err)
		}
		cmd.SetProcState(pState)
		chQuit <- quit{}
	}()

	select {
	case <-chQuit:
		return
	case <-ctx.Done():
		cmd.SetIsTimedOut(true)
		proc.Kill()
		<-chQuit
	}
}

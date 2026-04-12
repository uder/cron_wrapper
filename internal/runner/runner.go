package runner

import (
	"context"
	"os"
	"time"
	"crypto/md5"
	"fmt"


	"github.com/gofrs/flock"

	"cron_wrapper/internal/command"
)

type quit struct{}

func getLockFilePath(cmd *command.Command) string {
	hash := md5.Sum([]byte(cmd.CommandLine()))
	return cmd.TmpDir() + "/" + "_cron_wrapper_" + fmt.Sprintf("%x", hash) + ".lock"
}

func Run(cmd *command.Command) {
	if !cmd.IsParallel() {
		lockFilePath := getLockFilePath(cmd)
		fileLock := flock.New(lockFilePath)
		locked, err := fileLock.TryLock()

		//TODO: add debug logs
		if err != nil {
			panic(err)
		}

		//TODO: add debug logs
		if !locked {
			panic ("Another instance of the command is running")
		}
		defer fileLock.Close()
	}

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

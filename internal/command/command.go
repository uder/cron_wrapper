package command

import (
	"math/rand/v2"
	"os"
	"syscall"
	"time"

	"cron_wrapper/internal/args"
)

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknownHost"
	}
	return hostname
}

func generateId(length int) string {
	var hex = []rune("0123456789abcdef")
	b := make([]rune, length)
	for i := range b {
		b[i] = hex[rand.IntN(len(hex))]
	}
	return string(b)
}

type Command struct {
	runId     string
	startTs   time.Time
	endTs     time.Time
	hostname  string
	cliArgs   *args.Args
	procFiles *ProcFiles
	procState *os.ProcessState
}

func (c *Command) RunId() string {
	if c.runId == "" {
		c.runId = generateId(8)
	}
	return c.runId
}

func (c *Command) SetStartTs() {
	c.startTs = time.Now()
}

func (c *Command) StartTs() time.Time {
	return c.startTs
}

func (c *Command) SetEndTs() {
	c.endTs = time.Now()
}

func (c *Command) EndTs() time.Time {
	return c.endTs
}

func (c *Command) SetProcState(pState *os.ProcessState) {
	c.procState = pState
}

func (c *Command) ProcState() *os.ProcessState {
	return c.procState
}

func (c *Command) GetDuration() float64 {
	return float64((c.endTs.UnixMilli() - c.startTs.UnixMilli()) / 1000)
}

func (c *Command) Hostname() string {
	return c.hostname
}

func (c *Command) Timeout() int {
	return c.cliArgs.Timeout
}

func (c *Command) CommandLine() string {
	return c.cliArgs.Command
}

func (c *Command) CommandToExecute() *[]string {
	return &[]string{
		"/usr/bin/env bash",
		"-c",
		c.CommandLine(),
	}
}

func (c *Command) ProcAttrs() *os.ProcAttr {
	return &os.ProcAttr{
		Files: c.procFiles.getProcAttr(),
		Sys:   &syscall.SysProcAttr{Setpgid: true},
	}
}

func (c *Command) ReportsDir() string {
	return c.cliArgs.ReportsDir
}

func (c *Command) EnableBegin() bool {
	return c.cliArgs.EnableBegin
}

func (c *Command) EnableStdoutOnSuccess() bool {
	return c.cliArgs.EnableStdoutOnSuccess
}

func (c *Command) StdoutPath() string {
	return c.procFiles.stdout.Name()
}

func (c *Command) StderrPath() string {
	return c.procFiles.stderr.Name()
}

func (c *Command) Cleanup() {
	c.procFiles.cleanup()
}

func NewCommand(a *args.Args) *Command {
	return &Command{
		runId:     generateId(8),
		hostname:  getHostname(),
		cliArgs:   a,
		procFiles: NewProcFiles(a.TmpDir),
	}
}

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

package report

import (
	"time"

	"cron_wrapper/internal/command"
)

type EndingReport struct {
	startTs              time.Time
	exitCode             int
	duration             float64
	isTimedOut           bool
	commandLine          string
	runId                string
	hostname             string
	pid                  int
	stdout               string
	stderr               string
	reportsDir           string
	enableStoutOnSuccess bool
	isDisableChat        bool
}

func (r *EndingReport) Type() string {
	if r.isTimedOut {
		return "TIMEOUT"
	}

	if r.exitCode == 0 {
		return "INFO"
	} else {
		return "ERROR"
	}
}

func NewEndingReport(cmd *command.Command) *EndingReport {
	return &EndingReport{
		startTs:              cmd.StartTs(),
		exitCode:             cmd.ProcState().ExitCode(),
		duration:             cmd.GetDuration(),
		isTimedOut:           cmd.IsTimedOut(),
		commandLine:          cmd.CommandLine(),
		runId:                cmd.RunId(),
		hostname:             cmd.Hostname(),
		pid:                  cmd.Pid(),
		stdout:               readFile(cmd.StdoutPath()),
		stderr:               readFile(cmd.StderrPath()),
		reportsDir:           cmd.ReportsDir(),
		enableStoutOnSuccess: cmd.EnableStdoutOnSuccess(),
		isDisableChat:           cmd.IsDisableChat(),
	}
}

package report

import (
	"time"

	"cron_wrapper/internal/command"
)

type EndReport struct {
	startTs              time.Time
	exitCode             int
	duration             float64
	commandLine          string
	runId                string
	hostname             string
	pid                  int
	stdout               string
	stderr               string
	reportsDir           string
	enableStoutOnSuccess bool
}

func (r *EndReport) Type() string {
	if r.exitCode == 0 {
		return "INFO"
	} else {
		return "ERROR"
	}
}

func NewEndReport(cmd *command.Command) *EndReport {
	return &EndReport{
		startTs:              cmd.StartTs(),
		exitCode:             cmd.ProcState().ExitCode(),
		duration:             cmd.GetDuration(),
		commandLine:          cmd.CommandLine(),
		runId:                cmd.RunId(),
		hostname:             cmd.Hostname(),
		pid:                  cmd.ProcState().Pid(),
		stdout:               readFile(cmd.StdoutPath()),
		stderr:               readFile(cmd.StderrPath()),
		reportsDir:           cmd.ReportsDir(),
		enableStoutOnSuccess: cmd.EnableStdoutOnSuccess(),
	}
}

package report

import (
	"time"

	"cron_wrapper/internal/command"
)

type BeginReport struct {
	startTs     time.Time
	commandLine string
	runId       string
	hostname    string
	reportsDir  string
	enableBegin bool
}

func (r *BeginReport) Type() string {
	return "BEGIN"
}

func NewBeginReport(cmd *command.Command) *BeginReport {
	return &BeginReport{
		startTs:     cmd.StartTs(),
		commandLine: cmd.CommandLine(),
		runId:       cmd.RunId(),
		hostname:    cmd.Hostname(),
		reportsDir:  cmd.ReportsDir(),
		enableBegin: cmd.EnableBegin(),
	}
}

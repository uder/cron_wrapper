package report

import (
	"time"

	"cron_wrapper/internal/command"
)

type BeginningReport struct {
	startTs     time.Time
	commandLine string
	runId       string
	hostname    string
	reportsDir  string
	enableBegin bool
}

func (r *BeginningReport) Type() string {
	return "BEGIN"
}

func NewBeginningReport(cmd *command.Command) *BeginningReport {
	return &BeginningReport{
		startTs:     cmd.StartTs(),
		commandLine: cmd.CommandLine(),
		runId:       cmd.RunId(),
		hostname:    cmd.Hostname(),
		reportsDir:  cmd.ReportsDir(),
		enableBegin: cmd.EnableBegin(),
	}
}

package report

import (
	"fmt"
	"strings"
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

func (r *BeginReport) String() string {
	header := strings.Join([]string{
		time.Now().Format(time.DateTime),
		r.Type(),
		"\"" + r.commandLine + "\"",
		r.runId},
		" ")

	body := r.hostname

	return strings.Join([]string{header, body, delimiter()}, "\n")
}

func (r *BeginReport) print() {
	if r.enableBegin {
		fmt.Print(r.String())
	}
}

func (r *BeginReport) Write() {
	r.print()
	writeToFile(getReportFileName(r.reportsDir), r.String())
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

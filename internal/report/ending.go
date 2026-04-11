package report

import (
	"fmt"
	"strconv"
	"strings"
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

func (r *EndReport) String() string {
	header := strings.Join([]string{
		time.Now().Format(time.DateTime),
		r.Type(),
		strconv.Itoa(r.exitCode),
		strconv.FormatFloat(r.duration, 'f', 3, 32),
		"\"" + r.commandLine + "\"",
		r.runId},
		" ")

	meta := strings.Join([]string{
		r.hostname,
		"begin_at=" + r.startTs.Format(time.DateTime),
		"pid=" + strconv.Itoa(r.pid),
	},
		" ")
	stdout := strings.TrimSuffix(strings.Join([]string{"STDOUT:", r.stdout}, "\n"), "\n")
	stderr := strings.TrimSuffix(strings.Join([]string{"STDERR:", r.stderr}, "\n"), "\n")

	return strings.Join([]string{header, meta, stdout, stderr, delimiter()}, "\n")
}

func (r *EndReport) StringReport(withStdOut bool, withStdErr bool) string {
	header := strings.Join([]string{
		time.Now().Format(time.DateTime),
		r.Type(),
		strconv.Itoa(r.exitCode),
		strconv.FormatFloat(r.duration, 'f', 3, 32),
		"\"" + r.commandLine + "\"",
		r.runId},
		" ")

	meta := strings.Join([]string{
		r.hostname,
		"begin_at=" + r.startTs.Format(time.DateTime),
		"pid=" + strconv.Itoa(r.pid),
	},
		" ")

	var stdout, stderr string
	if withStdOut {
		stdout = strings.TrimSuffix(strings.Join([]string{"STDOUT:", r.stdout}, "\n"), "\n")
	}
	if withStdErr {
		stderr = strings.TrimSuffix(strings.Join([]string{"STDERR:", r.stderr}, "\n"), "\n")
	}

	rep := strings.Join([]string{header, meta}, "\n")
	for _, item := range []string{stdout, stderr, delimiter()} {
		if item != "" {
			rep = strings.Join([]string{rep, item}, "\n")
		}
	}
	return rep
}

func (r *EndReport) print() {
	var rep string
	if r.exitCode == 0 && r.enableStoutOnSuccess {
		rep = r.StringReport(true, true)
	}
	if r.exitCode == 0 && !r.enableStoutOnSuccess {
		rep = r.StringReport(false, false)
	}
	if r.exitCode != 0 {
		rep = r.StringReport(true, true)
	}
	fmt.Print(rep)
}

func (r *EndReport) Write() {
	r.print()
	WriteToFile(GetReportFileName(r.reportsDir), r.String())
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

package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func prepareReportDir(reportDir string) {
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		err = os.Mkdir(reportDir, 0755)
	}
}

func readFile(path string) string {
	body, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(body)
}
func writeToFile(filename string, record string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	_, err = f.WriteString(record)
	if err != nil {
		panic(err)
	}
}

func getReportFileName(dir string) string {
	filename := time.Now().Format(time.DateOnly) + ".log"
	return path.Join(dir, filename)
}

func delimiter() string {
	return "---\n"
}

type BeginReport struct {
	startTs     time.Time
	commandLine string
	runId       string
	hostname    string
	// db          database.DB
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
	// r.insertDb()
	writeToFile(getReportFileName(r.reportsDir), r.String())
}

func NewBeginReport(command *Command) *BeginReport {
	return &BeginReport{
		startTs:     command.StartTs(),
		commandLine: command.CommandLine(),
		runId:       command.RunId(),
		hostname:    command.Hostname(),
		reportsDir:  command.ReportsDir(),
		enableBegin: command.enableBegin(),
	}
}

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

	report := strings.Join([]string{header, meta}, "\n")
	for _, item := range []string{stdout, stderr, delimiter()} {
		if item != "" {
			report = strings.Join([]string{report, item}, "\n")
		}
	}
	return report
}

func (r *EndReport) print() {
	var report string
	if r.exitCode == 0 && r.enableStoutOnSuccess {
		report = r.StringReport(true, true)
	}
	if r.exitCode == 0 && !r.enableStoutOnSuccess {
		report = r.StringReport(false, false)
	}
	if r.exitCode != 0 {
		report = r.StringReport(true, true)
	}
	fmt.Print(report)
}

func (r *EndReport) Write() {
	r.print()
	writeToFile(getReportFileName(r.reportsDir), r.String())
}

func NewEndReport(command *Command) *EndReport {
	return &EndReport{
		startTs:              command.StartTs(),
		exitCode:             command.procState.ExitCode(),
		duration:             command.GetDuration(),
		commandLine:          command.CommandLine(),
		runId:                command.RunId(),
		hostname:             command.Hostname(),
		pid:                  command.procState.Pid(),
		stdout:               readFile(command.procFiles.stdout.Name()),
		stderr:               readFile(command.procFiles.stderr.Name()),
		reportsDir:           command.ReportsDir(),
		enableStoutOnSuccess: command.enableStdoutOnSuccess(),
	}
}

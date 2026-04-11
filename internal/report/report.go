package report

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"cron_wrapper/internal/command"
)

func PrepareReportDir(reportDir string) {
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		err = os.Mkdir(reportDir, 0755)
	}
}

func readFile(filePath string) string {
	body, err := os.ReadFile(filePath)
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
	writeToFile(getReportFileName(r.reportsDir), r.String())
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

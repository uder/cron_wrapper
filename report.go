package main

import (
	"fmt"
	"main/database"
	"main/database/model"
	"os"
	"strconv"
	"strings"
	"time"
)

func readFile(path string) string {
	body, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func delimiter() string {
	return "---\n"
}

type BeginReport struct {
	startTs     time.Time
	commandLine string
	runId       string
	hostname    string
	db          database.DB
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

func (r *BeginReport) insertDb() {
	record := model.NewCommandBeginRecord(r.Type(), r.runId, r.commandLine)
	r.db.Conn().Create(record)
}

func (r *BeginReport) print() {
	fmt.Print(r.String())
}

func (r *BeginReport) Write() {
	r.print()
	r.insertDb()
}

func NewBeginReport(command *Command, db database.DB) *BeginReport {
	return &BeginReport{
		startTs:     command.StartTs(),
		commandLine: command.CommandLine(),
		runId:       command.RunId(),
		hostname:    command.Hostname(),
		db:          db,
	}
}

type EndReport struct {
	startTs     time.Time
	exitCode    int
	duration    float64
	commandLine string
	runId       string
	hostname    string
	pid         int
	stdout      string
	stderr      string
	db          database.DB
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

func (r *EndReport) print() {
	fmt.Print(r.String())
}

func (r *EndReport) insertDb() {
	record := model.NewCommandEndRecord(r.Type(), r.runId, r.commandLine, r.exitCode, r.duration)
	r.db.Conn().Create(record)
}

func (r *EndReport) Write() {
	r.print()
	r.insertDb()
}

func NewEndReport(command *Command, db database.DB) *EndReport {
	return &EndReport{
		startTs:     command.StartTs(),
		exitCode:    command.procState.ExitCode(),
		duration:    command.GetDuration(),
		commandLine: command.CommandLine(),
		runId:       command.RunId(),
		hostname:    command.Hostname(),
		pid:         command.procState.Pid(),
		stdout:      readFile(command.procFiles.stdout.Name()),
		stderr:      readFile(command.procFiles.stderr.Name()),
		db:          db,
	}
}

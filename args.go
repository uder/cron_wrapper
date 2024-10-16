package main

import (
	"github.com/alexflint/go-arg"
	"strconv"
	"strings"
)

type Args struct {
	Command  string `arg:"positional, required" help:"command to execute"`
	Timeout  int    `arg:"-t,--timeout" default:"7200" help:"timeout in seconds"`
	Parallel int    `arg:"-p,--parallel" default:"1" help:"max number of parallel executions"`
	Migrate  bool   `arg:"--migrate" default:"false" help:"perform migration. Exit after migration complete"`

	// TODO: change default to /tmp after development. Revise the flag name
	TmpDir string `arg:"--tmpdir" default:"./tmp" help:"directory to store temporary files"`

	// TODO: change default to /var/logs/? after development. Revise the flag name
	ReportsDir  string `arg:"--reportsdir" default:"./log" help:"directory to store reports"`
	EnableBegin bool   `arg:"-b" default:"false" help:"whether to send begin reports"`

	// TODO: revise the flag name and defaults for this flag
	EnableStdoutOnSuccess bool `arg:"-s" default:"false" help:"whether to send stdout on success"`
	EnableDebug           bool `arg:"-d" default:"false" help:"enable debug logging"`

	SqliteDatabase string `arg:"--sqlite-db" default:"db.sqlite" help:"sqlite database file"`
	dbType         string
}

func (a *Args) String() string {
	return strings.Join([]string{
		"Command: " + a.Command,
		"Timeout: " + strconv.Itoa(a.Timeout),
		"Parallel: " + strconv.Itoa(a.Parallel),
		"Migrate: " + strconv.FormatBool(a.Migrate),
		"Tmp Dir: " + a.TmpDir,
		"Reports Dir: " + a.ReportsDir,
		"Enable Begin: " + strconv.FormatBool(a.EnableBegin),
		"Enable StdoutOnSuccess: " + strconv.FormatBool(a.EnableStdoutOnSuccess),
		"Enable Debug: " + strconv.FormatBool(a.EnableDebug),
	},
		"; ")
}

func (a *Args) setDbType(dbType string) {
	a.dbType = dbType
}

func (a *Args) DbType() string {
	return a.dbType
}

func parseArgs() *Args {
	var args Args
	arg.MustParse(&args)
	if args.SqliteDatabase != "" {
		args.setDbType("sqlite")
	}
	return &args
}

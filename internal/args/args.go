package args

import (
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"
)

type Args struct {
	Command string `arg:"positional, required" help:"command to execute"`
	Timeout int    `arg:"-t,--timeout" default:"43200" help:"timeout in seconds"`

	EnableParallel bool `arg:"-p,--parallel" default:"false" help:"whether to alllow running in parallel mode"`

	// TODO: change default to /tmp after development. Revise the flag name
	TmpDir string `arg:"--tmpdir" default:"/tmp" help:"directory to store temporary files"`

	// TODO: change default to /var/logs/? after development. Revise the flag name
	ReportsDir  string `arg:"--reportsdir" default:"/var/log/100sp/cron/" help:"directory to store reports"`
	EnableBegin bool   `arg:"-b" default:"false" help:"whether to send begin reports"`

	// TODO: revise the flag name and defaults for this flag
	EnableStdoutOnSuccess bool `arg:"-s" default:"false" help:"whether to send stdout on success"`
	EnableDebug           bool `arg:"-d" default:"false" help:"enable debug logging"`

	//TODO: revise flag name and default value. Maybe the flag should be replaced with something else
	DisableChat           bool `arg:"--disable-chat" default:"false" help:"whether to disable chat notifications"`

	MeaninglessO           bool `arg:"-o" default:"false" help:"Do nothing. It was added for compability reasons"`
}

func (a *Args) String() string {
	return strings.Join([]string{
		"Command: " + a.Command,
		"Timeout: " + strconv.Itoa(a.Timeout),
		"Enable Parallel: " + strconv.FormatBool(a.EnableParallel),
		"Tmp Dir: " + a.TmpDir,
		"Reports Dir: " + a.ReportsDir,
		"Enable Begin: " + strconv.FormatBool(a.EnableBegin),
		"Enable StdoutOnSuccess: " + strconv.FormatBool(a.EnableStdoutOnSuccess),
		"Enable Debug: " + strconv.FormatBool(a.EnableDebug),
		"Disable Chat: " + strconv.FormatBool(a.DisableChat),
	},
		"; ")
}

func ParseArgs() *Args {
	var a Args
	arg.MustParse(&a)
	return &a
}

package main

import (
	"flag"
)

type Args struct {
	command string
	timeout int
	tmpdir  string
}

func parseArgs() *Args {
	args := new(Args)
	flag.IntVar(&args.timeout, "timeout", 7200, "timeout")
	flag.IntVar(&args.timeout, "t", 7200, "timeout")
	// TODO: change default to /tmp after development
	flag.StringVar(&args.tmpdir, "tmpdir", "./tmp", "path to tmp dir")

	flag.Parse()
	args.command = flag.Arg(0)
	return args
}

func main() {
	args := parseArgs()
	files := NewProcFiles("./tmp")
	defer files.cleanup()

	command := NewCommand(args)
	defer command.cleanup()

	command.SetStartTs()
	beginReport := NewBeginReport(command)
	beginReport.Write()

	run(command)
	command.SetEndTs()
	endReport := NewEndReport(command)
	endReport.Write()
}

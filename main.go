package main

import (
	"flag"
)

type Args struct {
	command string
	timeout int
}

func parseArgs() *Args {
	args := new(Args)
	flag.IntVar(&args.timeout, "timeout", 7200, "timeout")
	flag.IntVar(&args.timeout, "t", 7200, "timeout")

	flag.Parse()
	args.command = flag.Arg(0)
	return args
}

func main() {
	args := parseArgs()
	files := NewProcFiles("./tmp")
	defer files.cleanup()

	procState, duration := run(args, files)
	report := NewReport(procState, files, duration)
	report.Print()
}

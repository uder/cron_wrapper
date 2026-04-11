package main

func main() {
	args := parseArgs()
	logger := NewLogger(args.EnableDebug)
	logger.Info(args.String())

	// TODO: check duration precision. Why duration is always an int number of Seconds
	// TODO: Wrap up actions with debug logs

	command := NewCommand(args)
	defer command.cleanup()

	command.SetStartTs()
	prepareReportDir(args.ReportsDir)
	beginReport := NewBeginReport(command)
	beginReport.Write()

	run(command)

	command.SetEndTs()
	endReport := NewEndReport(command)
	endReport.Write()
}

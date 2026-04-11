package main

import (
	"cron_wrapper/internal/args"
	"cron_wrapper/internal/command"
	"cron_wrapper/internal/logging"
	"cron_wrapper/internal/report"
	"cron_wrapper/internal/runner"
)

func main() {
	cliArgs := args.ParseArgs()
	logger := logging.NewLogger(cliArgs.EnableDebug)
	logger.Info(cliArgs.String())

	// TODO: check duration precision. Why duration is always an int number of Seconds
	// TODO: Wrap up actions with debug logs

	cmd := command.NewCommand(cliArgs)
	defer cmd.Cleanup()

	cmd.SetStartTs()
	report.PrepareReportDir(cliArgs.ReportsDir)
	beginReport := report.NewBeginReport(cmd)
	beginReport.Write()

	runner.Run(cmd)

	cmd.SetEndTs()
	endReport := report.NewEndReport(cmd)
	endReport.Write()
}

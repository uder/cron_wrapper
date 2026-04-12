package main

import (
	"fmt"
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

	// TODO: Wrap up actions with debug logs

	cmd := command.NewCommand(cliArgs)
	defer cmd.Cleanup()

	cmd.SetStartTs()
	report.PrepareReportDir(cliArgs.ReportsDir)

	beginningReport := report.NewBeginningReport(cmd)
	fmt.Print(report.GetBegginingForSending(beginningReport))
	report.WriteToFile(report.GetReportFileName(cliArgs.ReportsDir), report.GetBegginingForLogging(beginningReport))

	runner.Run(cmd)

	cmd.SetEndTs()
	endReport := report.NewEndingReport(cmd)
	fmt.Print(report.GetEndingForSending(endReport))
	report.WriteToFile(report.GetReportFileName(cliArgs.ReportsDir), report.GetEndingForLogging(endReport))
}

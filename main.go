package main

import (
	"main/database"
	"main/database/model"
	"main/database/sqlite"
	"os"
	"strconv"
)

func getDataBase() database.DB {
	// TODO: cover out SqliteOpts?
	opts := &sqlite.Opts{FileName: "db.sqlite"}
	db, err := database.Factory("sqlite", opts)
	if err != nil {
		panic(err)
	}
	return db
}

func migrate(db database.DB, logger Logger) {
	for _, table := range model.Models() {
		err := database.Migrate(db, table)
		if err != nil {
			panic(err)
		}
		logger.Warn("Migration complete: " + table.TableName())
	}
}

func getNumberUnfinishedTasks(db database.DB, args *Args) int {
	unfinishedTasks := database.GetUnfinishedRecords(db, args.Command, args.Timeout)
	return len(unfinishedTasks)
}

func main() {
	args := parseArgs()
	logger := NewLogger(args.EnableDebug)
	logger.Info(args.String())
	db := getDataBase()
	if args.Migrate {
		migrate(db, logger)
		os.Exit(0)
	}

	if getNumberUnfinishedTasks(db, args) >= args.Parallel {
		logger.Error("Too many tasks are running already: " + strconv.Itoa(args.Parallel))
		os.Exit(1)
	} else {
		logger.Info("Running: " + strconv.Itoa(getNumberUnfinishedTasks(db, args)))
	}

	// TODO: check duration precision. Why duration is always an int number of Seconds
	// TODO: Wrap up actions with debug logs
	// TODO: Make database detachable. Disable parallelism check if db disabled
	// TODO: implement database specific cli params

	command := NewCommand(args)
	defer command.cleanup()

	command.SetStartTs()
	prepareReportDir(args.ReportsDir)
	beginReport := NewBeginReport(command, db)
	beginReport.Write()

	run(command)

	command.SetEndTs()
	endReport := NewEndReport(command, db)
	endReport.Write()
}

package main

import (
	"main/database"
	"main/database/model"
	"main/database/sqlite"
	"os"
	"strconv"
)

func getDataBase(args *Args) database.DB {
	// TODO: cover out SqliteOpts?
	// TODO: Check if refactor is needed. It looks like this func can be partially moved to the Factory or vice versa
	switch args.DbType() {
	case "sqlite":
		opts := &sqlite.Opts{FileName: args.SqliteDatabase}
		db, err := database.Factory(args.DbType(), opts)
		if err != nil {
			panic(err)
		}
		return db
	default:
		panic("Unsupported database type: " + args.DbType())
	}
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
	db := getDataBase(args)
	if args.Migrate {
		migrate(db, logger)
		os.Exit(0)
	}

	// TODO: Review two-level conditional.
	if args.DbType() != "" {
		if getNumberUnfinishedTasks(db, args) >= args.Parallel {
			logger.Error("Too many tasks are running already: " + strconv.Itoa(args.Parallel))
			os.Exit(1)
		} else {
			logger.Info("Running: " + strconv.Itoa(getNumberUnfinishedTasks(db, args)))
		}
	}

	// TODO: check duration precision. Why duration is always an int number of Seconds
	// TODO: Wrap up actions with debug logs

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

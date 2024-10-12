package main

import (
	"flag"
	"log/slog"
	"main/database"
	"main/database/model"
	"main/database/sqlite"
	"os"
	"strconv"
)

type Args struct {
	command  string
	timeout  int
	parallel int
	migrate  bool
	tmpdir   string
}

func parseArgs() *Args {
	args := new(Args)
	flag.IntVar(&args.timeout, "timeout", 7200, "timeout")
	flag.IntVar(&args.timeout, "t", 7200, "timeout")
	flag.IntVar(&args.parallel, "p", 1, "max allowed parallel execution. Key is a command line. Default: 1")
	// TODO: change default to /tmp after development
	flag.StringVar(&args.tmpdir, "tmpdir", "./tmp", "path to tmp dir")

	flag.BoolVar(&args.migrate, "migrate", false, "whether to migrate. All other args will be ignored")

	flag.Parse()
	args.command = flag.Arg(0)
	return args
}

func getDataBase() database.DB {
	// TODO: cover out SqliteOpts?
	opts := &sqlite.Opts{FileName: "db.sqlite"}
	db, err := database.Factory("sqlite", opts)
	if err != nil {
		panic(err)
	}
	return db
}

func migrate(db database.DB, logger *slog.Logger) {
	for _, table := range model.Models() {
		err := database.Migrate(db, table)
		if err != nil {
			panic(err)
		}
		logger.Warn("Migration complete: " + table.TableName())
	}
}
func getLogger() *slog.Logger {
	logger := slog.Default()
	return logger
}

func getNumberUnfinishedTasks(db database.DB, args *Args) int {
	unfinishedTasks := database.GetUnfinishedRecords(db, args.command, args.timeout)
	return len(unfinishedTasks)
}

func main() {
	logger := getLogger()
	args := parseArgs()
	db := getDataBase()
	if args.migrate {
		migrate(db, logger)
		os.Exit(0)
	}

	if getNumberUnfinishedTasks(db, args) >= args.parallel {
		logger.Error("Too many tasks are running already: " + strconv.Itoa(args.parallel))
		os.Exit(1)
	} else {
		logger.Info("Running: " + strconv.Itoa(getNumberUnfinishedTasks(db, args)))
	}

	// TODO: check duration precision. Why duration is always an int number of Seconds
	// TODO: Implement logging to files.
	// TODO: Implement configuration of logging. Enable/Disable blocks and messages. eg Disable Begin messages or
	// TODO: disable STDOUT on exitCode == 0

	command := NewCommand(args)
	defer command.cleanup()

	command.SetStartTs()
	beginReport := NewBeginReport(command, db)
	beginReport.Write()

	run(command)

	command.SetEndTs()
	endReport := NewEndReport(command, db)
	endReport.Write()
}

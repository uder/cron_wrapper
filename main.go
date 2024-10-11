package main

import (
	"flag"
	"log/slog"
	"main/database"
	"main/database/model"
	"main/database/sqlite"
)

type Args struct {
	command string
	timeout int
	migrate bool
	tmpdir  string
}

func parseArgs() *Args {
	args := new(Args)
	flag.IntVar(&args.timeout, "timeout", 7200, "timeout")
	flag.IntVar(&args.timeout, "t", 7200, "timeout")
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
func main() {
	//logger := getLogger()
	args := parseArgs()
	//db := getDataBase()
	//if args.migrate {
	//	migrate(db, logger)
	//	os.Exit(0)
	//}

	// TODO: Make a command as an object and run it through the workflow
	// TODO: Parametrize tmp dir

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

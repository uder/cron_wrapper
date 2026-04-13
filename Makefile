.PHONY: run run-error run-timeout build-linux

run:
	go run ./cmd/cron_wrapper/ "./test_command/test_cron.sh 0 3" -b -s

run-error:
	go run ./cmd/cron_wrapper/ "./test_command/test_cron.sh 1 3" -b

run-timeout:
	go run ./cmd/cron_wrapper/ "./test_command/test_cron.sh 1 3" -t 3 -b

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/cron_wrapper ./cmd/cron_wrapper/

package report

import (
	"strconv"
	"strings"
	"time"
)

func GetEndingForSending(r *EndingReport) string {
	if r.isDisableChat {
		return ""
	}
	if r.exitCode != 0 {
		return getFullEndingReport(r)
	}
	if r.exitCode == 0 && r.enableStoutOnSuccess {
		return getFullEndingReport(r)
	}
	return getShortEndingReport(r)
}

func GetEndingForLogging(r *EndingReport) string {
	return getFullEndingReport(r)
}

func getFullEndingReport(r *EndingReport) string {
	return EndingStringReport(r, true, true)
}

func getShortEndingReport(r *EndingReport) string {
	return EndingStringReport(r, false, false)
}

func EndingStringReport(r *EndingReport, withStdOut bool, withStdErr bool) string {
	header := strings.Join([]string{
		time.Now().Format(time.DateTime),
		r.Type(),
		strconv.Itoa(r.exitCode),
		strconv.FormatFloat(r.duration, 'f', 3, 32),
		"\"" + r.commandLine + "\"",
		r.runId},
		" ")

	meta := strings.Join([]string{
		r.hostname,
		"begin_at=" + r.startTs.Format(time.DateTime),
		"pid=" + strconv.Itoa(r.pid),
	},
		" ")

	var stdout, stderr string
	if withStdOut {
		stdout = strings.TrimSuffix(strings.Join([]string{"STDOUT:", r.stdout}, "\n"), "\n")
	}
	if withStdErr {
		stderr = strings.TrimSuffix(strings.Join([]string{"STDERR:", r.stderr}, "\n"), "\n")
	}

	rep := strings.Join([]string{header, meta}, "\n")
	for _, item := range []string{stdout, stderr, delimiter()} {
		if item != "" {
			rep = strings.Join([]string{rep, item}, "\n")
		}
	}
	return rep
}

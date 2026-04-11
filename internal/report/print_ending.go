package report

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	)

func PrintEnding(r *EndReport) {
	var rep string
	if r.exitCode == 0 && r.enableStoutOnSuccess {
		rep = EndingStringReport(r, true, true)
	}
	if r.exitCode == 0 && !r.enableStoutOnSuccess {
		rep = EndingStringReport(r, false, false)
	}
	if r.exitCode != 0 {
		rep = EndingStringReport(r, true, true)
	}
	fmt.Print(rep)
}

func EndReportString(r *EndReport) string {
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
	stdout := strings.TrimSuffix(strings.Join([]string{"STDOUT:", r.stdout}, "\n"), "\n")
	stderr := strings.TrimSuffix(strings.Join([]string{"STDERR:", r.stderr}, "\n"), "\n")

	return strings.Join([]string{header, meta, stdout, stderr, delimiter()}, "\n")
}

func EndingStringReport(r *EndReport, withStdOut bool, withStdErr bool) string {
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

// func (r *EndReport) print() {
// 	var rep string
// 	if r.exitCode == 0 && r.enableStoutOnSuccess {
// 		rep = r.StringReport(true, true)
// 	}
// 	if r.exitCode == 0 && !r.enableStoutOnSuccess {
// 		rep = r.StringReport(false, false)
// 	}
// 	if r.exitCode != 0 {
// 		rep = r.StringReport(true, true)
// 	}
// 	fmt.Print(rep)
// }

// func (r *EndReport) Write() {
// 	r.print()
// 	WriteToFile(GetReportFileName(r.reportsDir), r.String())
// }

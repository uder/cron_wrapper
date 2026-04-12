package report

import (
	"strings"
	"time"
)

func GetBegginingForSending(r *BeginningReport) string {
	if r.enableBegin {
		return BegginingString(r)
	}
	return ""
}

func GetBegginingForLogging(r *BeginningReport) string {
	return BegginingString(r)
}

func BegginingString(r *BeginningReport) string {
	header := strings.Join([]string{
		time.Now().Format(time.DateTime),
		r.Type(),
		"\"" + r.commandLine + "\"",
		r.runId},
		" ")

	body := r.hostname

	return strings.Join([]string{header, body, delimiter()}, "\n")
}

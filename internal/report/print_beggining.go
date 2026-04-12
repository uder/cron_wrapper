package report

import (
	"strings"
	"time"
)

func BegginingString(r *BeginReport) string {
	if !r.enableBegin {
		return ""
	}

	header := strings.Join([]string{
		time.Now().Format(time.DateTime),
		r.Type(),
		"\"" + r.commandLine + "\"",
		r.runId},
		" ")

	body := r.hostname

	return strings.Join([]string{header, body, delimiter()}, "\n")
}

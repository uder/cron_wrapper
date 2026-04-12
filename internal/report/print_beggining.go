package report

import (
	"fmt"
	"strings"
	"time"
)

func PrintBeggining(r *BeginReport) {
	if r.enableBegin {
		fmt.Print(BegginingString(r))
	}
}

func BegginingString(r *BeginReport) string {
	header := strings.Join([]string{
		time.Now().Format(time.DateTime),
		r.Type(),
		"\"" + r.commandLine + "\"",
		r.runId},
		" ")

	body := r.hostname

	return strings.Join([]string{header, body, delimiter()}, "\n")
}

package report

import (
	"os"
	"path"
	"time"
)

func PrepareReportDir(reportDir string) {
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		err = os.Mkdir(reportDir, 0755)
	}
}

func readFile(filePath string) string {
	body, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func WriteToFile(filename string, record string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	_, err = f.WriteString(record)
	if err != nil {
		panic(err)
	}
}

func GetReportFileName(dir string) string {
	filename := time.Now().Format(time.DateOnly) + ".log"
	return path.Join(dir, filename)
}

func delimiter() string {
	return "---\n"
}

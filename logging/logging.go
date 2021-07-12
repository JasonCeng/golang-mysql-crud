package logging

import (
	"github.com/op/go-logging"
	"os"
)

var loggingDefaultLevel = logging.INFO

func Initialize() {
	consoleFormat := logging.MustStringFormatter(
		"%{color}%{time:2006-01-02T15:04:05.000}+0800 [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}",
	)

	consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	consoleBackendFormatter := logging.NewBackendFormatter(consoleBackend, consoleFormat)

	fileFormat := logging.MustStringFormatter(
		"%{color}%{time:2006-01-02T15:04:05.000}+0800 [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x} %{message}",
	)

	fileWriter := NewFileWriter()
	fileBackend := logging.NewLogBackend(fileWriter, "", 0)
	fileBackendFormatter := logging.NewBackendFormatter(fileBackend, fileFormat)

	logging.SetBackend(consoleBackendFormatter, fileBackendFormatter).SetLevel(loggingDefaultLevel, "")
}
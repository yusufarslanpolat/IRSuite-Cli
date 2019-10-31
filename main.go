package main

import (
	"os"
	"path/filepath"

)
import errors "github.com/pkg/errors"
import "gopkg.in/alecthomas/kingpin.v2"


func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(errors.Cause(err)) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "Create Directory Failed.")
		}
	}
	return nil
}
func createWorkingEnv(caseName string) {
	createDirIfNotExist(caseName)
	createDirIfNotExist(filepath.Join(caseName, "Windows Event Logs"))
}
func CollectEventLogs(path string) error{
	err := Dir(filepath.Join(os.Getenv("SystemRoot"), "System32", "Winevt", "Logs"),path)
	if err != nil {
		return errors.Wrap(err, "CollectEventLogs Failed.")
	}
	return nil
}

var (
	app      = kingpin.New("irsuite", "A incident response suite.")
	caseCommand     = app.Command("case", "Create a new case.")
	caseName = caseCommand.Arg("Case Name", "Case Name.").Required().String()
	debug   = kingpin.Flag("debug", "Enable debug mode.").Bool()
)
func main() {
	kingpin.Version("0.0.1")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case caseCommand.FullCommand():
		if _, err := os.Stat(*caseName); !os.IsNotExist(err) {
			panic("Directory exists.")
		}
		createWorkingEnv(*caseName)
		CollectEventLogs(filepath.Join(*caseName, "Windows Event Logs"))
	}
}
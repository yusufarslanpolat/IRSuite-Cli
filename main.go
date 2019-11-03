// +build windows
// +build 386

package main

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
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
func CollectEventLogs(path string) error {
	err := Dir(filepath.Join(os.Getenv("SystemRoot"), "System32", "Winevt", "Logs"), path)
	if err != nil {
		return errors.Wrap(err, "CollectEventLogs Failed.")
	}
	return nil
}

type Config struct {
	Targets []string `yaml:"Targets"`
}
type Artifact struct {
	Description string  `yaml:"Description"`
	Author      string  `yaml:"Author"`
	Version     float64 `yaml:"Version"`
	Targets     []struct {
		Name        string `yaml:"Name"`
		Category    string `yaml:"Category"`
		Path        string `yaml:"Path"`
		IsDirectory bool   `yaml:"IsDirectory"`
		Recursive   bool   `yaml:"Recursive"`
		Comment     string `yaml:"Comment"`
	} `yaml:"Targets"`
}

var (
	app         = kingpin.New("irsuite", "A incident response suite.")
	caseCommand = app.Command("case", "Create a new case.")
	caseName    = caseCommand.Arg("Case Name", "New case.").Required().String()
	config      = caseCommand.Arg("Config file", "YAML config file.(Optinal)").ExistingFile()
	debug       = kingpin.Flag("debug", "Enable debug mode.").Bool()
)

func main() {
	start := time.Now()
	// Set up the logging.
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
	log.SetOutput(colorable.NewColorableStdout())
	log.Info("Started acquisition ", *caseName)
	kingpin.Version("0.0.0.1")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case caseCommand.FullCommand():
		if _, err := os.Stat(*caseName); !os.IsNotExist(err) {
			panic("Directory exists. Please use empty folder for the case.")
		}
		createWorkingEnv(*caseName)
		//CollectEventLogs(filepath.Join(*caseName, "Windows Event Logs"))
		//case configCommand.FullCommand():
		//	if _, err := os.Stat(*caseName); !os.IsNotExist(err) {
		//		panic("Config file doesn't exists.")
		//	}
		//	println("Config ok.")
		//	//CollectEventLogs(filepath.Join(*caseName, "Windows Event Logs"))
	}
	configfile, _ := filepath.Abs(*config)
	yamlFile, err := ioutil.ReadFile(configfile)
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	for _, element := range config.Targets {
		//fmt.Println(index, "=>", element)
		artifactfile, _ := filepath.Abs(filepath.Join("Artifacts", element))
		yamlFile, err := ioutil.ReadFile(artifactfile)
		var artifact Artifact
		err = yaml.Unmarshal(yamlFile, &artifact)
		if err != nil {
			panic(err)
		}
		createDirIfNotExist(filepath.Join(*caseName, artifact.Description))
		for _, elementArtifact := range artifact.Targets {
			//fmt.Println(indexArtifact, "=>", elementArtifact.Path)
			matches, err := doublestar.Glob(elementArtifact.Path)
			//matches, err := filepath.Glob(elementArtifact.Path)

			if err != nil {
				errors.Wrap(err, "Can not find evidence for"+element)
			}

			for i := 0; i < len(matches); i++ {
				//fmt.Println(matches[i])
				filename := filepath.Base(matches[i])
				Info("\n" + filename)
				if elementArtifact.IsDirectory == false {
					err := File(matches[i], filepath.Join(*caseName, artifact.Description, filename))
					if err != nil {
						errors.Wrap(err, "Copy Evidence Failed.")
					}
				} else {
					err := Dir(matches[i], filepath.Join(*caseName, artifact.Description, filename))
					if err != nil {
						errors.Wrap(err, "Copy Evidence Failed.")
					}
				}
			}
		}
	}
	getClipboard(*caseName)
	getAutoruns(*caseName)
	getSystemInfo(*caseName)
	getArpTable(*caseName)
	getProxyConfig(*caseName)
	getVSSConfig(*caseName)
	getRecyclebinitems(*caseName)
	elapsed := time.Since(start)
	fmt.Printf("\nCollecting evidence done. %s", elapsed)
}

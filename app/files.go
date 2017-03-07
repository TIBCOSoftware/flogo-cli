package app

import (
	"os"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileDescriptor     string = "flogo.json"
	fileMainGo         string = "main.go"
	fileImportsGo      string = "imports.go"

	pathFlogoLib string = "github.com/TIBCOSoftware/flogo-lib"
)


func createMainGoFile(codeSourcePath string) {

	var data struct{}
	f, _ := os.Create(fgutil.Path(codeSourcePath, fileMainGo))
	fgutil.RenderTemplate(f, tplNewMainGoFile, &data)
	f.Close()
}

var tplNewMainGoFile = `package main

import (
	"fmt"
    "encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/types"
)

func main() {

	flogo, err := os.Open("flogo.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    jsonParser := json.NewDecoder(flogo)
    app := &types.AppConfig{}
    jsonParser.Decode(&app)
    e, err := engine.New(app)

	e.Start()

	exitChan := setupSignalHandling()

	code := <-exitChan

	e.Stop()

	os.Exit(code)
}

func setupSignalHandling() chan int {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGHUP
			case syscall.SIGHUP:
				exitChan <- 0
			// kill -SIGINT/Ctrl+c
			case syscall.SIGINT:
				exitChan <- 0
			// kill -SIGTERM
			case syscall.SIGTERM:
				exitChan <- 0
			// kill -SIGQUIT
			case syscall.SIGQUIT:
				exitChan <- 0
			default:
				logger.Debug("Unknown signal.")
				exitChan <- 1
			}
		}
	}()

	return exitChan
}
`


func createImportsGoFile(codeSourcePath string, refs []string) {
	f, _ := os.Create(fgutil.Path(codeSourcePath, fileImportsGo))
	fgutil.RenderTemplate(f, tplNewImportsGoFile, refs)
	f.Close()
}

var tplNewImportsGoFile = `package main

import (

{{range $i, $ref := .}}	_ "{{ $ref }}"
{{end}}
)
`

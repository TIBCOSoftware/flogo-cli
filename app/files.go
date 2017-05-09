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


func createMainGoFile(codeSourcePath string, flogoJSON string) {

	data := struct {
		FlogoJSON string
	}{
		flogoJSON,
	}

	f, _ := os.Create(fgutil.Path(codeSourcePath, fileMainGo))
	fgutil.RenderTemplate(f, tplNewMainGoFile, &data)
	f.Close()
}

var tplNewMainGoFile = `package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var (
	cp app.ConfigProvider
)

func main() {

	if cp == nil {
		// Use default config provider
		cp = app.DefaultConfigProvider()
	}

	app, err := cp.GetApp()
	if err != nil {
        	fmt.Println(err.Error())
        	os.Exit(1)
    	}

    	e, err := engine.New(app)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

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

func createImportsGoFile(codeSourcePath string, deps []*Dependency) error {
	f, err := os.Create(fgutil.Path(codeSourcePath, fileImportsGo))

	if err != nil {
		return err
	}

	fgutil.RenderTemplate(f, tplNewImportsGoFile, deps)
	f.Close()

	return nil
}

var tplNewImportsGoFile = `package main

import (

{{range $i, $dep := .}}	_ "{{ $dep.Ref }}"
{{end}}
)
`
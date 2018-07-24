package app

import (
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/flogo-cli/config"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileMainGo        string = "main.go"
	fileEmbeddedAppGo string = "embeddedapp.go"
	makeFile          string = "Makefile"
	gobuildFile       string = "build.go"
	fileShimGo        string = "shim.go"
	fileShimSupportGo string = "shim_support.go"

	dirShim      string = "shim"
	pathFlogoLib string = "github.com/TIBCOSoftware/flogo-lib"
)

func createMainGoFile(codeSourcePath string, flogoJSON string) {

	data := struct {
		FlogoJSON string
	}{
		flogoJSON,
	}

	f, _ := os.Create(filepath.Join(codeSourcePath, fileMainGo))
	fgutil.RenderTemplate(f, tplNewMainGoFile, &data)
	f.Close()
}

func removeMainGoFile(codeSourcePath string) {
	os.Remove(filepath.Join(codeSourcePath, fileMainGo))
}

var tplNewMainGoFile = `package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
    "runtime/pprof"
    "flag"
    "runtime"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("main-engine")
var cpuprofile = flag.String("cpuprofile", "", "Writes CPU profiling for the current process to the specified file")
var memprofile = flag.String("memprofile", "", "Writes memory profiling for the current process to the specified file")
var (
	cp app.ConfigProvider
)

func main() {

	var flogoApp *app.Config
	var err error

	if cp != nil {
		flogoApp, err = cp.GetApp()
	} else {
		flogoApp, err = app.LoadConfig("")
	}

	if err != nil {
        	fmt.Println(err.Error())
        	os.Exit(1)
    }

    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            fmt.Println(fmt.Sprintf("Failed to create CPU profiling file due to error - %s", err.Error()))
        	os.Exit(1)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    
    e, err := engine.New(app)
	if err != nil {
		log.Errorf("Failed to create engine instance due to error: %s", err.Error())
		os.Exit(1)
	}

	err = e.Start()
	if err != nil {
		log.Errorf("Failed to start engine due to error: %s", err.Error())
		os.Exit(1)
	}

	exitChan := setupSignalHandling()

	code := <-exitChan

	e.Stop()

    if *memprofile != "" {
        f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to create memory profiling file due to error - %s", err.Error()))
            os.Exit(1)
		}
		
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
          fmt.Println(fmt.Sprintf("Failed to write memory profiling data to file due to error - %s", err.Error()))
          os.Exit(1)
        }
        f.Close()
    }

	os.Exit(code)
}

func setupSignalHandling() chan int {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int, 1)
	select {
	case s := <-signalChan:
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			exitChan <- 0
		default:
			logger.Debug("Unknown signal.")
			exitChan <- 1
		}
	}
	return exitChan
}

`

func createImportsGoFile(codeSourcePath string, deps []*config.Dependency) error {
	f, err := os.Create(filepath.Join(codeSourcePath, config.FileImportsGo))

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

func createEmbeddedAppGoFile(codeSourcePath string, flogoJSON string) {

	data := struct {
		FlogoJSON string
	}{
		flogoJSON,
	}

	f, _ := os.Create(filepath.Join(codeSourcePath, fileEmbeddedAppGo))
	fgutil.RenderTemplate(f, tplEmbeddedAppGoFile, &data)
	f.Close()
}

func removeEmbeddedAppGoFile(codeSourcePath string) {
	os.Remove(filepath.Join(codeSourcePath, fileEmbeddedAppGo))
}

var tplEmbeddedAppGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/app"
)

// embedded flogo app descriptor file
const flogoJSON string = ` + "`{{.FlogoJSON}}`" + `

func init () {
	cp = EmbeddedProvider()
}

// embeddedConfigProvider implementation of ConfigProvider
type embeddedProvider struct {
}

//EmbeddedProvider returns an app config from a compiled json file
func EmbeddedProvider() (app.ConfigProvider){
	return &embeddedProvider{}
}

// GetApp returns the app configuration
func (d *embeddedProvider) GetApp() (*app.Config, error){
     return app.LoadConfig(flogoJSON)
}
`

func createShimSupportGoFile(codeSourcePath string, flogoJSON string, embeddedConfig bool) {

	configJson := ""

	if embeddedConfig {
		configJson = flogoJSON
	}

	data := struct {
		FlogoJSON string
	}{
		configJson,
	}

	f, _ := os.Create(filepath.Join(codeSourcePath, fileShimSupportGo))
	fgutil.RenderTemplate(f, tplShimSupportGoFile, &data)
	f.Close()
}

func removeShimGoFiles(codeSourcePath string) {
	os.Remove(filepath.Join(codeSourcePath, fileShimGo))
	os.Remove(filepath.Join(codeSourcePath, fileShimSupportGo))
}

var tplShimSupportGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/logger"

)

// embedded flogo app descriptor file
const flogoJSON string = ` + "`{{.FlogoJSON}}`" + `

func init() {
	config.SetDefaultLogLevel("ERROR")
	logger.SetLogLevel(logger.ErrorLevel)

	var cp app.ConfigProvider

	if flogoJSON != "" {
		cp = EmbeddedProvider()
	} else {
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

	e.Init(true)
}

// embeddedConfigProvider implementation of ConfigProvider
type embeddedProvider struct {
}

//EmbeddedProvider returns an app config from a compiled json file
func EmbeddedProvider() (app.ConfigProvider){
	return &embeddedProvider{}
}

// GetApp returns the app configuration
func (d *embeddedProvider) GetApp() (*app.Config, error){

	appCfg := &app.Config{}
	err := json.Unmarshal([]byte(flogoJSON), appCfg)
	if err != nil {
		return nil, err
	}
	return appCfg, nil
}
`

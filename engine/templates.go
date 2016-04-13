package engine

import (
	"os"

	"github.com/TIBCOSoftware/flogo/fgutil"
)

var tplMainGoFile = `package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

func main() {

	var format = logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}",
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	config := GetEngineConfig()

	logLevel, _ := logging.LogLevel(config.LogLevel)

	logging.SetLevel(logLevel, "")

	env := GetEngineEnvironment(config)

	engine := engine.NewEngine(env)
	engine.Start()

	exitChan := setupSignalHandling()

	code := <-exitChan

	engine.Stop()

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
				log.Debug("Unknown signal.")
				exitChan <- 1
			}
		}
	}()

	return exitChan
}
`
var tplEngineEnvGoFile = `package main

import (
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/service/processprovider/ppsremote"
	"github.com/TIBCOSoftware/flogo-lib/service/staterecorder/srsremote"
	"github.com/TIBCOSoftware/flogo-lib/service/tester"
)

// GetEngineEnvironment gets the engine environment
func GetEngineEnvironment(engineConfig *engine.Config) *engine.Environment {

	processProvider := ppsremote.NewRemoteProcessProvider()
	stateRecorder := srsremote.NewRemoteStateRecorder()
	engineTester := tester.NewRestEngineTester()

	env := engine.NewEnvironment(processProvider, stateRecorder, engineTester, engineConfig)
	env.SetEmbeddedJSONFlows(EmeddedFlowsAreCompressed(), EmeddedJSONFlows())

	return env
}
`

var tplEngineConfigGoFile = `package main

import (
	"github.com/TIBCOSoftware/flogo-lib/engine"
)

const configFileName string = "config.json"

// can be used to compile in config file
const configJSON string = ""

// GetEngineConfig gets the engine configuration
func GetEngineConfig() *engine.Config {

	config := engine.LoadConfigFromFile(configFileName)
	//config := engine.LoadConfigFromJSON(configJSON)

	if config == nil {
		config = engine.DefaultConfig()
		log.Warningf("Configuration file '%s' not found, using defaults", configFileName)
	}

	return config
}
`

var tplImportsGoFile = `package main

import (

	// activities
{{range .Activities}}{{if .Local}}	_ "activity/{{.Name}}/rt"{{end}}{{if not .Local}}	_ "{{.Path}}/rt"{{end}}
{{end}}
	// triggers
{{range .Triggers}}	_ "{{.Path}}/rt"
{{end}}
	// models
{{range .Models}}	_ "{{.Path}}"
{{end}}
)
`

func createFlowsGoFile(dir string, flows map[string]string) {
	// create flows Go file
	f, _ := os.Create(path(dir, fileFlowsGo))
	fgutil.RenderTemplate(f, tplFlowsGoFile, flows)
	f.Close()
}

var tplFlowsGoFile = `package main

var embeddedJSONFlows map[string]string

func init() {

	embeddedJSONFlows = make(map[string]string)

{{ range $key, $value := . }}	embeddedJSONFlows["{{ $key }}"] = "{{ $value }}"
{{ end }}
}

func EmeddedFlowsAreCompressed() bool {
	return true
}

func EmeddedJSONFlows() map[string]string {
	return embeddedJSONFlows
}
`

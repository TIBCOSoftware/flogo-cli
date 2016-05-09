package main

import (
	"os"

	"github.com/TIBCOSoftware/flogo/util"
)

const (
	fileDescriptor     string = "flogo.json"
	fileEngineConfig   string = "config.json"
	fileTriggersConfig string = "triggers.json"
	fileMainGo         string = "main.go"
	fileEnvGo          string = "env.go"
	fileConfigGo       string = "config.go"
	fileImportsGo      string = "imports.go"
	fileExprsGo        string = "exprs.go"

	dirFlows string = "flows"

	pathFlogoLib string = "github.com/TIBCOSoftware/flogo-lib"
)

func createMainGoFile(codeSourcePath string, projectDescriptor *FlogoProjectDescriptor) {
	f, _ := os.Create(path(codeSourcePath, fileMainGo))
	fgutil.RenderTemplate(f, tplMainGoFile, projectDescriptor)
	f.Close()
}

var tplMainGoFile = `package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/op/go-logging"
)

func init() {
	var format = logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{color:reset} %{message}",
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logging.SetLevel(logging.INFO, "")
}

var log = logging.MustGetLogger("main")

func main() {

	config := GetEngineConfig()
	triggersConfig := GetTriggersConfig()

	logLevel, _ := logging.LogLevel(config.LogLevel)

	logging.SetLevel(logLevel, "")

	env := GetEngineEnvironment(config, triggersConfig)

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

func createEngineEnvGoFile(codeSourcePath string, projectDescriptor *FlogoProjectDescriptor) {
	f, _ := os.Create(path(codeSourcePath, fileEnvGo))
	fgutil.RenderTemplate(f, tplEngineEnvGoFile, projectDescriptor)
	f.Close()
}

var tplEngineEnvGoFile = `package main

import (
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/service/flowprovider/ppsremote"
	"github.com/TIBCOSoftware/flogo-lib/service/staterecorder/srsremote"
	"github.com/TIBCOSoftware/flogo-lib/service/tester"
)

// GetEngineEnvironment gets the engine environment
func GetEngineEnvironment(engineConfig *engine.Config, triggersConfig *engine.TriggersConfig) *engine.Environment {

	flowProvider := ppsremote.NewRemoteFlowProvider()
	stateRecorder := srsremote.NewRemoteStateRecorder()
	engineTester := tester.NewRestEngineTester()

	env := engine.NewEnvironment(flowProvider, stateRecorder, engineTester, engineConfig, triggersConfig)
	env.SetEmbeddedJSONFlows(EmeddedFlowsAreCompressed(), EmeddedJSONFlows())

	return env
}
`
type ConfigInfo struct {
	Include     bool
	ConfigJSON  string
	TriggerJSON string
}

func createEngineConfigGoFile(codeSourcePath string, configInfo *ConfigInfo) {

	if configInfo == nil {
		configInfo = &ConfigInfo{Include:false, ConfigJSON:"", TriggerJSON:""}
	}

	f, _ := os.Create(path(codeSourcePath, fileConfigGo))
	fgutil.RenderTemplate(f, tplEngineConfigGoFile, configInfo)
	f.Close()
}

var tplEngineConfigGoFile = `package main

import (
	"github.com/TIBCOSoftware/flogo-lib/engine"
)

const configFileName string = "config.json"
const triggersConfigFileName string = "triggers.json"

// can be used to compile in config file
const configJSON string = ` + "`{{.ConfigJSON}}`" + `

// can be used to compile in triggers config file
const triggersConfigJSON string = ` + "`{{.TriggerJSON}}`" + `

// GetEngineConfig gets the engine configuration
func GetEngineConfig() *engine.Config {

	{{ if .Include }}//{{ end }}config := engine.LoadConfigFromFile(configFileName)
	{{ if not .Include }}//{{ end }}config := engine.LoadConfigFromJSON(configJSON)

	if config == nil {
		config = engine.DefaultConfig()
		log.Warningf("Configuration file '%s' not found, using defaults", configFileName)
	}

	return config
}

// GetTriggersConfig gets the triggers configuration
func GetTriggersConfig() *engine.TriggersConfig {

	{{ if .Include }}//{{ end }}config := engine.LoadTriggersConfigFromFile(triggersConfigFileName)
	{{ if not .Include }}//{{ end }}config := engine.LoadTriggersConfigFromJSON(triggersConfigJSON)

	if config == nil {
		config = engine.DefaultTriggersConfig()
		log.Warningf("Configuration file '%s' not found, using defaults", triggersConfigFileName)
	}

	return config
}
`

func createImportsGoFile(codeSourcePath string, projectDescriptor *FlogoProjectDescriptor) {
	f, _ := os.Create(path(codeSourcePath, fileImportsGo))
	fgutil.RenderTemplate(f, tplImportsGoFile, projectDescriptor)
	f.Close()
}

var tplImportsGoFile = `package main

import (

	// activities
{{range .Activities}}	_ "{{.Path}}/rt"
{{end}}
	// triggers
{{range .Triggers}}	_ "{{.Path}}/rt"
{{end}}
	// models
{{range .Models}}	_ "{{.Path}}"
{{end}}
)
`

func createFlowsGoFile(codeSourcePath string, flows map[string]string) {
	f, _ := os.Create(path(codeSourcePath, fileFlowsGo))
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

func createExprsGoFile(codeSourcePath string, flows map[string]map[int]string) {
	f, _ := os.Create(path(codeSourcePath, fileExprsGo))
	fgutil.RenderTemplate(f, tplExprsGoFile, flows)
	f.Close()
}


var tplExprsGoFile = `package main

import (
	"github.com/TIBCOSoftware/flogo-lib/script/fgn"
)


func init() {

	allFlowExprs := make(map[string]map[int]fgn.ExprFunc)
	var flowExprs map[int]fgn.ExprFunc

{{ range $key, $value := . }}	flowExprs = make(map[int]fgn.ExprFunc)
    {{ range $k, $v := $value }}flowExprs[{{ $k }}] = func(v map[string]interface{}) bool {
		return {{ $v }}
	}
	{{ end }}allFlowExprs["{{ $key }}"] = flowExprs

{{ end }}
}

`
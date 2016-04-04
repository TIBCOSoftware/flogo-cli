package engine

var tplMainGoFile = `package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/engine/ext/trigger"
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

	config := engine.LoadConfigurationFromFile("config.json")

	if config == nil {
		config = engine.NewConfiguration()
		log.Warning("Configuration file not found, using defaults")
	}

	logLevel, _ := logging.LogLevel(config.LogLevel)

	logging.SetLevel(logLevel, "")

	processRegistry := engine.NewProcessRegistry()
	stateService := engine.NewRestStateService(config.StateServiceURI)

	system := engine.NewSystem(processRegistry, stateService, config.EngineConfig)

	log.Info("Starting Engine...")

	engine := system.GetEngine()

	triggers := trigger.Triggers()

	// initialize triggers
	for _, trigger := range triggers {

		triggerConfig := config.Triggers[trigger.Metadata().ID]
		trigger.Init(engine, triggerConfig.Config)
	}

	engine.Start()

	// start triggers
	for _, trigger := range triggers {
		trigger.Start()
	}

	exitChan := setupSignalHandling()

	log.Info("Engine Running...")
	code := <-exitChan

	log.Infof("\nShutting Down Engine...\n")

	// stop triggers
	for _, trigger := range triggers {
		trigger.Stop()
	}
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
var tplImportsGoFile = `package main

import (

	// activities
{{range .Activities}}
	_ "{{.Path}}/rt"{{end}}

	// triggers
{{range .Triggers}}
	_ "{{.Path}}/rt"{{end}}

	// models
{{range .Models}}
	_ "{{.Path}}"{{end}}

)
`
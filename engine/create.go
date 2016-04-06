package engine

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TIBCOSoftware/flogo/fg"
	"github.com/TIBCOSoftware/flogo/fgutil"
)

const fileProjectConfig string = "engine.json"
const fileEngineConfig string = "config.json"
const fileMainGo string = "main.go"
const fileImportsGo string = "imports.go"

var optCreate = &flogo.OptionInfo{
	Name:      "create",
	UsageLine: "create [-with-ui] [-gb] engineName",
	Short:     "create an engine project",
	Long: `Creates a flogo engine project.

Options:
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option *flogo.OptionInfo
}

func (c *cmdCreate) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
}

func (c *cmdCreate) Exec(ctx *flogo.Context, args []string) error {

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Engine name not specified\n\n")
		Tool().CmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	engineName := args[0]

	if _, err := os.Stat(engineName); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot create engine project, directory '%s' already exists\n\n", engineName)
		os.Exit(2)
	}

	if !fgutil.ExecutableExists("gb") {
		fmt.Fprintf(os.Stderr, "Error: Cannot create engine project [%s], gb is not installed\n\n", engineName)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "Creating flogo engine '%s'...\n", engineName)

	basePath := engineName
	sourcePath := path(engineName, "src", engineName)
	vendorPath := path(engineName, "vendor", "src")

	os.MkdirAll(sourcePath, 0777)
	os.MkdirAll(vendorPath, 0777)

	fmt.Fprint(os.Stdout, "Installing flogo lib...\n")
	os.Chdir(engineName)
	cmd := exec.Command("gb", "vendor", "fetch", "github.com/TIBCOSoftware/flogo-lib")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	os.Chdir("..")

	// create engine.json file
	projectConfig := &EngineProjectConfig{
		Name:        engineName,
		Version:     "0.0.1",
		Description: "My engine description",
		Activities:  make([]*ItemConfig, 0),
		Triggers:    make([]*ItemConfig, 0),
		Models:      make([]*ItemConfig, 0),
	}

	// todo: add default model
	// todo: make a .flogo directory in user home, were people can put a default engine.json (use -default on create, or specify a json?)
	fgutil.WriteJSONtoFile(path(basePath, fileProjectConfig), projectConfig)

	// create engine.json file
	engineConfig := &EngineConfig{
		LogLevel:    "info",
		StateServiceURI: "",
		WorkerConfig: &WorkerConfig{NumWorkers:5, WorkQueueSize:50, MaxStepCount:32000},
		Triggers: make([]*TriggerConfig,0),
	}

	binPath := path(engineName, "bin")
	os.MkdirAll(binPath, 0777)

	fgutil.WriteJSONtoFile(path(binPath, fileEngineConfig), engineConfig)

	// create main Go file
	f, _ := os.Create(path(sourcePath, fileMainGo))
	fgutil.RenderTemplate(f, tplMainGoFile, projectConfig)
	f.Close()

	// create imports test Go file
	f, _ = os.Create(path(sourcePath, fileImportsGo))
	fgutil.RenderTemplate(f, tplImportsGoFile, projectConfig)
	f.Close()

	return nil
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

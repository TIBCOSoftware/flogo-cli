package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo/cli"
	"github.com/TIBCOSoftware/flogo/util"
)

var optCreate = &cli.OptionInfo{
	Name:      "create",
	UsageLine: "create AppName",
	Short:     "create a flogo project",
	Long: `Creates a flogo project.
`,
}

func init() {
	commandRegistry.RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option *cli.OptionInfo
}

func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
}

func (c *cmdCreate) Exec(args []string) error {

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Application name not specified\n\n")
		cmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	appName := args[0]

	if _, err := os.Stat(appName); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot create flogo project, directory '%s' already exists\n\n", appName)
		os.Exit(2)
	}

	if !fgutil.ExecutableExists("gb") {
		fmt.Fprintf(os.Stderr, "Error: Cannot create flogo project [%s], gb is not installed\n\n", appName)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "Creating flogo project '%s'...\n", appName)

	os.MkdirAll(appName, 0777)
	os.Chdir(appName)

	gb := fgutil.NewGb(appName)
	gb.Init(true)

	os.MkdirAll("flows", 0777)

	fmt.Fprint(os.Stdout, "Installing flogo lib...\n")

	err := gb.VendorFetch(pathFlogoLib)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// create flogo.json file
	projectConfig := &FlogoProjectConfig{
		Name:        appName,
		Version:     "0.0.1",
		Description: "My flogo application description",
		Activities:  make([]*ItemConfig, 0),
		Triggers:    make([]*ItemConfig, 0),
		Models:      make([]*ItemConfig, 0),
	}

	// todo: add default model
	// todo: make a .flogo directory in user home, were people can put a default flogo.json (use -default on create, or specify a json?)
	fgutil.WriteJSONtoFile(fileProjectConfig, projectConfig)

	createMainGoFile(gb.CodeSourcePath, projectConfig)
	createEngineEnvGoFile(gb.CodeSourcePath, projectConfig)
	createEngineConfigGoFile(gb.CodeSourcePath, projectConfig)
	createImportsGoFile(gb.CodeSourcePath, projectConfig)

	// create empty "flows" Go file
	createFlowsGoFile(gb.CodeSourcePath, make(map[string]string))

	// create config.json file
	engineConfig := DefaultEngineConfig()
	fgutil.WriteJSONtoFile(gb.NewBinFilePath(fileEngineConfig), engineConfig)

	return nil
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

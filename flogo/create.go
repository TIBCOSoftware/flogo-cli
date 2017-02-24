package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

var optCreate = &cli.OptionInfo{
	Name:      "create",
	UsageLine: "create [-flv version] AppName",
	Short:     "create a flogo project",
	Long: `Creates a flogo project.

Options:
    -flv specify the flogo-lib version
`,
}

func init() {
	commandRegistry.RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option     *cli.OptionInfo
	libVersion string
}

func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.libVersion), "flv", "", "flogo-lib version")
}

func (c *cmdCreate) Exec(args []string) error {

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "Error: Application name not specified\n\n")
		cmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
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

	err := gb.VendorFetch(pathFlogoLib, c.libVersion)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// create flogo.json file
	projectDescriptor := &FlogoProjectDescriptor{
		Name:        appName,
		Version:     "0.0.1",
		Description: "My flogo application description",
		Activities:  make([]*ItemDescriptor, 0),
		Triggers:    make([]*ItemDescriptor, 0),
		Models:      make([]*ItemDescriptor, 0),
	}

	// todo: add ability to create project from existing project descriptor

	// install default model
	fmt.Fprint(os.Stdout, "Installing default flow model...\n")
	installItem(projectDescriptor, itModel, "github.com/TIBCOSoftware/flogo-contrib/model/simple", "", false)

	createMainGoFile(gb.CodeSourcePath, projectDescriptor)
	createEngineConfigGoFile(gb.CodeSourcePath, nil)
	createImportsGoFile(gb.CodeSourcePath, projectDescriptor)

	// create empty "flows" Go file
	createFlowsGoFile(gb.CodeSourcePath, make(map[string]string))

	// create config.json file
	engineConfig := DefaultEngineConfig()
	fgutil.WriteJSONtoFile(gb.NewBinFilePath(fileEngineConfig), engineConfig)

	// create triggers.json file
	triggersConfig := DefaultTriggersConfig()
	fgutil.WriteJSONtoFile(gb.NewBinFilePath(fileTriggersConfig), triggersConfig)

	return nil
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

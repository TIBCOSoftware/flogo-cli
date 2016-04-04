package engine

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-tools/fg"
	"encoding/json"
	"os/exec"
)

var optBuild = &flogo.OptionInfo{
	Name:      "build",
	UsageLine: "build [-validate]",
	Short:     "build the engine using gb",
	Long: `Build the engine project using gb.

Options:
    -validate   validate if the engine is buildable
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option *flogo.OptionInfo
	validate bool
}

func (c *cmdBuild) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.validate), "validate", false, "only validate if buildable")
}

func (c *cmdBuild) Exec(ctx *flogo.Context, args []string) error {

	configFile, err := os.Open(fileProjectConfig)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	projectConfig := &EngineProjectConfig{}
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(projectConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse engine.json, file may be corrupted.\n\n")
		os.Exit(2)
	}

	configFile.Close()

	if len(projectConfig.Models) == 0 {
		fmt.Fprint(os.Stderr, "Error: Engine must have a least one model.\n\n")
		os.Exit(2)
	}

	if len(projectConfig.Triggers) == 0 {
		fmt.Fprint(os.Stderr, "Error: Engine must have a least one trigger.\n\n")
		os.Exit(2)
	}

	if c.validate {
		return nil
	}

	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		os.Exit(2)
	}

	return nil
}

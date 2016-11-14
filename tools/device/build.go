package device

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/config"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build",
	Short:     "build the device firmware",
	Long: `Build the device firmware.
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option     *cli.OptionInfo
	optimize   bool
	includeCfg bool
	configDir  string
}

func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {

}

func (c *cmdBuild) Exec(args []string) error {

	if len(args) != 0 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	validateDependencies()

	if PioIsProject() {
		PioBuild()
	} else {

		if DevicesEmpty() {
			fmt.Fprint(os.Stderr, "Error: No devices to build.\n\n")
			os.Exit(2)
		}

		config.LoadProjectDescriptor()
		triggersConfig := config.LoadTriggersConfig()

		workingDir, _ := os.Getwd()

		for  _, trigger := range triggersConfig.Triggers {

			if trigger.Type == "device" {

				if PioDirIsProject("devices/"+trigger.Name) {
					os.Chdir(workingDir + "/devices/" + trigger.Name)
					PioBuild()
				} else {
					fmt.Fprintf(os.Stdout, "Warning: Device Trigger %s has not been prepared.\n", trigger.Name)
				}

				os.Chdir(workingDir)
			}
		}
	}

	return nil
}

func DevicesEmpty() bool {

	if _, err := os.Stat("devices"); os.IsNotExist(err) {
		return true
	}

	empty,_ := fgutil.IsDirectoryEmpty("devices");

	return empty;
}


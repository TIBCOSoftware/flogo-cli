package device

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-cli/cli"
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

	//if cwd has platformio.ini
	//   platformio run

	//else
	// load triggers file
	// load project file

	//  validate prepare step
	//  for each trigger
	//    if "device" trigger
	//      check if devices/trigger/platformio.ini file exists
    //
	//  for each trigger
	//    if "device" trigger
	//      cd devices/trigger
	//      platformio run

	return nil
}


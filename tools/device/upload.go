package device

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"os"
	"fmt"
)

var optUpload = &cli.OptionInfo{
	Name:      "upload",
	UsageLine: "upload",
	Short:     "upload the device firmware",
	Long: `Upload the device firmware.
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdUpload{option: optUpload})
}

type cmdUpload struct {
	option     *cli.OptionInfo
	optimize   bool
	includeCfg bool
	configDir  string
}

func (c *cmdUpload) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdUpload) AddFlags(fs *flag.FlagSet) {

}

func (c *cmdUpload) Exec(args []string) error {

	if len(args) != 0 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	validateDependencies()

	//todo support named triggers, for now only allow upload from device/trigger directory

	if !PioIsProject() {
		fmt.Fprint(os.Stderr, "Error: upload can only be run within the device trigger directory\n\n")
		os.Exit(2)
	}

	PioUpload()

	return nil
}


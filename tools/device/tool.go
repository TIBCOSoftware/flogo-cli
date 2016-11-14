package device

import (
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"fmt"
	"os"
)

var optDevice = &cli.OptionInfo{
	IsTool:    true,
	Name:      "device",
	UsageLine: "device <command>",
	Short:     "tool to manage project devices",
	Long:      "Tool for managing project devices.",
}

var deviceTool *cli.Tool

// Tool gets or create the device tool
func Tool() *cli.Tool {
	if deviceTool == nil {
		deviceTool = cli.NewTool(optDevice)
		cli.RegisterTool(deviceTool)
	}

	return deviceTool
}

func init() {
	Tool()
}

func validateDependencies() {

	if !PioInstalled(){
		fmt.Fprint(os.Stderr, "Error: platformio not installed on your system\n\n")
		os.Exit(2)
	}
}
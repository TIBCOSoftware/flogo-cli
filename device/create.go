package device

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"path"
)

var optCreate = &cli.OptionInfo{
	Name:      "create",
	UsageLine: "create device",
	Short:     "create a device project",
	Long: `Creates a flogo device project.

Options:
    -f       specify the device.json to create device project from
 `,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option   *cli.OptionInfo
	fileName string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.fileName), "f", "", "flogo device file")
}

// Exec implementation of cli.Command.Exec
func (c *cmdCreate) Exec(args []string) error {

	var deviceJson string
	var deviceName string
	var err error

	if c.fileName != "" {

		if fgutil.IsRemote(c.fileName) {

			deviceJson, err = fgutil.LoadRemoteFile(c.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading device file '%s' - %s\n\n", c.fileName, err.Error())
				os.Exit(2)
			}
		} else {
			deviceJson, err = fgutil.LoadLocalFile(c.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Error loading device file '%s' - %s\n\n", c.fileName, err.Error())
				os.Exit(2)
			}

			if len(args) != 0 {
				deviceName = args[0]
			}
		}
	} else {
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "Error: Device name not specified\n\n")
			cmdUsage(c)
		}

		if len(args) != 1 {
			fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
			cmdUsage(c)
		}

		deviceName = args[0]
		deviceJson = tplSimpleDevice
	}

	currentDir, err := os.Getwd()

	if err != nil {
		return err
	}

	deviceDir := path.Join(currentDir, deviceName)

	return CreateDevice(SetupNewProjectEnv(), deviceJson, deviceDir, deviceName)
}

var tplSimpleDevice =`{
  "name": "mydevice",
  "type": "flogo:device",
  "version": "0.0.1",
  "description": "My flogo device application description",
  "device_profile": "github.com/TIBCOSoftware/flogo-contrib/device/profile/feather_m0_wifi",
  "mqtt_enabled":true,
  "settings": {
    "mqtt:server":"192.168.1.50",
    "mqtt:port":"1883",
    "mqtt:user":"",
    "mqtt:pass":"",
    "wifi:ssid":"mynetwork",
    "wifi:password": "mypass"
  },
  "triggers": [
    {
      "id": "mqtt_trigger",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/device-mqtt",
      "actionId": "pin_on",
      "settings": {
        "topic": "recievetopic"
      }
    }
  ],
  "actions": [
    {
      "id": "pin_on",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/action/device-activity",
      "data": {
        "activity": {
          "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/device-pin",
          "settings": {
            "pin": "A1",
            "digital": "true",
            "value": "HIGH"
          }
        }
      }
    }
  ]
}`
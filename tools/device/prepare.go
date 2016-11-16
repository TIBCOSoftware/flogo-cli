package device

import (
	"flag"
	"fmt"
	"os"
	"text/template"
	"strings"
	"io"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/config"
	"io/ioutil"
)

var optPrepare = &cli.OptionInfo{
	Name:      "prepare",
	UsageLine: "prepare",
	Short:     "prepare the device code",
	Long: `Prepare the device code.
`,
}


func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdPrepare{option: optPrepare})
}

type cmdPrepare struct {
	option     *cli.OptionInfo
	optimize   bool
	includeCfg bool
	configDir  string
}

func (c *cmdPrepare) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdPrepare) AddFlags(fs *flag.FlagSet) {

}

func (c *cmdPrepare) Exec(args []string) error {

	if len(args) != 0 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	validateDependencies()

	descriptor := config.LoadProjectDescriptor();
	triggersConfig := config.LoadTriggersConfig()

	workingDir, _ := os.Getwd()

	os.Mkdir("devices", 0777);

	for  _, trigger := range triggersConfig.Triggers {

		if trigger.Type == "device" {

			dirName := path(workingDir,"devices",trigger.Name)

			fmt.Printf("PioDir File: %s\n", dirName)

			if !PioDirIsProject(dirName) {
				os.Mkdir(dirName, 0777);
				os.Chdir(dirName)

				deviceType := trigger.Settings["device:type"]

				triggerSourcePath := path(workingDir,findTriggerSourcePath(descriptor, trigger.Name))
				devicesConfig := LoadDevicesConfig(triggerSourcePath)

				var device *DeviceConfig

				for  _, deviceConfig := range devicesConfig.Devices {

					fmt.Printf("Device: %v\n", deviceConfig)

					if deviceConfig.Type == deviceType {
						device = deviceConfig
						break
					}
				}

				if device == nil {
					fmt.Fprintf(os.Stderr, "Error: device type [%s] not supported\n\n", deviceType)
					os.Exit(2)
				}

				PioInit(device.Board)

				epSettings := make([]map[string]string, len(trigger.Endpoints))

				//var epSettings = [len(trigger.Endpoints)]map[string]string{}

				for  i, endpoint := range trigger.Endpoints {
					epSettings[i] = endpoint.Settings
				}

				settingsConfig := &SettingsConfig{Settings:trigger.Settings, EndpointSettings:epSettings}

				createSource(triggerSourcePath, path(dirName, "src"), device, settingsConfig)

				for  _, libConfig := range devicesConfig.Libs {

					fmt.Fprintf(os.Stdout, "Installing Lib: %v\n", libConfig)
					PioInstallLib(libConfig.ID)
				}

			} else {
				fmt.Fprintf(os.Stdout, "Warning: Device Trigger %s has already been prepared.\n", trigger.Name)
			}

			os.Chdir(workingDir)
		}
	}

	return nil
}

func createSource(triggerSourcePath string, devicePath string, deviceConfig *DeviceConfig, settings *SettingsConfig) {

	f, _ := os.Create(path(devicePath, deviceConfig.Source))
	RenderFileTemplate(f, path(triggerSourcePath,deviceConfig.Template), settings)
	f.Close()
}

func findTriggerSourcePath(descriptor *config.FlogoProjectDescriptor, triggerName string) string {

	var triggerPath string;

	for  _, trigger := range descriptor.Triggers {
		if trigger.Name == triggerName {
			triggerPath = "vendor/src/" + trigger.Path
			break
		}
	}

	return triggerPath
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

//RenderFileTemplate renders the specified template
func RenderFileTemplate(w io.Writer, templateFile string, data interface{}) {

	if (!fileExists(templateFile)) {
		fmt.Fprint(os.Stderr, "Error: template file not found\n\n")
		os.Exit(2)
	}

	t := template.New("source")
	t.Funcs(DeviceFuncMap)

	fmt.Printf("Template File: %s\n", templateFile)

	b, err := ioutil.ReadFile(templateFile)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: unable to read template file\n\n")
		os.Exit(2)
	}
	s := string(b)

	t.Parse(s)

	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
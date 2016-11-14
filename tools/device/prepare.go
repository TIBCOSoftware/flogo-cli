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

			dirName := "devices/"+trigger.Name

			if !PioDirIsProject(dirName) {
				os.Mkdir(dirName, 0777);
				os.Chdir(dirName)

				boardName := trigger.Settings["device:board"]

				triggerSourcePath := path(workingDir,findTriggerSourcePath(descriptor, trigger.Name))
				devicesConfig := LoadDevicesConfig(triggerSourcePath)

				var device *DeviceConfig

				for  _, deviceConfig := range devicesConfig.Devices {

					if deviceConfig.Board == boardName {
						device = deviceConfig
						break
					}
				}

				if device == nil {
					fmt.Fprint(os.Stderr, "Error: device [%s] not supported\n\n")
					os.Exit(2)
				}

				PioInit(boardName)
				createSource(triggerSourcePath, path(dirName, "src"), device, trigger.Settings)

				for  _, libConfig := range devicesConfig.Libs {

					PioInstallLib(libConfig.ID)
				}

			} else {
				fmt.Fprintf(os.Stdout, "Warning: Device Trigger %s has not been prepared.\n", trigger.Name)
			}

			os.Chdir(workingDir)
		}
	}

	return nil
}

func createSource(triggerSourcePath string, devicePath string, deviceConfig DeviceConfig, settings map[string]string) {

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

	t := template.New("source")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	t.ParseFiles(templateFile)

	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

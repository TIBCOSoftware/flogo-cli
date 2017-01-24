package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build [-o][-i][-c configDir][-f appFile]",
	Short:     "build the flogo application",
	Long: `Build the flogo application.

Options:
    -o   optimize for embedded flows
    -i   incorporate config into application
    -c   specifiy configration directory
    -f   Application configration file 
`,
}

const fileFlowsGo string = "flows.go"

func init() {
	commandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option     *cli.OptionInfo
	optimize   bool
	includeCfg bool
	configDir  string
	appFile    string
}

func (c *cmdBuild) OptionInfo() *cli.OptionInfo {

	return c.option
}

func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.optimize), "o", false, "optimize build")
	fs.BoolVar(&(c.includeCfg), "i", false, "include config")
	fs.StringVar(&(c.configDir), "c", "bin", "config directory")
	fs.StringVar(&(c.appFile), "f", "", "application file")
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *cmdBuild) Exec(args []string) error {

	if len(c.appFile) > 0 {
		if Exists(c.appFile) {
			fgutil.CopyFile(c.appFile, "flogo.json")
		} else {
			fmt.Fprint(os.Stderr, "Error: Invalid application configuration file.\n\n", c.appFile)
			os.Exit(2)
		}
	}

	projectDescriptor := loadProjectDescriptor()

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	gb := fgutil.NewGb(projectDescriptor.Name)

	if len(c.appFile) > 0 {
		
		os.MkdirAll(projectDescriptor.Name, 0777)
		os.Chdir(projectDescriptor.Name)
		gb.Init(true)
		err := gb.VendorFetch(pathFlogoLib, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		for _, triggerConfig := range projectDescriptor.Triggers {
			err := gb.VendorFetch(triggerConfig.Ref+"/runtime", "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}

		}

		for _, actionConfig := range projectDescriptor.Actions {
			err := gb.VendorFetch(actionConfig.Ref, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}
			for _, taskConfig := range actionConfig.Data.Flows.RootTask.Tasks {
				err := gb.VendorFetch(taskConfig.Ref+"/runtime", "")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(2)
				}
			}

		}
		createMainGoFile(gb.CodeSourcePath, projectDescriptor)
	}

	if len(projectDescriptor.Triggers) == 0 {
		fmt.Fprint(os.Stderr, "Error: Project must have a least one trigger.\n\n")
		os.Exit(2)
	}

	createImportsGoFile(gb.CodeSourcePath, projectDescriptor)

	if c.includeCfg {

		engineCfg, err := ioutil.ReadFile(filepath.Join(c.configDir, fileEngineConfig))

		if err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to read engine.config -\n%s\n", err.Error())
			os.Exit(2)
		}

		configInfo := &ConfigInfo{Include: true, ConfigJSON: string(engineCfg)}

		createEngineConfigGoFile(gb.CodeSourcePath, configInfo)

	} else {
		createEngineConfigGoFile(gb.CodeSourcePath, nil)
	}

	err := gb.Build()
	if err != nil {
		os.Exit(2)
	}

	return nil
}

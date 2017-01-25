package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

var optBuildApp = &cli.OptionInfo{
	Name:      "build-app",
	UsageLine: "build-app [-i][-c configDir][-flv version][-f application]",
	Short:     "Build new flogo application",
	Long: `Build new flogo application.
Options:
    -i   incorporate engine config into application
    -c   specifiy configration directory
    -flv specify the flogo-lib version
    -f   specify application
`,
}


func init() {
	commandRegistry.RegisterCommand(&cmdBuildApp{option: optBuildApp})
}

type cmdBuildApp struct {
	option     *cli.OptionInfo
	includeCfg bool
	configDir  string
	flvVersion string
	appFile    string
}

func (c *cmdBuildApp) OptionInfo() *cli.OptionInfo {

	return c.option
}

func (c *cmdBuildApp) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.includeCfg), "i", false, "include config")
	fs.StringVar(&(c.configDir), "c", "bin", "config directory")
	fs.StringVar(&(c.flvVersion), "flv", "", "flogo-lib version")
	fs.StringVar(&(c.appFile), "f", "", "application")
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Remote(appPath string) bool {
	return strings.HasPrefix(appPath, "http")
}

func Local(appPath string) bool {
	return appPath == "flogo.json"
}

func (c *cmdBuildApp) Exec(args []string) error {


	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}
	
	if len(c.appFile) > 0 {
		if Remote(c.appFile) {
			fgutil.CopyRemoteFile(c.appFile, "flogo.json")
		} else if Local(c.appFile) == false {
			if Exists(c.appFile) {
				fgutil.CopyFile(c.appFile, "flogo.json")
			} else {
				fmt.Fprint(os.Stderr, "Error: Invalid application configuration file.\n\n", c.appFile)
				os.Exit(2)
			}
		}
	}

	projectDescriptor := loadAppDescriptor()

	gb := fgutil.NewGb(projectDescriptor.Name)

	if len(c.appFile) > 0 {
		
		os.MkdirAll(projectDescriptor.Name, 0777)
		os.Chdir(projectDescriptor.Name)
		gb.Init(true)
		err := gb.VendorFetch(pathFlogoLib, c.flvVersion)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		for _, triggerConfig := range projectDescriptor.Triggers {
			gb.VendorDeleteSilent(triggerConfig.Ref)
			err := gb.VendorFetch(triggerConfig.Ref, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}

		}
//
		for _, actionConfig := range projectDescriptor.Actions {
			gb.VendorDeleteSilent(actionConfig.Ref)
			err := gb.VendorFetch(actionConfig.Ref, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}
			
			gb.VendorDeleteSilent(actionConfig.Data.Ref)
			err = gb.VendorFetch(actionConfig.Data.Ref, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}
			for _, taskConfig := range actionConfig.Data.RootTask.Tasks {
				gb.VendorDeleteSilent(taskConfig.Ref)
				err := gb.VendorFetch(taskConfig.Ref, "")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(2)
				}
			}

		}
		createNewMainGoFile(gb.CodeSourcePath, projectDescriptor)
	}

	if len(projectDescriptor.Triggers) == 0 {
		fmt.Fprint(os.Stderr, "Error: Project must have a least one trigger.\n\n")
		os.Exit(2)
	}

	createNewImportsGoFile(gb.CodeSourcePath, projectDescriptor)

	if c.includeCfg {

		engineCfg, err := ioutil.ReadFile(filepath.Join(c.configDir, fileEngineConfig))

		if err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to read engine.config -\n%s\n", err.Error())
			os.Exit(2)
		}

		configInfo := &ConfigInfo{Include: true, ConfigJSON: string(engineCfg)}

		createNewEngineConfigGoFile(gb.CodeSourcePath, configInfo)

	} else {
		createNewEngineConfigGoFile(gb.CodeSourcePath, nil)
	}

	err := gb.Build()
	if err != nil {
		os.Exit(2)
	}

	return nil
}

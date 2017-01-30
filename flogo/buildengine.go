package main

import (
	"flag"
	"fmt"
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var optBuildApp = &cli.OptionInfo{
	Name:      "build-engine",
	UsageLine: "build-engine [-app application]",
	Short:     "Build engine based on flogo application",
	Long: `Build Flogo engine based on flogo application.
Options:
    -flv specify the flogo-lib version
    -app specify application location
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
	ctbVersion string
	appFile    string
	localDep   string
}

func (c *cmdBuildApp) OptionInfo() *cli.OptionInfo {

	return c.option
}

func (c *cmdBuildApp) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.includeCfg), "i", false, "include config")
	fs.StringVar(&(c.configDir), "c", "bin", "config directory")
	fs.StringVar(&(c.flvVersion), "flv", "", "flogo-lib version")
	fs.StringVar(&(c.ctbVersion), "cv", "", "contrib version")
	fs.StringVar(&(c.appFile), "app", "", "application")
	fs.StringVar(&(c.localDep), "d", "", "copy dependencies from directory")
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
	os.MkdirAll(projectDescriptor.Name, 0777)
	os.Chdir(projectDescriptor.Name)
	gb.Init(true)

	if len(c.localDep) > 0 {
	     if Exists(c.localDep) {
				fgutil.CopyDir(c.localDep, gb.VendorPath)
			} else {
				fmt.Fprint(os.Stderr, "Error: Invalid dependency location.\n\n", c.localDep)
				os.Exit(2)
			}
	} else {
		gb.VendorDeleteSilent(pathFlogoLib)
		err := gb.VendorFetch(pathFlogoLib, c.flvVersion)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		for _, triggerConfig := range projectDescriptor.Triggers {
			gb.VendorDeleteSilent(triggerConfig.Ref)
			err = gb.VendorFetch(triggerConfig.Ref, c.ctbVersion)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}
		}
		//
		for _, actionConfig := range projectDescriptor.Actions {
			gb.VendorDeleteSilent(actionConfig.Ref)
			err := gb.VendorFetch(actionConfig.Ref, c.flvVersion)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}

			gb.VendorDeleteSilent(actionConfig.Data.Flow.Ref)
			err = gb.VendorFetch(actionConfig.Data.Flow.Ref, c.flvVersion)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(2)
			}
			
			for _, taskConfig := range actionConfig.Data.Flow.RootTask.Tasks {
				gb.VendorDeleteSilent(taskConfig.Ref)
				err = gb.VendorFetch(taskConfig.Ref, c.ctbVersion)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(2)
				}
			}
		}
	}
	createNewMainGoFile(gb.CodeSourcePath, projectDescriptor)
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

	if len(c.localDep) > 0 {

	}

	err := gb.Build()
	if err != nil {
		os.Exit(2)
	}

	return nil
}

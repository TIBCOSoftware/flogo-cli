package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/flogo/cli"
	"github.com/TIBCOSoftware/flogo/util"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build [-o][-i]",
	Short:     "build the flogo application",
	Long: `Build the flogo application.

Options:
    -o   optimize for embedded flows
    -i   incorporate config into application
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
}

func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.optimize), "o", false, "optimize build")
	fs.BoolVar(&(c.includeCfg), "i", false, "include config")
}

func (c *cmdBuild) Exec(args []string) error {

	projectDescriptor := loadProjectDescriptor()

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	gb := fgutil.NewGb(projectDescriptor.Name)

	flows := ImportFlows(projectDescriptor, dirFlows)
	createFlowsGoFile(gb.CodeSourcePath, flows)

	if len(projectDescriptor.Models) == 0 {
		fmt.Fprint(os.Stderr, "Error: Project must have a least one model.\n\n")
		os.Exit(2)
	}

	if len(projectDescriptor.Triggers) == 0 {
		fmt.Fprint(os.Stderr, "Error: Project must have a least one trigger.\n\n")
		os.Exit(2)
	}

	//allFlowExprs := getAllFlowExprs(dirFlows)
	//fmt.Printf("all flow exprs: %v\n", allFlowExprs)
	//
	//allFlowTransExprs := make(map[string]map[int]string)
	//
	//if len(allFlowExprs) > 0 {
	//
	//	for flowURI, exprs := range allFlowExprs {
	//
	//		transExprs := convertExprsToGo(exprs)
	//		allFlowTransExprs[flowURI] = transExprs
	//	}
	//}
	//
	//createExprsGoFile(gb.CodeSourcePath, allFlowTransExprs)

	if c.optimize {

		//todo optimize triggers

		activityTypes := getAllActivityTypes(dirFlows)

		var activities []*ItemDescriptor

		for  _, activity := range projectDescriptor.Activities {

			if _, ok := activityTypes[activity.Name]; ok {
				activities = append(activities, activity)
			}
		}

		projectDescriptor.Activities = activities

		createImportsGoFile(gb.CodeSourcePath, projectDescriptor)
	}

	if c.includeCfg {

		engineCfg, err := ioutil.ReadFile(filepath.Join("bin", fileEngineConfig))

		if err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to read engine.config -\n%s\n", err.Error())
			os.Exit(2)
		}

		triggersCfg, err := ioutil.ReadFile(filepath.Join("bin", fileTriggersConfig))

		if err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to read triggers.config -\n%s\n", err.Error())
			os.Exit(2)
		}

		configInfo := &ConfigInfo{Include:true, ConfigJSON:string(engineCfg), TriggerJSON:string(triggersCfg)}

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


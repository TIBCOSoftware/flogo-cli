package device

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build",
	Short:     "build the device firmware",
	Long: `Build the device firmware.
`,
}

func init() {
	commandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
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
	fs.BoolVar(&(c.optimize), "o", false, "optimize build")
	fs.BoolVar(&(c.includeCfg), "i", false, "include config")
	fs.StringVar(&(c.configDir), "c", "bin", "config directory")
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

	//if len(projectDescriptor.Triggers) == 0 {
	//	fmt.Fprint(os.Stderr, "Error: Project must have a least one trigger.\n\n")
	//	os.Exit(2)
	//}

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

		//optimize activities

		activityTypes := getAllActivityTypes(dirFlows)

		var activities []*ItemDescriptor

		for  _, activity := range projectDescriptor.Activities {

			if _, ok := activityTypes[activity.Name]; ok {
				activities = append(activities, activity)
			}
		}

		projectDescriptor.Activities = activities

		//optimize triggers

		triggersConfigPath := gb.NewBinFilePath(fileTriggersConfig)
		triggersConfigFile, err := os.Open(triggersConfigPath)

		triggersConfig := &TriggersConfig{}
		jsonParser := json.NewDecoder(triggersConfigFile)

		if err = jsonParser.Decode(triggersConfig); err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to parse application triggers.json, file may be corrupted.\n\n")
			os.Exit(2)
		}

		triggersConfigFile.Close()

		if triggersConfig.Triggers == nil {
			triggersConfig.Triggers = make([]*TriggerConfig, 0)
		}

		var triggers []*ItemDescriptor

		for  _, trigger := range projectDescriptor.Triggers {

			if ContainsTriggerConfig(triggersConfig.Triggers, trigger.Name)  {
				triggers = append(triggers, trigger)
			}
		}

		projectDescriptor.Triggers = triggers

		createImportsGoFile(gb.CodeSourcePath, projectDescriptor)
	}

	if c.includeCfg {

		engineCfg, err := ioutil.ReadFile(filepath.Join(c.configDir, fileEngineConfig))

		if err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to read engine.config -\n%s\n", err.Error())
			os.Exit(2)
		}

		triggersCfg, err := ioutil.ReadFile(filepath.Join(c.configDir, fileTriggersConfig))

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


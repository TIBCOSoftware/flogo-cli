package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TIBCOSoftware/flogo-tools/fg"
	"github.com/TIBCOSoftware/flogo-tools/fgutil"
)

type getItems func(cfg *EngineConfig) []*ItemConfig

// AddEngineItem adds an item(activity, model or trigger) to the engine
func AddEngineItem(c flogo.Command, itemType string, args []string, gi getItems, useSrc bool) (itemConfig *ItemConfig, engineConfig *EngineConfig) {

	configFile, err := os.Open(fileDescriptor)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: %s path not specified\n\n", fgutil.Capitalize(itemType))
		Tool().CmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	engineConfig = &EngineConfig{}
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(engineConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse engine.json, file maybe corrupted.\n\n")
		os.Exit(2)
	}

	configFile.Close()

	itemPath := args[0]

	if ContainsItem(itemPath, engineConfig.Models) {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' is already in engine project.\n\n", fgutil.Capitalize(itemType), itemPath)
		os.Exit(2)
	}

	var altPath string
	var localPath string

	if strings.HasPrefix(itemPath, "local://") {
		localPath = itemPath[8:]
		altPath = "file://" + localPath
	} else if strings.HasPrefix(itemPath, "file://") {
		localPath = itemPath[6:]
		altPath = "local://" + localPath
	}

	if len(localPath) > 0 {

		if ContainsItem(altPath, gi(engineConfig)) {
			fmt.Fprintf(os.Stderr, "Error: %s '%s' is already in engine project.\n\n", fgutil.Capitalize(itemType), itemPath)
			os.Exit(2)
		}

		usesGb := false

		itemFile, err := os.Open(path(localPath, itemType+".json"))

		if err != nil {
			itemFile.Close()
			itemFile, err = os.Open(path(localPath, "src", itemType+".json"))

			usesGb = true
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a flogo-based %s.\n\n", itemPath, itemType)
				itemFile.Close()
				os.Exit(2)
			}
		}

		itemConfig := &struct {
			Name string `json:"name"`
		}{}

		jsonParser := json.NewDecoder(itemFile)

		if err = jsonParser.Decode(itemConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to parse %s.json, file maybe corrupted.\n\n", itemType)
			itemFile.Close()
			os.Exit(2)
		}

		itemFile.Close()

		sourcePath := path("src")
		vendorPath := path("vendor", "src")

		toDir := path(vendorPath, itemType, itemConfig.Name)

		if useSrc {
			toDir = path(sourcePath, itemType, itemConfig.Name)
		}

		fromDir := localPath

		if usesGb {
			fromDir = path(localPath, "src")
		}

		fgutil.CopyDir(fromDir, toDir)

		if usesGb {

			fgutil.CopyDir(path(localPath, "src", itemType), path("src", itemType, itemConfig.Name))
		}

	} else {

		cmd := exec.Command("gb", "vendor", "fetch", itemPath)
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		//check if it contains model.json
	}

	//engineConfig.Models = append(engineConfig.Models, itemConfig)
	//fgutil.WriteJsonToFile(fileDescriptor, engineConfig)

	return &ItemConfig{Path: itemPath, Version: "latest"}, engineConfig
}


func updateConfigFiles(engineConfig *EngineConfig) {
	fgutil.WriteJSONtoFile(fileDescriptor, engineConfig)

	sourcePath := path("src", engineConfig.Name)

	// create imports test Go file
	f, _ := os.Create(path(sourcePath, fileImportsGo))
	fgutil.RenderTemplate(f, tplImportsGoFile, engineConfig)
	f.Close()
}
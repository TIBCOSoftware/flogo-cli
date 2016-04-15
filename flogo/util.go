package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo/util"
)

func updateProjectConfigFiles(gb *fgutil.Gb, projectConfig *FlogoProjectConfig) {
	fgutil.WriteJSONtoFile(fileProjectConfig, projectConfig)
	createImportsGoFile(gb.CodeSourcePath, projectConfig)
}

func loadProjectConfig() *FlogoProjectConfig {

	projectConfigFile, err := os.Open(fileProjectConfig)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	projectConfig := &FlogoProjectConfig{}
	jsonParser := json.NewDecoder(projectConfigFile)

	if err = jsonParser.Decode(projectConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse flogo.json, file may be corrupted.\n\n")
		os.Exit(2)
	}

	projectConfigFile.Close()

	return projectConfig
}

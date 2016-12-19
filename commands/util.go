package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

func updateProjectFiles(gb *fgutil.Gb, projectDescriptor *FlogoProjectDescriptor) {
	fgutil.WriteJSONtoFile(fileDescriptor, projectDescriptor)
	createImportsGoFile(gb.CodeSourcePath, projectDescriptor)
}

func loadProjectDescriptor() *FlogoProjectDescriptor {

	projectDescriptorFile, err := os.Open(fileDescriptor)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	projectDescriptor := &FlogoProjectDescriptor{}
	jsonParser := json.NewDecoder(projectDescriptorFile)

	if err = jsonParser.Decode(projectDescriptor); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to parse flogo.json, file may be corrupted.\n - %s\n", err.Error())
		os.Exit(2)
	}

	projectDescriptorFile.Close()

	return projectDescriptor
}

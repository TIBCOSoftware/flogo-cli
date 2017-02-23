package app

import (
	"encoding/json"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/flogo-cli/env"
	"strings"
)

// BuildPreProcessor interface for build pre-processors
type BuildPreProcessor interface {
	PrepareForBuild(env env.Project)
}

// CreateApp creates an application from the specified json application descriptor
func CreateApp(env env.Project, appJson string) error {

	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()

	if err != nil {
		return err
	}

	appDir := fgutil.Path(currentDir, descriptor.Name)
	env.Init(appDir)
	err = env.Create(false)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromString(fgutil.Path(appDir,"flogo.json"), appJson)
	if err != nil {
		return err
	}

	//
	//todo allow ability to specify flogo-lib version
	env.InstallDependency(pathFlogoLib, "")

	refs := ExtractRefs(descriptor)

	for _, ref := range refs {
		path, version :=  splitVersion(ref)
		err = env.InstallDependency(path, version)
		if err != nil {
			return err
		}
	}

	// create source files
	cmdPath := fgutil.Path(env.GetSourceDir(), strings.ToLower(descriptor.Name))
	os.MkdirAll(cmdPath, 0777)

	createMainGoFile(cmdPath)
	createImportsGoFile(cmdPath, refs)

	return nil
}

// BuildApp build the flogo application
func BuildApp(env env.Project, customPreProcessor BuildPreProcessor) (err error) {

	if customPreProcessor != nil {
		customPreProcessor.PrepareForBuild(env)
	}

	//todo do standard pre-processing
	// regenerate imports?

	err = env.Build()
	if err != nil {
		return err
	}

	err = fgutil.MoveFiles(env.GetBinDir(), env.GetRootDir())
	if err != nil {
		return err
	}

	return
}

// InstallDependency install a dependency
func InstallDependency(env env.Project, path string, version string) error {

	return env.InstallDependency(path, version)
}

// ParseAppDescriptor parse the application descriptor
func ParseAppDescriptor(appJson string) (*FlogoAppDescriptor, error) {
	descriptor := &FlogoAppDescriptor{}

	err := json.Unmarshal([]byte(appJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor,nil
}
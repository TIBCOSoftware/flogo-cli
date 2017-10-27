package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/dep"
	"github.com/TIBCOSoftware/flogo-cli/env"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/flogo-cli/config"
	"os/exec"
)

// BuildPreProcessor interface for build pre-processors
type BuildPreProcessor interface {
	PrepareForBuild(env env.Project) error
}

// CreateApp creates an application from the specified json application descriptor
func CreateApp(env env.Project, appJson , appDir , appName , vendorDir string) error {
	if IsBuildExperimental(){
		return doCreate(env, appJson, appDir, appName, vendorDir)
	}


	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	if appName != "" {
		// override the application name

		altJson := strings.Replace(appJson, `"`+descriptor.Name+`"`, `"`+appName+`"`, 1)
		altDescriptor, err := ParseAppDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Name == appName {
			appJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(appJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = appName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			appJson = string(updApp)
		}

		descriptor.Name = appName
	} else {
		appName = descriptor.Name
		appDir = path.Join(appDir, appName)
	}

	env.Init(appDir)
	err = env.Create(false, vendorDir)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromString(path.Join(appDir, "flogo.json"), appJson)
	if err != nil {
		return err
	}

	//todo allow ability to specify flogo-lib version
	env.InstallDependency(pathFlogoLib, "")

	deps := config.ExtractDependencies(descriptor)

	for _, dep := range deps {
		path, version := splitVersion(dep.Ref)
		err = env.InstallDependency(path, version)
		if err != nil {
			return err
		}
	}

	// create source files
	cmdPath := path.Join(env.GetSourceDir(), strings.ToLower(descriptor.Name))
	os.MkdirAll(cmdPath, 0777)

	createMainGoFile(cmdPath, "")
	createImportsGoFile(cmdPath, deps)

	return nil
}

// doCreate performs the app creation
func doCreate(enviro env.Project, appJson , rootDir , appName , vendorDir string) error{
	fmt.Printf("Creating initial project structure, this migh take a few seconds ... \n")
	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}
	if appName != "" {
		// override the application name

		altJson := strings.Replace(appJson, `"`+descriptor.Name+`"`, `"`+appName+`"`, 1)
		altDescriptor, err := ParseAppDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Name == appName {
			appJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(appJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = appName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			appJson = string(updApp)
		}

		descriptor.Name = appName
	} else {
		appName = descriptor.Name
		rootDir = path.Join(rootDir, appName)
	}

	err = enviro.Init(rootDir)
	if err != nil {
		return err
	}

	err = enviro.Create(false, vendorDir)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromString(path.Join(rootDir, "flogo.json"), appJson)
	if err != nil {
		return err
	}

	deps := config.ExtractDependencies(descriptor)
	// create source files
	appDir := path.Join(enviro.GetSourceDir(), descriptor.Name)
	os.MkdirAll(appDir, os.ModePerm)

	createMainGoFile(appDir, "")
	createImportsGoFile(appDir, deps)

	// Create the dep manager
	depManager := &dep.DepManager{AppDir:appDir}
	// Initialize the dep manager
	err = depManager.Init(rootDir)
	if err != nil {
		return err
	}

	return nil
}

type PrepareOptions struct {
	PreProcessor    BuildPreProcessor
	OptimizeImports bool
	EmbedConfig     bool
	Shim            string
}

// PrepareApp do all pre-build setup and pre-processing
func PrepareApp(env env.Project, options *PrepareOptions) (err error) {
	if IsBuildExperimental(){
		return doPrepare(env, options)
	}
	if options == nil {
		options = &PrepareOptions{}
	}

	if options.PreProcessor != nil {
		err = options.PreProcessor.PrepareForBuild(env)
		if err != nil {
			return err
		}
	}

	//generate metadata
	err = generateGoMetadata(env)
	if err != nil {
		return err
	}

	//load descriptor
	appJson, err := fgutil.LoadLocalFile(path.Join(env.GetRootDir(), "flogo.json"))

	if err != nil {
		return err
	}
	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	//generate imports file
	var deps []*config.Dependency

	if options.OptimizeImports {

		deps = config.ExtractDependencies(descriptor)

	} else {
		deps, err = ListDependencies(env, 0)
	}

	cmdPath := path.Join(env.GetSourceDir(), strings.ToLower(descriptor.Name))
	createImportsGoFile(cmdPath, deps)

	removeEmbeddedAppGoFile(cmdPath)

	wasUsingShim := fgutil.FileExists(path.Join(cmdPath, fileShimGo))
	removeShimGoFiles(cmdPath)

	if options.Shim != "" {

		removeMainGoFile(cmdPath) //todo maybe rename if it exists
		createShimSupportGoFile(cmdPath, appJson, options.EmbedConfig)

		fmt.Println("Shim:", options.Shim)

		for _, value := range descriptor.Triggers {

			fmt.Println("Id:", value.ID)
			if value.ID == options.Shim {
				triggerPath := path.Join(env.GetVendorSrcDir(), value.Ref, "trigger.json")

				mdJson, err := fgutil.LoadLocalFile(triggerPath)
				if err != nil {
					return err
				}
				metadata, err := ParseTriggerMetadata(mdJson)
				if err != nil {
					return err
				}

				fmt.Println("Shim Metadata:", metadata.Shim)

				if metadata.Shim != "" {

					//todo blow up if shim file not found
					shimFilePath := path.Join(env.GetVendorSrcDir(), value.Ref, dirShim, fileShimGo)
					fmt.Println("Shim File:", shimFilePath)
					fgutil.CopyFile(shimFilePath, path.Join(cmdPath, fileShimGo))

					if metadata.Shim == "plugin" {
						//look for Makefile and execute it
						makeFilePath := path.Join(env.GetVendorSrcDir(), value.Ref, dirShim, makeFile)
						fmt.Println("Make File:", makeFilePath)
						fgutil.CopyFile(makeFilePath, path.Join(cmdPath, makeFile))

						// Copy the vendor folder (Ugly workaround, this will go once our app is golang structure compliant)
						vendorDestDir := path.Join(cmdPath, "vendor")
						_, err = os.Stat(vendorDestDir)
						if err == nil {
							// We don't support existing vendor folders yet
							return fmt.Errorf("Unsupported vendor folder found for function build, please create an issue on https://github.com/TIBCOSoftware/flogo")
						}
						// Create vendor folder
						err = CopyDir(env.GetVendorSrcDir(), vendorDestDir)
						if err != nil {
							return err
						}
						defer os.RemoveAll(vendorDestDir)

						// Execute make
						cmd := exec.Command("make", "-C", cmdPath)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Env = append(os.Environ(),
							fmt.Sprintf("GOPATH=%s", env.GetRootDir()),
						)

						err = cmd.Run()
						if err != nil {
							return err
						}
					}
				}

				break
			}
		}

	} else {

		if wasUsingShim {
			createMainGoFile(cmdPath, "")
		}

		if options.EmbedConfig {
			createEmbeddedAppGoFile(cmdPath, appJson)
		}
	}

	return
}

// doPrepare performs all the prepare functionality
func doPrepare(env env.Project, options *PrepareOptions) (err error) {
	// Create the dep manager
	depManager := &dep.DepManager{AppDir:env.GetAppDir()}
	if !depManager.IsInitialized(){
		// This is an old app
		fmt.Println("Initializing dependency management files ....")
		err := depManager.Init(env.GetRootDir())
		if err != nil {
			return err
		}
		// Move old vendor folder to _old_vendor
		oldVendorDir := path.Join(env.GetRootDir(), "vendor")
		_, err = os.Stat(oldVendorDir)
		if err == nil {
			newVendorDir, _ := path.Split(env.GetVendorDir())
			newVendorDir = path.Join(newVendorDir, "_old_vendor")
			fmt.Printf("Moving old vendoring directory %s to %s \n", oldVendorDir, newVendorDir)
			// Vendor found, move it
			err = CopyDir(oldVendorDir, newVendorDir)
			if err != nil {
				return err
			}
			defer os.RemoveAll(oldVendorDir)
		}

	}

	if options == nil {
		options = &PrepareOptions{}
	}

	// Call external preprocessor
	if options.PreProcessor != nil {
		err = options.PreProcessor.PrepareForBuild(env)
		if err != nil {
			return err
		}
	}


	//generate metadata
	err = generateGoMetadata(env)
	if err != nil {
		return err
	}

	//load descriptor
	appJson, err := fgutil.LoadLocalFile(path.Join(env.GetRootDir(), "flogo.json"))

	if err != nil {
		return err
	}
	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	removeEmbeddedAppGoFile(env.GetAppDir())
	removeShimGoFiles(env.GetAppDir())

	if options.Shim != "" {

		removeMainGoFile(env.GetAppDir()) //todo maybe rename if it exists
		createShimSupportGoFile(env.GetAppDir(), appJson, options.EmbedConfig)

		fmt.Println("Shim:", options.Shim)

		for _, value := range descriptor.Triggers {

			fmt.Println("Id:", value.ID)
			if value.ID == options.Shim {
				triggerPath := path.Join(env.GetVendorSrcDir(), value.Ref, "trigger.json")

				mdJson, err := fgutil.LoadLocalFile(triggerPath)
				if err != nil {
					return err
				}
				metadata, err := ParseTriggerMetadata(mdJson)
				if err != nil {
					return err
				}

				fmt.Println("Shim Metadata:", metadata.Shim)

				if metadata.Shim != "" {

					//todo blow up if shim file not found
					shimFilePath := path.Join(env.GetVendorSrcDir(), value.Ref, dirShim, fileShimGo)
					fmt.Println("Shim File:", shimFilePath)
					fgutil.CopyFile(shimFilePath, path.Join(env.GetAppDir(), fileShimGo))

					if metadata.Shim == "plugin" {
						//look for Makefile and execute it
						makeFilePath := path.Join(env.GetVendorSrcDir(), value.Ref, dirShim, makeFile)
						fmt.Println("Make File:", makeFilePath)
						fgutil.CopyFile(makeFilePath, path.Join(env.GetAppDir(), makeFile))

						// Execute make
						cmd := exec.Command("make", "-C", env.GetAppDir())
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Env = append(os.Environ(),
							fmt.Sprintf("GOPATH=%s", env.GetRootDir()),
						)

						err = cmd.Run()
						if err != nil {
							return err
						}
					}
				}

				break
			}
		}

	} else if options.EmbedConfig {
		createEmbeddedAppGoFile(env.GetAppDir(), appJson)
	}
	return
}

type BuildOptions struct {
	*PrepareOptions

	SkipPrepare bool
}

// BuildApp build the flogo application
func BuildApp(env env.Project, options *BuildOptions) (err error) {

	if IsBuildExperimental(){
		return doBuild(env, options)
	}

	if options == nil {
		options = &BuildOptions{}
	}

	if !options.SkipPrepare {
		err = PrepareApp(env, options.PrepareOptions)

		if err != nil {
			return err
		}
	}

	err = env.Build()
	if err != nil {
		return err
	}

	if !options.EmbedConfig {
		fgutil.CopyFile(path.Join(env.GetRootDir(), config.FileDescriptor), path.Join(env.GetBinDir(), config.FileDescriptor))
		if err != nil {
			return err
		}
	} else {
		os.Remove(path.Join(env.GetBinDir(), config.FileDescriptor))
	}

	return
}

// doBuildApp performs the build functionality
func doBuild(env env.Project, options *BuildOptions) (err error) {
	if options == nil {
		options = &BuildOptions{}
	}

	if !options.SkipPrepare {
		err = PrepareApp(env, options.PrepareOptions)

		if err != nil {
			return err
		}
	}

	err = env.Build()
	if err != nil {
		return err
	}

	if !options.EmbedConfig {
		fgutil.CopyFile(path.Join(env.GetRootDir(), config.FileDescriptor), path.Join(env.GetBinDir(), config.FileDescriptor))
		if err != nil {
			return err
		}
	} else {
		os.Remove(path.Join(env.GetBinDir(), config.FileDescriptor))
	}

	return
}

// InstallPalette install a palette
func InstallPalette(env env.Project, path string) error {

	var file []byte

	file, _ = ioutil.ReadFile(path)

	var paletteDescriptor *config.FlogoPaletteDescriptor
	err := json.Unmarshal(file, &paletteDescriptor)

	var deps []config.Dependency

	if err != nil {
		err = json.Unmarshal(file, &deps)
	} else {
		deps = paletteDescriptor.Extensions
	}

	if err != nil {
		return err
		//fmt.Fprint(os.Stderr, "Error: Unable to parse palette descriptor, file may be corrupted.\n\n")
		//os.Exit(2)
	}

	for _, dep := range deps {
		err = env.InstallDependency(dep.Ref, "")
		if err != nil {
			return err
		}
	}

	//fmt.Fprintf(os.Stdout, "Adding Palette: %s [%s]\n\n", paletteDescriptor.Name, paletteDescriptor.Description)

	return nil
}

// InstallDependency install a dependency
func InstallDependency(environ env.Project, path string, version string) error {
	if IsBuildExperimental(){
		// Create the dep manager
		depManager := &dep.DepManager{AppDir:environ.GetRootDir()}
		return depManager.InstallDependency(environ.GetRootDir(), environ.GetAppDir(), path, version)
	}
	return environ.InstallDependency(path, version)
}

// UninstallDependency uninstall a dependency
func UninstallDependency(environ env.Project, path string) error {
	if IsBuildExperimental(){
		// Create the dep manager
		depManager := &dep.DepManager{AppDir:environ.GetRootDir()}
		return depManager.UninstallDependency(environ.GetRootDir(), environ.GetAppDir(), path)
	}

	return environ.UninstallDependency(path)
}

// ListDependencies lists all installed dependencies
func ListDependencies(env env.Project, cType config.ContribType) ([]*config.Dependency, error) {

	vendorSrc := env.GetVendorSrcDir()
	var deps []*config.Dependency

	err := filepath.Walk(vendorSrc, func(filePath string, info os.FileInfo, _ error) error {

		if !info.IsDir() {

			switch info.Name() {
			case "action.json":
				if cType == 0 || cType == config.ACTION {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:action" {
						deps = append(deps, &config.Dependency{ContribType: config.ACTION, Ref: ref})
					}
				}
			case "trigger.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0: len(filePath)-12]
				if _, err := os.Stat(fmt.Sprintf("%s/../trigger.json", dir)); err == nil {
					//old trigger.json, ignore
					return nil
				}
				if cType == 0 || cType == config.TRIGGER {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:trigger" {
						deps = append(deps, &config.Dependency{ContribType: config.TRIGGER, Ref: ref})
					}
				}
			case "activity.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0: len(filePath)-13]
				if _, err := os.Stat(fmt.Sprintf("%s/../activity.json", dir)); err == nil {
					//old activity.json, ignore
					return nil
				}
				if cType == 0 || cType == config.ACTIVITY {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:activity" {
						deps = append(deps, &config.Dependency{ContribType: config.ACTIVITY, Ref: ref})
					}
				}
			case "flow-model.json":
				if cType == 0 || cType == config.FLOW_MODEL {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:flow-model" {
						deps = append(deps, &config.Dependency{ContribType: config.FLOW_MODEL, Ref: ref})
					}
				}
			}

		}

		return nil
	})

	return deps, err
}

func refPath(vendorSrc string, filePath string) string {

	startIdx := len(vendorSrc) + 1
	endIdx := strings.LastIndex(filePath, string(os.PathSeparator))

	return strings.Replace(filePath[startIdx:endIdx], string(os.PathSeparator), "/", -1)
}

func readDescriptor(path string, info os.FileInfo) (*config.Descriptor, error) {

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return nil, err
	}

	return ParseDescriptor(string(raw))
}

func generateGoMetadata(env env.Project) error {
	//todo optimize metadata recreation to minimize compile times
	dependencies, err := ListDependencies(env, 0)

	if err != nil {
		return err
	}

	for _, dependency := range dependencies {
		createMetadata(env, dependency)
	}

	return nil
}

func createMetadata(env env.Project, dependency *config.Dependency) error {

	vendorSrc := env.GetVendorSrcDir()
	mdFilePath := path.Join(vendorSrc, dependency.Ref)
	mdGoFilePath := path.Join(vendorSrc, dependency.Ref)
	pkg := path.Base(mdFilePath)

	tplMetadata := tplMetadataGoFile

	switch dependency.ContribType {
	case config.ACTION:
		mdFilePath = path.Join(mdFilePath, "action.json")
		mdGoFilePath = path.Join(mdGoFilePath, "action_metadata.go")
	case config.TRIGGER:
		mdFilePath = path.Join(mdFilePath, "trigger.json")
		mdGoFilePath = path.Join(mdGoFilePath, "trigger_metadata.go")
		tplMetadata = tplTriggerMetadataGoFile
	case config.ACTIVITY:
		mdFilePath = path.Join(mdFilePath, "activity.json")
		mdGoFilePath = path.Join(mdGoFilePath, "activity_metadata.go")
		tplMetadata = tplActivityMetadataGoFile
	default:
		return nil
	}

	raw, err := ioutil.ReadFile(mdFilePath)
	if err != nil {
		return err
	}

	info := &struct {
		Package      string
		MetadataJSON string
	}{
		Package:      pkg,
		MetadataJSON: string(raw),
	}

	f, _ := os.Create(mdGoFilePath)
	fgutil.RenderTemplate(f, tplMetadata, info)
	f.Close()

	return nil
}

var tplMetadataGoFile = `package {{.Package}}

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

func getJsonMetadata() string {
	return jsonMetadata
}
`

var tplActivityMetadataGoFile = `package {{.Package}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(NewActivity(md))
}
`

var tplTriggerMetadataGoFile = `package {{.Package}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.RegisterFactory(md.ID, NewFactory(md))
}
`

// ParseDescriptor parse a descriptor
func ParseDescriptor(descJson string) (*config.Descriptor, error) {
	descriptor := &config.Descriptor{}

	err := json.Unmarshal([]byte(descJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

// ParseAppDescriptor parse the application descriptor
func ParseAppDescriptor(appJson string) (*config.FlogoAppDescriptor, error) {
	descriptor := &config.FlogoAppDescriptor{}

	err := json.Unmarshal([]byte(appJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

// ParseTriggerMetadata parse the trigger metadata
func ParseTriggerMetadata(metadataJson string) (*config.TriggerMetadata, error) {
	metadata := &config.TriggerMetadata{}

	err := json.Unmarshal([]byte(metadataJson), metadata)

	if err != nil {
		return nil, err
	}

	return metadata, nil
}

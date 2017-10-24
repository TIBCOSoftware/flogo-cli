package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"bytes"
	"github.com/TIBCOSoftware/flogo-cli/config"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

type DepProject struct {
	BinDir             string
	RootDir            string
	SourceDir          string
	VendorDir          string
	VendorSrcDir       string
	CodeSourcePath     string
	AppDir             string
	FileDescriptorPath string
}

type ConstraintDef struct {
	ProjectRoot string
	Version     string
}

type DepManager struct {
	AppDir string
}

// Init initializes the dependency manager
func (b *DepManager) Init(rootDir, appDir string) error {
	exists := fgutil.ExecutableExists("dep")
	if !exists {
		return errors.New("dep not installed")
	}

	cmd := exec.Command("dep", "init")
	cmd.Dir = appDir
	newEnv := os.Environ()
	newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", rootDir))
	cmd.Env = newEnv

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	// TODO remove this prune cmd once it gets absorved into dep ensure https://github.com/golang/dep/issues/944
	cmd = exec.Command("dep", "prune")
	cmd.Dir = appDir
	cmd.Env = newEnv

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// InstallDependency installs the given dependency
func (b *DepManager) InstallDependency(rootDir, appDir, depPath , depVersion string) error {
	exists := fgutil.ExecutableExists("dep")
	if !exists {
		return errors.New("dep not installed")
	}
	fmt.Println("Validating existing dependencies, this might take a few seconds...")

	// Load imports file
	importsPath := path.Join(appDir, config.FileImportsGo)
	// Validate that it exists
	_, err := os.Stat(importsPath)

	if err != nil {
		return fmt.Errorf("Error installing dependency, import file '%s' doesn't exists", importsPath)
	}

	fset := token.NewFileSet()

	importsFileAst, err := parser.ParseFile(fset, importsPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("Error parsing import file '%s', %s", importsPath, err)
	}

	//Validate that the install does not exist in imports.go file
	for _, imp := range importsFileAst.Imports {
		if imp.Path.Value == strconv.Quote(depPath) {
			return fmt.Errorf("Error installing dependency, import '%s' already exists", depPath)
		}
	}

	existingConstraint, err := GetExistingConstraint(rootDir, appDir, depPath)
	if err != nil {
		return err
	}

	if existingConstraint != nil {
		if len(depVersion) > 0 {
			fmt.Printf("Existing root package version found '%s', to update it please change Gopkg.toml manually\n", existingConstraint.Version)
		}
	} else {
		// Contraint does not exist add it
		fmt.Printf("Adding new dependency '%s' version '%s' \n", depPath, depVersion)
		cmd := exec.Command("dep", "ensure", "-add", fmt.Sprintf("%s@%s", depPath, depVersion))
		cmd.Dir = appDir
		newEnv := os.Environ()
		newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", rootDir))
		cmd.Env = newEnv

		// Only show errors
		//cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Error adding dependency '%s', '%s'", depPath, err.Error())
		}
	}

	// Add the import
	for i := 0; i < len(importsFileAst.Decls); i++ {
		d := importsFileAst.Decls[i]

		switch d.(type) {
		case *ast.FuncDecl:
		// No action
		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				// Add the new import
				newSpec := &ast.ImportSpec{Name: &ast.Ident{Name: "_"}, Path: &ast.BasicLit{Value: strconv.Quote(depPath)}}
				dd.Specs = append(dd.Specs, newSpec)
				break
			}
		}
	}

	ast.SortImports(fset, importsFileAst)

	out, err := GenerateFile(fset, importsFileAst)
	if err != nil {
		return fmt.Errorf("Error creating import file '%s', %s", importsPath, err)
	}

	err = ioutil.WriteFile(importsPath, out, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating import file '%s', %s", importsPath, err)
	}

	// Sync up
	fmt.Printf("Synching up Gopkg.yaml and imports \n")
	cmd := exec.Command("dep", "ensure")
	cmd.Dir = appDir
	newEnv := os.Environ()
	newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", rootDir))
	cmd.Env = newEnv

	// Only show errors
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error Synching up Gopkg.yaml and imports '%s', '%s'", depPath, err.Error())
	}

	fmt.Printf("'%s' installed successfully \n", depPath)

	return nil
}


// UninstallDependency deletes the given dependency
func (b *DepManager) UninstallDependency(rootDir, appDir , depPath string) error {
	exists := fgutil.ExecutableExists("dep")
	if !exists {
		return errors.New("dep not installed")
	}

	// Load imports file
	importsPath := path.Join(appDir, config.FileImportsGo)
	// Validate that it exists
	_, err := os.Stat(importsPath)

	if err != nil {
		return fmt.Errorf("Error installing dependency, import file '%s' doesn't exists", importsPath)
	}

	fset := token.NewFileSet()

	importsFileAst, err := parser.ParseFile(fset, importsPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("Error parsing import file '%s', %s", importsPath, err)
	}

	exists = false

	//Validate that the install exists in imports.go file
	for _, imp := range importsFileAst.Imports {
		if imp.Path.Value == strconv.Quote(depPath) {
			exists = true
			break
		}
	}

	if !exists{
		fmt.Printf("No import '%s' found in import file \n", depPath)
		// Just sync up and return
		// Sync up
		fmt.Printf("Synching up Gopkg.yaml and imports \n")
		cmd := exec.Command("dep", "ensure")
		cmd.Dir = appDir
		newEnv := os.Environ()
		newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", rootDir))
		cmd.Env = newEnv

		// Only show errors
		//cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("Error Synching up Gopkg.yaml and imports '%s', '%s'", depPath, err.Error())
		}

		fmt.Printf("'%s' uninstalled successfully \n", depPath)
		return nil
	}

	fmt.Printf("Deleting import from imports file \n")
	// Delete the import
	for i := 0; i < len(importsFileAst.Decls); i++ {
		d := importsFileAst.Decls[i]

		switch d.(type) {
		case *ast.FuncDecl:
		// No action
		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				var newSpecs []ast.Spec
				for _, spec := range dd.Specs {
					importSpec, ok := spec.(*ast.ImportSpec)
					if !ok{
						newSpecs = append(newSpecs, spec)
						continue
					}
					// Check Path
					if importPath := importSpec.Path; importPath.Value != strconv.Quote(depPath) {
						// Add import
						newSpecs = append(newSpecs, spec)
						continue
					}
				}
				// Update specs
				dd.Specs = newSpecs
				break
			}
		}
	}

	ast.SortImports(fset, importsFileAst)

	out, err := GenerateFile(fset, importsFileAst)
	if err != nil {
		return fmt.Errorf("Error creating import file '%s', %s", importsPath, err)
	}

	err = ioutil.WriteFile(importsPath, out, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating import file '%s', %s", importsPath, err)
	}

	// Sync up
	fmt.Printf("Synching up Gopkg.yaml and imports \n")
	cmd := exec.Command("dep", "ensure")
	cmd.Dir = appDir
	newEnv := os.Environ()
	newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", rootDir))
	cmd.Env = newEnv

	// Only show errors
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error Synching up Gopkg.yaml and imports '%s', '%s'", depPath, err.Error())
	}

	fmt.Printf("'%s' uninstalled successfully \n", depPath)
	return nil
}

// GetExistingConstraint returns the constraint definition if it already exists
func GetExistingConstraint(rootDir, appDir, depPath string) (*ConstraintDef, error) {
	// Validate that the install project does not exist in Gopkg.toml
	cmd := exec.Command("dep", "status", "-json")
	cmd.Dir = appDir
	newEnv := os.Environ()
	newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", rootDir))
	cmd.Env = newEnv

	status, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error checking project dependency status '%s'", err)
	}

	var statusMap []map[string]interface{}

	err = json.Unmarshal(status, &statusMap)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling project dependency status '%s'", err)
	}

	var existingConstraint map[string]interface{}

	for _, constraint := range statusMap {
		// Get project root
		projectRoot, ok := constraint["ProjectRoot"]
		if !ok {
			continue
		}
		pr := projectRoot.(string)
		if strings.HasPrefix(depPath, pr) {
			// Constraint already exists
			existingConstraint = constraint
			break
		}
	}

	var constraint *ConstraintDef

	if existingConstraint != nil {
		constraint = &ConstraintDef{ProjectRoot: existingConstraint["ProjectRoot"].(string), Version: existingConstraint["Version"].(string)}
	}

	return constraint, nil
}

func GenerateFile(fset *token.FileSet, file *ast.File) ([]byte, error) {
	var output []byte
	buffer := bytes.NewBuffer(output)
	if err := printer.Fprint(buffer, fset, file); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func NewDepProject() Project {
	return &DepProject{}
}

func (e *DepProject) Init(rootDir string) error {

	exists := fgutil.ExecutableExists("dep")

	if !exists {
		return errors.New("dep not installed")
	}
	e.RootDir = rootDir
	e.SourceDir = path.Join(e.RootDir, "src")
	return nil
}

// Create creates directories for the project
func (e *DepProject) Create(createBin bool, vendorDir string) error {

	if _, err := os.Stat(e.RootDir); err == nil {
		return fmt.Errorf("Cannot create project, directory '%s' already exists", e.RootDir)
	}

	os.MkdirAll(e.RootDir, os.ModePerm)
	os.MkdirAll(e.SourceDir, os.ModePerm)

	return nil
}

// Open the project directory and validate its structure
func (e *DepProject) Open() error {

	// Check root dir
	info, err := os.Stat(e.RootDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Cannot open project, directory '%s' doesn't exists", e.RootDir)
	}

	// Check source dir
	info, err = os.Stat(e.SourceDir)

	if err != nil || !info.IsDir() {
		return errors.New("Invalid project, source directory doesn't exists")
	}

	// Check file descriptor
	fd := path.Join(e.RootDir, config.FileDescriptor)
	_, err = os.Stat(fd)

	if err != nil {
		return fmt.Errorf("Invalid project, file descriptor '%s' doesn't exists", fd)
	}
	e.FileDescriptorPath = fd

	fdbytes, err := ioutil.ReadFile(fd)
	if err != nil {
		return fmt.Errorf("Invalid reading file descriptor '%s' error: %s", fd, err)
	}

	descriptor, err := ParseAppDescriptor(string(fdbytes))
	if err != nil {
		return fmt.Errorf("Invalid parsing file descriptor '%s' error: %s", fd, err)
	}

	appName := descriptor.Name

	// Validate that there is an app dir
	e.AppDir = path.Join(e.SourceDir, appName)
	info, err = os.Stat(e.AppDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Invalid project, app directory '%s' doesn't exists", e.AppDir)
	}

	return nil
}

func (e *DepProject) GetBinDir() string {
	return e.BinDir
}

func (e *DepProject) GetRootDir() string {
	return e.RootDir
}

func (e *DepProject) GetSourceDir() string {
	return e.SourceDir
}

func (e *DepProject) GetVendorDir() string {
	return e.VendorDir
}

func (e *DepProject) GetVendorSrcDir() string {
	return e.VendorSrcDir
}

// GetAppDir returns the directory of the app
func (e *DepProject) GetAppDir() string {
	return e.AppDir
}

func (e *DepProject) InstallDependency(depPath string, version string) error {
	var cmd *exec.Cmd

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//check if dependency is installed
	if _, err := os.Stat(path.Join(e.VendorSrcDir, depPath)); err == nil {
		//todo ignore installed dependencies for now
		//exists, return
		return nil
	}

	if version == "" {
		//if strings.HasPrefix(depPath,"github.com/TIBCOSoftware/flogo-") {
		//	cmd = exec.Command("gb", "vendor", "fetch", "-branch", "entrypoint", depPath)
		//} else {
		cmd = exec.Command("gb", "vendor", "fetch", depPath)
		//}
	} else {
		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, depPath)
	}

	os.Chdir(e.RootDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *DepProject) UninstallDependency(depPath string) error {

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//check if dependency is installed
	if _, err := os.Stat(path.Join(e.VendorSrcDir, depPath)); err != nil {
		//todo ignore dependencies that are not installed for now
		//exists, return
		return nil
	}

	os.Chdir(e.RootDir)

	cmd := exec.Command("gb", "vendor", "delete", depPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *DepProject) Build() error {
	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	os.Chdir(e.RootDir)
	fmt.Println(e.RootDir)

	return cmd.Run()
}

func IsDepProject(projectPath string) bool {

	sourceDir := path.Join(projectPath, "src")
	vendorDir := path.Join(projectPath, "vendor", "src")

	info, err := os.Stat(sourceDir)

	if err != nil || !info.IsDir() {
		return false
	}

	info, err = os.Stat(vendorDir)

	if err != nil || !info.IsDir() {
		return false
	}

	return true
}

//Env checker?
//IsProject(path.Join string) bool

// Gb structure that contains gb project paths
type Dep struct {
	BinPath        string
	SourcePath     string
	VendorPath     string
	CodeSourcePath string
}

// NewGb creates a new Gb struct
func NewDep(codePath string) *Gb {

	env := &Gb{}
	env.BinPath = "bin"
	env.SourcePath = "src"
	env.VendorPath = path.Join("vendor", "src")
	env.CodeSourcePath = path.Join("src", codePath)

	return env
}

// Init creates directories for the gb project
func (e *Dep) Init(createBin bool) {
	os.MkdirAll(e.SourcePath, 0777)
	os.MkdirAll(e.VendorPath, 0777)
	os.MkdirAll(e.CodeSourcePath, 0777)

	if createBin {
		os.MkdirAll(e.BinPath, 0777)
	}
}

// Installed indicates if gb is installed
func (e *Dep) Installed() bool {
	return fgutil.ExecutableExists("gb")
}

// NewBinFilepath.Join creates a new file path.Join in the bin directory
func (e *Dep) NewBinFilePath(fileName string) string {
	return path.Join(e.BinPath, fileName)
}

// VendorFetch performs a 'gb vendor fetch'
func (e *Dep) VendorFetch(depPath string, version string) error {

	var cmd *exec.Cmd

	if version == "" {
		cmd = exec.Command("gb", "vendor", "fetch", depPath)
	} else {

		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, depPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// VendorDeleteSilent performs a 'gb vendor delete' silently
func (e *Dep) VendorDeleteSilent(depPath string) error {
	cmd := exec.Command("gb", "vendor", "delete", depPath)

	return cmd.Run()
}

// VendorDelete performs a 'gb vendor delete'
func (e *Dep) VendorDelete(depPath string) error {
	cmd := exec.Command("gb", "vendor", "delete", depPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Build performs a 'gb build'
func (e *Dep) Build() error {
	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
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

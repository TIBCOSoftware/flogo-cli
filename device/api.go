package device

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"io"
	"text/template"
	"fmt"
	"strconv"
	"errors"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

// BuildPreProcessor interface for build pre-processors
type BuildPreProcessor interface {
	PrepareForBuild(env Project) error
}

// CreateDevice creates an device project from the specified json device descriptor
func CreateDevice(project Project, deviceJson string, deviceDir string, deviceName string) error {

	descriptor, err := ParseDeviceDescriptor(deviceJson)
	if err != nil {
		return err
	}

	if deviceName != "" {
		// override the application name

		altJson := strings.Replace(deviceJson, `"`+descriptor.Name+`"`, `"`+deviceName+`"`, 1)
		altDescriptor, err := ParseDeviceDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Name == deviceName {
			deviceJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(deviceJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = deviceName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			deviceJson = string(updApp)
		}

		descriptor.Name = deviceName
	}

	project.Init(deviceDir)
	project.Create()

	profile, err := GetDeviceProfile(project, descriptor.Device.Profile)
	if err != nil {
		return err
	}

	err = project.Setup(profile.Board)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromString(path.Join(deviceDir, "device.json"), deviceJson)
	if err != nil {
		return err
	}

	_, err = generateSourceFiles(project, descriptor, profile)
	if err != nil {
		return err
	}

	return nil
}

type PrepareOptions struct {
	PreProcessor BuildPreProcessor
}

// PrepareDevice do all pre-build setup and pre-processing
func PrepareDevice(project Project, options *PrepareOptions) (err error) {

	if options == nil {
		options = &PrepareOptions{}
	}

	if options.PreProcessor != nil {
		err = options.PreProcessor.PrepareForBuild(project)
		if err != nil {
			return err
		}
	}

	//load descriptor
	appJson, err := fgutil.LoadLocalFile(path.Join(project.GetRootDir(), "device.json"))
	if err != nil {
		return err
	}

	descriptor, err := ParseDeviceDescriptor(appJson)
	if err != nil {
		return err
	}

	profile, err := GetDeviceProfile(project, descriptor.Device.Profile)
	if err != nil {
		return err
	}

	libs, err := generateSourceFiles(project, descriptor, profile)
	if err != nil {
		return err
	}

	err = InstallLibs(project, libs)
	if err != nil {
		return err
	}

	return nil
}

func InstallLibs(project Project, libs []*Lib) error {
	for _, lib := range libs {

		//assume platformio for now
		if lib.LibType != "platformio" {
			return errors.New("Unsupported lib type: " + lib.LibType)
		}

		libId, err := strconv.Atoi(lib.Ref)
		if err != nil {
			return errors.New("Error parsing lib id: " + lib.Ref)
		}
		err = project.InstallLib("", libId)
		if err != nil {
			return err
		}
	}

	return nil

}

func GetDeviceProfile(proj Project, ref string) (*DeviceProfile, error) {

	proj.InstallContribution(ref, "")

	descFile := path.Join(proj.GetContributionDir(), ref, "profile.json")
	profJson, err := fgutil.LoadLocalFile(descFile)

	if err != nil {
		return nil, err
	}

	profile, err := ParseDeviceProfile(profJson)

	return profile, err
}

func GetDevicePlatform(proj Project, ref string) (*DevicePlatform, error) {

	proj.InstallContribution(ref, "")

	platformFile := path.Join(proj.GetContributionDir(), ref, "platform.json")
	platformJson, err := fgutil.LoadLocalFile(platformFile)

	if err != nil {
		return nil, err
	}

	platform, err := ParseDevicePlatform(platformJson)

	return platform, err
}

type BuildOptions struct {
	*PrepareOptions

	SkipPrepare bool
}

// BuildDevice build the flogo application
func BuildDevice(env Project, options *BuildOptions) (err error) {

	if options == nil {
		options = &BuildOptions{}
	}

	if !options.SkipPrepare {
		err = PrepareDevice(env, options.PrepareOptions)

		if err != nil {
			return err
		}
	}

	err = env.Build()
	if err != nil {
		return err
	}

	return
}

// InstallDependency install a dependency
func InstallContribution(project Project, path string, version string) error {

	return project.InstallContribution(path, version)
}

// UploadDevice upload the device application
func UploadDevice(env Project) error {

	err := env.Upload()

	return err
}

func generateSourceFiles(proj Project, descriptor *FlogoDeviceDescriptor, profile *DeviceProfile) ([]*Lib, error) {

	libMap := make(map[string]*Lib)

	generatePlatformCode(proj, descriptor, profile)

	srcDir := proj.GetSourceDir()

	//generate triggers
	for _, trgCfg := range descriptor.Triggers {

		trgContrib, err := LoadTriggerContrib(proj, trgCfg.Ref)

		if err != nil {
			panic("Trigger '" + trgCfg.Ref + "' not found")
		}

		if len(trgContrib.Libs()) > 0 {
			for _, lib := range trgContrib.Libs() {
				libMap[lib.LibType+lib.Ref] = lib
			}
		}

		f, _ := os.Create(path.Join(srcDir, trgCfg.Id+".ino"))
		tpl := trgContrib.Template
		err = RenderTemplate(f, tpl, trgCfg)
		f.Close()
		if err != nil {
			return nil, err
		}
	}

	// generate actions
	for _, actCfg := range descriptor.Actions {

		// Action should define how to do this

		if strings.HasSuffix(actCfg.Ref, "flow") {
			var flowCfg *FlowActionConfig
			err := json.Unmarshal(actCfg.Data, &flowCfg)
			if err != nil {
				errorMsg := fmt.Sprintf("Error while loading flow '%s' error '%s'", actCfg.Id, err.Error())
				panic(errorMsg)
			}

			flowTree := toFlowTree(actCfg.Id, flowCfg.Flow)

			//generate activities
			for _, task := range flowTree.AllTasks {

				f, _ := os.Create(path.Join(srcDir, actCfg.Id+"_"+strconv.Itoa(task.Id)+".ino"))

				actContrib, err := LoadActivityContrib(proj, task.ActivityRef)

				if err != nil {
					panic("Activity '" + task.ActivityRef + "' not found")
				}

				if len(actContrib.Libs()) > 0 {
					for _, lib := range actContrib.Libs() {
						libMap[lib.LibType+lib.Ref] = lib
					}
				}

				tpl := actContrib.Template

				data := struct {
					Id       string
					Activity *Task
				}{
					actCfg.Id,
					task,
				}

				err = RenderTemplate(f, tpl, data)
				f.Close()
				if err != nil {
					return nil, err
				}
			}

			f, _ := os.Create(path.Join(srcDir, actCfg.Id+".ino"))

			t := template.New("action")
			t.Funcs(DeviceFuncMap)
			_, err = t.Parse(actionContribs[actCfg.Ref].Template)
			if err != nil {
				return nil, err
			}

			tmpl := t.New("actioneval")
			_, err = tmpl.Parse(tplActionDeviceFlowEval)
			if err != nil {
				return nil, err
			}
			err = t.Execute(f, flowTree)
			f.Close()
			if err != nil {
				return nil, err
			}

		} else {
			var aaCfg *ActivityActionConfig
			err := json.Unmarshal(actCfg.Data, aaCfg)
			if err != nil {
				errorMsg := fmt.Sprintf("Error while loading activity action '%s' error '%s'", actCfg.Id, err.Error())
				panic(errorMsg)
			}

			if _, ok := activityContribs[aaCfg.Activity.Ref]; !ok {
				panic("Activity '" + aaCfg.Activity.Ref + "' not found")
			}

			if len(activityContribs[aaCfg.Activity.Ref].Libs()) > 0 {
				for _, lib := range activityContribs[aaCfg.Activity.Ref].Libs() {
					libMap[lib.LibType+lib.Ref] = lib
				}
			}

			f, _ := os.Create(path.Join(srcDir, "ac_ "+actCfg.Id+".ino"))
			tpl := activityContribs[aaCfg.Activity.Ref].Template
			err = RenderTemplate(f, tpl, aaCfg.Activity)
			f.Close()
			if err != nil {
				return nil, err
			}

			f, _ = os.Create(path.Join(srcDir, actCfg.Id+".ino"))

			t := template.New("action")
			t.Funcs(DeviceFuncMap)
			_, err = t.Parse(actionContribs[actCfg.Ref].Template)
			if err != nil {
				return nil, err
			}
			err = t.Execute(f, actCfg)
			f.Close()
			if err != nil {
				return nil, err
			}
		}
	}


	libs := make([]*Lib, 0, len(libMap))

	for  _, value := range libMap {
		libs = append(libs, value)
	}

	return libs, nil
}

func generatePlatformCode(proj Project, descriptor *FlogoDeviceDescriptor, profile *DeviceProfile) error {

	platform, err := GetDevicePlatform(proj, profile.Platform)
	if err != nil {
		return err
	}

	var actionIds []string

	for _, value := range descriptor.Actions {
		actionIds = append(actionIds, value.Id)
	}

	var triggerIds []string
	mqttTriggers := make(map[string]string)

	for _, value := range descriptor.Triggers {

		//todo fix mqtt determination
		if !strings.Contains(value.Ref, "mqtt") {
			triggerIds = append(triggerIds, value.Id)
		} else {
			mqttTriggers[value.Id] = value.Settings["topic"]
		}
	}

	data := struct {
		MqttEnabled  bool
		Actions      []string
		Triggers     []string
		MqttTriggers map[string]string
	}{
		descriptor.Device.MqttEnabled,
		actionIds,
		triggerIds,
		mqttTriggers,
	}

	platformDir := path.Join(proj.GetContributionDir(), profile.Platform)

	tpl,err := fgutil.LoadLocalFile(path.Join(platformDir, platform.MainTemplate))
	if err != nil {
		return err
	}

	//todo fix main file name generation
	f, _ := os.Create(path.Join(proj.GetSourceDir(), "main.ino"))
	err = RenderTemplate(f, tpl, data)
	f.Close()
	if err != nil {
		return err
	}

	if descriptor.Device.MqttEnabled {
		err = generateWifiCode(proj, descriptor, platform, profile)
		if err != nil {
			return err
		}

		err = generateMqttCode(proj, descriptor, platform, profile)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateWifiCode(project Project, descriptor *FlogoDeviceDescriptor, platform *DevicePlatform, profile *DeviceProfile) error {

	settings := &SettingsConfig{DeviceName: descriptor.Name, Settings: descriptor.Device.Settings}

	for _, value := range platform.WifiDetails {
		if value.Name == profile.PlatformWifi {

			tpl,err := fgutil.LoadLocalFile(path.Join(project.GetContributionDir(), profile.Platform, value.Template))
			if err != nil {
				return err
			}

			//todo fix mqtt file name generation
			f, _ := os.Create(path.Join(project.GetSourceDir(), "wifi.ino"))
			err = RenderTemplate(f, tpl, settings)
			f.Close()
			if err != nil {
				return err
			}

			InstallLibs(project, value.Libs)

			return nil
		}

	}

	return nil
}

func generateMqttCode(project Project, descriptor *FlogoDeviceDescriptor, platform *DevicePlatform, profile *DeviceProfile) error {

	settings := &SettingsConfig{DeviceName: descriptor.Name, Settings: descriptor.Device.Settings}

	tpl,err := fgutil.LoadLocalFile(path.Join(project.GetContributionDir(), profile.Platform, platform.MqttDetails.Template))
	if err != nil {
		return err
	}

	//todo fix mqtt file name generation
	f, _ := os.Create(path.Join(project.GetSourceDir(), "mqtt.ino"))
	err = RenderTemplate(f, tpl, settings)
	f.Close()
	if err != nil {
		return err
	}

	InstallLibs(project, platform.MqttDetails.Libs)

	return nil
}

//RenderFileTemplate renders the specified template
func RenderTemplate(w io.Writer, tpl string, data interface{}) error {

	t := template.New("source")
	t.Funcs(DeviceFuncMap)

	t.Parse(tpl)
	err := t.Execute(w, data)

	return err
}

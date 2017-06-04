package device

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"io"
	"text/template"
	"fmt"
	"strconv"
	"errors"
)

// BuildPreProcessor interface for build pre-processors
type BuildPreProcessor interface {
	PrepareForBuild(env Project) error
}

// CreateDevice creates an device project from the specified json device descriptor
func CreateDevice(env Project, deviceJson string, deviceDir string, deviceName string) error {

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

	env.Init(deviceDir)

	details := GetDevice(descriptor.DeviceType)

	err = env.Create(details.Board)
	if err != nil {
		return err
	}
	err = fgutil.CreateFileFromString(path.Join(deviceDir, "device.json"), deviceJson)
	if err != nil {
		return err
	}

	_, err = generateSourceFiles(env, descriptor, false)
	if err != nil {
		return err
	}

	return nil
}

type PrepareOptions struct {
	PreProcessor    BuildPreProcessor
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
	appJson, err := fgutil.LoadLocalFile(path.Join(project.GetRootDir(),"device.json"))
	if err != nil {
		return err
	}

	descriptor, err := ParseDeviceDescriptor(appJson)
	if err != nil {
		return err
	}

	libs, err := generateSourceFiles(project, descriptor, true)
	if err != nil {
		return err
	}

	details := GetDevice(descriptor.DeviceType)

	for name, id := range details.Libs {

		err := project.InstallLib(name, id)
		if err != nil {
			return err
		}
	}

	for _, lib := range libs {

		//assume platformio for now
		libId,err :=strconv.Atoi(lib.Ref)
		if err != nil {
			return errors.New("Unsupported lib type: " + lib.LibType)
		}
		err = project.InstallLib("",libId)
		if err != nil {
			return err
		}
	}

	return nil
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

// UploadDevice upload the device application
func UploadDevice(env Project) error {

	err := env.Upload()

	return err
}

func generateSourceFiles(env Project, descriptor *FlogoDeviceDescriptor, skipMain bool) (map[string]*Lib,error) {

	libs := make(map[string]*Lib)

	details := GetDevice(descriptor.DeviceType)

	srcDir := env.GetSourceDir()
	//if !skipMain {
		generateMainCode(srcDir, descriptor, details)
	//}

	if descriptor.MqttEnabled {

		generateMqttCode(srcDir, descriptor, details)
	}

	//generate triggers
	for _, trgCfg := range descriptor.Triggers {

		if _,ok:= triggerContribs[trgCfg.Ref]; !ok {
			panic("Trigger '" + trgCfg.Ref + "' not found")
		}

		if len(triggerContribs[trgCfg.Ref].libs) > 0 {
			for _, lib := range triggerContribs[trgCfg.Ref].libs {
				libs[lib.LibType + lib.Ref] = lib
			}
		}

		f, _ := os.Create(path.Join(srcDir, trgCfg.Id + ".ino"))
		tpl := triggerContribs[trgCfg.Ref].Template
		err := RenderTemplate(f, tpl, trgCfg)
		f.Close()
		if err != nil {
			return nil,err
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

				f, _ := os.Create(path.Join(srcDir, actCfg.Id + "_" + strconv.Itoa(task.Id) + ".ino"))

				if _,ok:= activityContribs[task.ActivityRef]; !ok {
					panic("Activity '" + task.ActivityRef + "' not found")
				}

				if len(activityContribs[task.ActivityRef].libs) > 0 {
					for _, lib := range activityContribs[task.ActivityRef].libs {
						libs[lib.LibType + lib.Ref] = lib
					}
				}

				tpl := activityContribs[task.ActivityRef].Template

				data := struct {
					Id       string
					Activity *Task
				}{
					actCfg.Id,
					task,
				}

				err := RenderTemplate(f, tpl, data)
				f.Close()
				if err != nil {
					return nil,err
				}
			}

			f, _ := os.Create(path.Join(srcDir, actCfg.Id + ".ino"))

			t := template.New("action")
			t.Funcs(DeviceFuncMap)
			_, err = t.Parse(actionContribs[actCfg.Ref].Template)
			if err != nil {
				return nil,err
			}

			tmpl := t.New("actioneval")
			_, err = tmpl.Parse(tplActionDeviceFlowEval)
			if err != nil {
				return nil,err
			}
			err = t.Execute(f, flowTree)
			f.Close()
			if err != nil {
				return nil,err
			}

		} else {
			var aaCfg *ActivityActionConfig
			err := json.Unmarshal(actCfg.Data, aaCfg)
			if err != nil {
				errorMsg := fmt.Sprintf("Error while loading activity action '%s' error '%s'", actCfg.Id, err.Error())
				panic(errorMsg)
			}

			if _,ok:= activityContribs[aaCfg.Activity.Ref]; !ok {
				panic("Activity '" + aaCfg.Activity.Ref + "' not found")
			}

			if len(activityContribs[aaCfg.Activity.Ref].libs) > 0 {
				for _, lib := range activityContribs[aaCfg.Activity.Ref].libs {
					libs[lib.LibType + lib.Ref] = lib
				}
			}

			f, _ := os.Create(path.Join(srcDir, "ac_ " + actCfg.Id + ".ino"))
			tpl := activityContribs[aaCfg.Activity.Ref].Template
			err = RenderTemplate(f, tpl, aaCfg.Activity)
			f.Close()
			if err != nil {
				return nil, err
			}

			f, _ = os.Create(path.Join(srcDir, actCfg.Id + ".ino"))

			t := template.New("action")
			t.Funcs(DeviceFuncMap)
			_, err = t.Parse(actionContribs[actCfg.Ref].Template)
			if err != nil {
				return nil,err
			}
			err = t.Execute(f, actCfg)
			f.Close()
			if err != nil {
				return nil,err
			}
		}
	}

	return libs,nil
}

func generateMainCode(srcDir string, descriptor *FlogoDeviceDescriptor, details *DeviceDetails) error {
	var actionIds []string

	for _, value := range descriptor.Actions {
		actionIds = append(actionIds, value.Id)
	}

	var triggerIds []string
	mqttTriggers := make(map[string]string)

	for _, value := range descriptor.Triggers {

		if !strings.Contains(value.Ref, "mqtt") {
			triggerIds = append(triggerIds, value.Id)
		} else {
			mqttTriggers[value.Id] = value.Settings["topic"]
			//mqttTriggerIds = append(mqttTriggerIds, value.Id)
		}
	}

	data := struct {
		MqttEnabled bool
		Actions []string
		Triggers []string
		MqttTriggers map[string]string
	}{
		descriptor.MqttEnabled,
		actionIds,
		triggerIds,
		mqttTriggers,
	}

	tpl := details.MainFile
	//todo fix name resolution
	f, _ := os.Create(path.Join(srcDir, "main.ino"))
	err := RenderTemplate(f, tpl, data)
	f.Close()
	if err != nil {
		return err
	}

	return nil
}

func generateMqttCode(srcDir string, descriptor *FlogoDeviceDescriptor, details *DeviceDetails) error {

	settings := &SettingsConfig{DeviceName: descriptor.Name, Settings:descriptor.Settings}

	for name, tpl := range details.MqttFiles {

		f, _ := os.Create(path.Join(srcDir, name))
		err := RenderTemplate(f, tpl, settings)
		f.Close()
		if err != nil {
			return err
		}
	}

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

package gen

import (
	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileTriggerDescriptor string = "trigger.json"
	fileTriggerGo         string = "trigger.go"
	fileTriggerGoTest     string = "trigger_test.go"
)

type TriggerGenerator struct {

}

func (g *TriggerGenerator) Description() string {
	return "generates a trigger project"
}

func (g *TriggerGenerator) Generate(basePath string, data interface{}) error {

	err := fgutil.CreateFileFromTemplate(basePath, fileTriggerDescriptor, tplTriggerDescriptor, data)
	if err != nil {
		return err
	}
	err = fgutil.CreateFileFromTemplate(basePath, fileTriggerGo, tplTriggerGo, data)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromTemplate(basePath, fileTriggerGoTest, tplTriggerGoTestGo, data)
	if err != nil {
		return err
	}

	return nil
}

var tplTriggerDescriptor = `{
  "name": "{{.Name}}",
  "version": "0.0.1",
  "type": "flogo:trigger",
  "description": "trigger description",
  "author": "Your Name <you.name@example.org>",
  "settings":[
    {
      "name": "setting",
      "type": "string",
      "value": "default"
    }
  ],
  "outputs": [
    {
      "name": "output",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "handler_setting",
        "type": "string"
      }
    ]
  }
}`

var tplTriggerGo = `package {{.Name}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

// MyTriggerFactory My Trigger factory
type MyTriggerFactory struct{
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MyTriggerFactory{metadata:md}
}

//New Creates a new trigger instance for a given id
func (t *MyTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &MyTrigger{metadata: t.metadata, config:config}
}

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {
	// start the trigger
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	return nil
}
`

var tplTriggerGoTestGo = `package {{.Name}}

import (
	"context"
	"io/ioutil"
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

func getJsonMetadata() string{
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil{
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	return 0, nil, nil
}

const testConfig string = ` + "`" + `{
  "id": "mytrigger",
  "settings": {
    "setting": "somevalue"
  },
  "handlers": [
    {
      "actionId": "test_action",
      "settings": {
        "handler_setting": "somevalue"
      }
    }
  ]
}` + "`" +`

func TestInit(t *testing.T) {

	// New factory
	md := trigger.NewMetadata(getJsonMetadata())
	f := NewFactory(md)

	// New Trigger
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)
}
`
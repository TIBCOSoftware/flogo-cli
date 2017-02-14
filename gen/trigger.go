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
      "name": "input",
      "type": "string",
      "value": "default"
    }
  ],
  "outputs": [
    {
      "name": "output",
      "type": "string"
    }
  ]
}`

var tplTriggerGo = `package {{.Name}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&MyTrigger{metadata: md})
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(config *trigger.Config, runner action.Runner) {
	t.config = config
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
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	return 0, nil, nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("{{.Name}}")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}
`
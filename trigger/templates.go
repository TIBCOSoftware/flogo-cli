package trigger

var tplTriggerJSON = `{
  "name": "{{.Name}}",
  "version": "0.0.1",
  "description": "trigger description",
  "config":[
    {
      "name": "input",
      "type": "string",
      "value": "default"
    }
  ]
}`

var tplTriggerGoFile = `package {{.Name}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
)

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata       *trigger.Metadata
	flowStarter flowinst.Starter
	config         *trigger.Config
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&MyTrigger{metadata: md})
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(flowStarter flowinst.Starter, config *trigger.Config) {
	t.flowStarter = flowStarter
	t.config = config
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() {
	// start the trigger
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() {
	// stop the trigger
}
`

var tplTriggerTestGoFile = `package {{.Name}}

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
)

func TestRegistered(t *testing.T) {
	act := trigger.Get("{{.Name}}")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}
`

var tplMetadataGoFile = `package {{.Name}}

var jsonMetadata = ` + "`" + tplTriggerJSON + "`"

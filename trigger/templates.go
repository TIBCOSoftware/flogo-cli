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
	"github.com/TIBCOSoftware/flogo-lib/engine/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/starter"
)

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata       *trigger.Metadata
	processStarter starter.ProcessStarter
	config         map[string]string
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&MyTrigger{metadata: md})
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(processStarter starter.ProcessStarter, config map[string]string) {
	t.processStarter = processStarter
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
	"github.com/TIBCOSoftware/flogo-lib/engine/ext/trigger"
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

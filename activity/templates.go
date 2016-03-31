package activity

var tplActivityJSON = `{
  "name": "{{.Name}}",
  "version": "0.0.1",
  "description": "activity description",
  "inputs":[
    {
      "name": "input",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "output",
      "type": "string"
    }
  ]
}`

var tplActivityGoFile = `package {{.Name}}

import (
	"github.com/TIBCOSoftware/flogo/golib/core/ext/activity"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&MyActivity{metadata: md})
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) bool {

	// do eval

	return true //done
}
`
var tplActivityTestGoFile = `package {{.Name}}

import (
	"testing"
	"github.com/TIBCOSoftware/flogo/golib/core/ext/activity"
	"github.com/TIBCOSoftware/flogo/golib/test"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("{{.Name}}")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)
	act := &MyActivity{metadata: md}

	tc := test.NewTestActivityContext()
	//setup attrs

	act.Eval(tc)

	//check result attr
}
`

var tplMetadataGoFile = `package {{.Name}}

var jsonMetadata = ` + "`" + tplActivityJSON + "`"

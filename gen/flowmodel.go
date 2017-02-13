package gen

import (
	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileFlowModelDescriptor string = "model.json"
	fileFlowModelGo         string = "model.go"
	fileFlowModelGoTest     string = "model_test.go"
)

type FlowModelGenerator struct {

}

func (g *FlowModelGenerator) Description() string {
	return "generates a flow-model project"
}

func (g *FlowModelGenerator) Generate(basePath string, data interface{}) error {

	err := fgutil.CreateFileFromTemplate(basePath, fileFlowModelDescriptor, tplFlowModelDescriptor, data)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromTemplate(basePath, fileFlowModelGo, tplFlowModelGo, data)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromTemplate(basePath, fileFlowModelGoTest, tplFlowModelGoTest, data)
	if err != nil {
		return err
	}

	return nil
}

var tplFlowModelDescriptor = `{
  "name": "{{.Name}}",
  "version": "0.0.1",
  "type": "flogo:flow-model",
  "description": "model description"
}`

var tplFlowModelGo = `package {{.Name}}

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/flow/model"
)

func init() {
	m := model.New("{{.Name}}")
	m.RegisterFlowBehavior(1, &MyFlowBehavior{})
	m.RegisterTaskBehavior(1, &MyTaskBehavior{})
	m.RegisterLinkBehavior(1, &MyLinkBehavior{})

	model.Register(m)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

// MyFlowBehavior implements model.FlowBehavior
type MyFlowBehavior struct {
}

// Start implements model.FlowBehavior.Start
func (pb *MyFlowBehavior) Start(context model.FlowContext, data interface{}) (start bool, evalCode int) {
	// just schedule the root task
	return true, 0
}

// Resume implements model.FlowBehavior.Resume
func (pb *MyFlowBehavior) Resume(context model.FlowContext, data interface{}) bool {
	return true
}

// TasksDone implements model.FlowBehavior.TasksDone
func (pb *MyFlowBehavior) TasksDone(context model.FlowContext, doneCode int) {
	// all tasks are done
}

// Done implements model.FlowBehavior.Done
func (pb *MyFlowBehavior) Done(context model.FlowContext) {
	fmt.Printf("Flow Done\n")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

// MyTaskBehavior implements model.TaskBehavior
type MyTaskBehavior struct {
}

// Enter implements model.TaskBehavior.Enter
func (tb *MyTaskBehavior) Enter(context model.TaskContext, enterCode int) (eval bool, evalCode int) {

	context.SetState(STATE_ENTERED)
	linkContexts := context.FromLinks()
	ready := true

	if len(linkContexts) == 0 {
		// has no predecessor links, so task is ready
		ready = true
	} else {
		// check if all pedecessor links are done
		for _, linkContext := range linkContexts {

			if linkContext.State() != STATE_LINK_TRUE {
				ready = false
				break
			}
		}
	}

	if ready {
		context.SetState(STATE_READY)
	}

	return ready, 0
}

// Eval implements model.TaskBehavior.Eval
func (tb *MyTaskBehavior) Eval(context model.TaskContext, evalCode int) (done bool, doneCode int) {

	task := context.Task()

	if len(task.ChildTasks()) > 0 {
		//has children, so set to waiting
		context.SetState(STATE_WAITING)

		//for now enter all children (bpel style) - todo: change to enter leading chlidren
		context.EnterChildren(nil)

		return false, 0
	}

	activity, activityContext := context.Activity()

	if activity != nil {
		done := activity.Eval(activityContext)
		return done, 0
	}

	// doesn't have an activity so treat as no-op
	return true, 0
}

// PostEval implements model.TaskBehavior.PostEval
func (tb *MyTaskBehavior) PostEval(context model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int) {
	// ignore, just mark done
	return true, 0
}

// Done implements model.TaskBehavior.Done
func (tb *MyTaskBehavior) Done(context model.TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*model.TaskEntry) {

	task := context.Task()

	context.SetState(STATE_DONE)
	//context.SetTaskDone() for task garbage collection

	links := task.ToLinks()
	numLinks := len(links)

	// flow outgoing links
	if numLinks > 0 {

		taskEntries := make([]*model.TaskEntry, 0, numLinks)

		for _, link := range links {

			linkContext := context.EvalLink(link, 0)
			if linkContext.State() == STATE_LINK_TRUE {

				taskEntry := &model.TaskEntry{Task: link.ToTask(), EnterCode: 0}
				taskEntries = append(taskEntries, taskEntry)
			}
		}

		//continue on to successor tasks
		return false, 0, taskEntries
	}

	// there are no outgoing links, so just notify parent that we are done
	return true, 0, nil
}

// ChildDone implements model.TaskBehavior.ChildDone
func (tb *MyTaskBehavior) ChildDone(context model.TaskContext, childTask *flow.Task, childDoneCode int) (done bool, doneCode int) {

	// our children are done, so just transition ourselves to done
	return true, 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

// MyLinkBehavior implements model.LinkBehavior
type MyLinkBehavior struct {
}

// Eval implements model.LinkBehavior.Eval
func (lb *MyLinkBehavior) Eval(context model.LinkContext, evalCode int) {

	context.SetState(STATE_LINK_TRUE)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// State
const (
	STATE_NOT_STARTED int = 0

	STATE_LINK_FALSE int = 1
	STATE_LINK_TRUE  int = 2

	STATE_ENTERED int = 10
	STATE_READY   int = 20
	STATE_WAITING int = 30
	STATE_DONE    int = 40
)
`

var tplFlowModelGoTest = `package {{.Name}}

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/flow/model"
)

func TestRegistered(t *testing.T) {
	act := model.Get("{{.Name}}")

	if act == nil {
		t.Error("Model Not Registered")
		t.Fail()
		return
	}
}
`

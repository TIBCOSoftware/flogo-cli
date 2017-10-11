package app

import (
	"testing"
	"flag"
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"os"
)

type TestEnv struct {
	currentDir string
}

func (t *TestEnv) getTestwd() (dir string, err error){
	return t.currentDir, nil
}

func (t *TestEnv) cleanup(){
	os.RemoveAll(t.currentDir)
}


// TestCmdCreate_Exec test the default cmd create, create new app
func TestCmdCreate_Exec (t *testing.T) {
	// TODO remote this after merging
	err := os.Setenv("FLOGO_BUILD_EXPERIMENTAL", "true")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("FLOGO_BUILD_EXPERIMENTAL")
	testEnv := &TestEnv{currentDir:os.TempDir()}
	defer testEnv.cleanup()
	t.Logf("Current dir '%s'", testEnv.currentDir)
	cmd := &cmdCreate{option: optCreate, currentDir: testEnv.getTestwd}
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	args := []string{"testApp"}

	err = cli.ExecCommand(fs, cmd, args)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the structure
	currentDir := testEnv.currentDir
	_, err = os.Stat(currentDir)
	if err != nil {
		t.Fatal(err)
	}


}

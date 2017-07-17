package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/device"
)

var (
	fs = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

func init() {
	fs.Usage = device.Usage
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "FATAL: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) < 2 || args[1] == "-h" {
		device.Usage()
	}
	name := args[1]

	var remainingArgs []string

	cmd, exists := device.CommandRegistry.Command(name)

	if !exists {
		fmt.Fprintf(os.Stderr, "FATAL: unknown command %q\n\n", name)
		device.Usage()
	}
	remainingArgs = args[2:]

	if err := cli.ExecCommand(fs, cmd, remainingArgs); err != nil {
		fatalf("command %q failed: %v", name, err)
	}

	os.Exit(0)
}

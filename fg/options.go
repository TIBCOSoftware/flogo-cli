package flogo

// OptionInfo is the option information for a command
type OptionInfo struct {
	// Denotes if tool or command
	IsTool bool

	// Name of the tool/command
	Name string

	// UsageLine demonstrates how to invoke the tool/command
	UsageLine string

	// Short description of tool/command
	Short string

	// Description of this tool/command
	Long string
}

// HasOptionInfo is an interface for an object that
// has Option Information
type HasOptionInfo interface {
	OptionInfo() *OptionInfo
}

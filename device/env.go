package device

type Project interface {

	// GetRootDir get the root directory of the project
	GetRootDir() string

	// GetSourceDir get the source directory of the project
	GetSourceDir() string

	// GetLibDir get the lib directory of the project
	GetLibDir() string

	// GetContributionDir get the contribution directory of the project
	GetContributionDir() string

	// Init initializes the project settings an validates it requirements
	Init(path string) error

	// Create the project directory and its structure
	Create() error

	// Setup the project directory
	Setup(board string) error

	// Open the project directory and validate its structure
	Open() error

	InstallLib(name string, id int) error

	// InstallContribution install a contribution
	InstallContribution(path string, version string) error

	// UninstallContribution uninstall a contribution
	UninstallContribution(path string) error

	Build() error

	Upload() error

	Clean() error
}

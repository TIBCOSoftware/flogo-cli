package device


type Project interface {

	// GetRootDir get the root directory of the project
	GetRootDir() string

	// GetSourceDir get the source directory of the project
	GetSourceDir() string

	// GetLibDir get the lib directory of the project
	GetLibDir() string

	// Init initializes the project settings an validates it requirements
	Init(path string) error

	// Create the project directory and its structure, optional existing vendor dir to copy
	Create(board string) error

	// Open the project directory and validate its structure
	Open() error

	InstallLib(name string, id int) error

	Build() error

	Upload() error

	Clean() error
}

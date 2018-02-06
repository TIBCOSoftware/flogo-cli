package gen

type CodeGenerator interface {
	Description() string

	Generate(basePath string, data interface{}) error
}

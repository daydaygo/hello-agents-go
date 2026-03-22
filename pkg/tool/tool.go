package tool

type Tool interface {
	Name() string
	Description() string
	Execute(input string) (string, error)
}

package tools

type Tool struct {
	Name        string
	Description string

	Check func() bool
}

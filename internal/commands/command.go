package commands

type Command interface {
	Description() string
	GetName() string
}

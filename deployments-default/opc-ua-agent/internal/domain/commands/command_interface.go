package domain_commands

type ICommand interface {
	Execute() (interface{}, error)
}

package domain_commands

type ICommand interface {
	Execute() (interface{}, error)

	GetCorrelationId() string
}

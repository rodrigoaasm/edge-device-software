package domain_interfaces

type ITransactionManager interface {
	Commit() error
	Rollback() error
}

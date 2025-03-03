package protocol

type NodeInterface interface {
	Execute() error
	Rollback() error
}

package command

type Interface interface {
	GetHeader() Header
	Invoke() any
}
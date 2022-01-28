package common

type Adapter interface {
	CheckParam() bool

	Execute()

	GetMode() AdapterMode
}

type AdapterMode int

const (
	TelnetAdapterMode AdapterMode = iota
	GeneratorAdapterMode
)

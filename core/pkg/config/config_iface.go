package config

type ILoader interface {
	Load(name string) ([]byte, error)
}

type ITemplate interface {
	Load(loader ILoader)
	Update() bool
	Check() bool
	GetName() string
	GetFileName() string
}

type IVoList interface {
	Update() bool
	Check() bool
}

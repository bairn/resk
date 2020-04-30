package infra

type Initializer interface {
	Init()
}

type InitializeRegister struct {
	Initializers []Initializer
}

func (i *InitializeRegister) Register(ai Initializer) {
	i.Initializers = append(i.Initializers, ai)
}

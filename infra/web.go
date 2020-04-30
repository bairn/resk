package infra

var apiInitializerRegister *InitializeRegister = new(InitializeRegister)

func RegisterApi(ai Initializer) {
	apiInitializerRegister.Register(ai)
}

func GetApiInitializes() []Initializer {
	return apiInitializerRegister.Initializers
}

type WebApiStarter struct {
	BaseStarter
}

func (w *WebApiStarter) Setup(ctx StarterContext) {
	for _, v := range GetApiInitializes() {
		v.Init()
	}
}

package adapters

type Adapter[UA any, SA any] interface {
	CreateUser()
}
